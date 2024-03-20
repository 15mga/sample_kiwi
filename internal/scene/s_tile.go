package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/worker"
	"math/rand"
	"strconv"
)

type TileState int8

const (
	TileStateEntry TileState = iota
	TileStateExit
	TileStateStay
)

func NewTile(x, y int) *Tile {
	t := &Tile{
		id: util.Vec2Int{
			X: x,
			Y: y,
		},
		interest: ds.NewArray[util.Vec2Int](8),
		events:   ds.NewArray[*pb.SceneEvent](8),
		stayInvisibleEvents: ds.NewKSet[string, *pb.SceneEvent](8, func(event *pb.SceneEvent) string {
			return event.Id
		}),
		stayVisibleEvents: ds.NewKSet[string, *pb.SceneEvent](8, func(event *pb.SceneEvent) string {
			return event.Id
		}),
		stayDelCache: ds.NewArray[*CTile](8),
		cTiles: ds.NewKSet[string, *CTile](64, func(tile *CTile) string {
			return tile.Entity().Id()
		}),
	}
	t.tag = getTileTag(t.id)
	return t
}

type Tile struct {
	tag                 string
	id                  util.Vec2Int
	interest            *ds.Array[util.Vec2Int]
	events              *ds.Array[*pb.SceneEvent]
	stayInvisibleEvents *ds.KSet[string, *pb.SceneEvent]
	stayVisibleEvents   *ds.KSet[string, *pb.SceneEvent]
	stayDelCache        *ds.Array[*CTile]
	cTiles              *ds.KSet[string, *CTile]
}

func (t *Tile) AddCTile(ct *CTile) {
	_ = t.cTiles.AddNX(ct)
}

func (t *Tile) CacheDelCTile(ct ...*CTile) {
	t.stayDelCache.AddRange(ct...)
}

func (t *Tile) Clean() {
	for _, ct := range t.stayDelCache.Values() {
		t.stayInvisibleEvents.Del(ct.cEvent.invisible.Id)
		t.stayVisibleEvents.Del(ct.cEvent.visible.Id)
		t.cTiles.Del(ct.Entity().Id())
	}
	t.events.Reset()
	t.stayDelCache.Reset()
}

func NewSTile(tileSize, sceneWidth, sceneHeight, fovLaps int) *STile {
	s := &STile{
		System:      ecs.NewSystem(S_Tile),
		tileSize:    tileSize,
		sceneWidth:  float32(sceneWidth),
		sceneHeight: float32(sceneHeight),
		fovLaps:     fovLaps,
	}
	s.tileXCount = sceneWidth / tileSize
	if sceneWidth%tileSize > 0 {
		s.tileXCount++
	}
	s.tileYCount = sceneHeight / tileSize
	if sceneHeight%tileSize > 0 {
		s.tileYCount++
	}
	return s
}

type STile struct {
	ecs.System
	tiles                   *ds.KSet[string, *Tile]
	fovLaps                 int
	tileSize                int
	sceneWidth, sceneHeight float32
	tileXCount, tileYCount  int
}

func (s *STile) OnAfterStart() {
	s.tiles = ds.NewKSet[string, *Tile](s.tileXCount*s.tileYCount, func(t *Tile) string {
		return t.tag
	})
	for x := 0; x < s.tileXCount; x++ {
		for y := 0; y < s.tileYCount; y++ {
			_ = s.tiles.Add(NewTile(x, y))
		}
	}
	for _, tile := range s.tiles.Values() {
		s.GetInterestTiles(tile.id, tile.interest)
	}
}

func (s *STile) OnUpdate() {
	s.processSceneEntry()
	s.processSceneExit()
	s.processEvents()
}

func (s *STile) processSceneExit() {
	_, _ = s.PTagComponents(TagCompSceneExit, func(component ecs.IComponent) {
		tile := component.(*CTile)
		tile.state = TileStateExit
	})
}

func (s *STile) processSceneEntry() {
	components, ok := s.PTagComponents(TagCompSceneEntry, func(component ecs.IComponent) {
		tile := component.(*CTile)
		tile.state = TileStateEntry
		tile.InitTile(s.posToTile(tile.cTnf.position))
		s.GetInterestTiles(tile.currTile, tile.interest)
	})
	if !ok {
		return
	}

	for _, component := range components {
		ct := component.(*CTile)
		s.getTile(ct.currTile).AddCTile(ct)
	}
}

func (s *STile) processEvents() {
	//components, ok := s.PTagComponentsToFnLink(string(C_Tile), func(component ecs.IComponent, link *ds.FnLink) {
	//	ct := component.(*CTile)
	//	if !ct.cTnf.IsMoved() {
	//		return
	//	}
	//	t := s.posToTile(ct.cTnf.position)
	//	ct.UpdateTile(t)
	//	if !ct.IsTileChanged() {
	//		return
	//	}
	//	s.getInterestTileChanged(ct)
	//	link.Push(func() {
	//		s.getTile(ct.currTile).AddCTile(ct)
	//	})
	//})
	components, ok := s.Scene().GetTagComponents(string(C_Tile))
	if !ok {
		return
	}
	worker.PFilter[ecs.IComponent, *CTile](components, func(component ecs.IComponent) (*CTile, bool) {
		ct := component.(*CTile)
		if !ct.cTnf.IsMoved() {
			return nil, false
		}
		t := s.posToTile(ct.cTnf.position)
		ct.UpdateTile(t)
		if !ct.IsTileChanged() {
			return nil, false
		}
		s.getInterestTileChanged(ct)
		return ct, true
	}, func(cTiles []*CTile) {
		for _, ct := range cTiles {
			s.getTile(ct.currTile).AddCTile(ct)
		}
	})
	//添加事件到格子
	//now := s.Frame().NowMillSecs()
	worker.P[*Tile](s.tiles.Values(), func(tile *Tile) {
		if tile.cTiles.Count() == 0 {
			return
		}
		//for _, event := range tile.stayInvisibleEvents.Values() {
		//	event.Event.(*pb.SceneEvent_Invisible).Invisible.Timestamp = now
		//}
		//for _, event := range tile.stayVisibleEvents.Values() {
		//	event.Event.(*pb.SceneEvent_Visible).Visible.Timestamp = now
		//}
		var delCTiles []*CTile
		for _, ct := range tile.cTiles.Values() {
			switch ct.state {
			case TileStateExit:
				delCTiles = append(delCTiles, ct)
				tile.events.Add(ct.cEvent.invisible)
			case TileStateEntry:
				_ = tile.stayInvisibleEvents.AddNX(ct.cEvent.invisible)
				_ = tile.stayVisibleEvents.AddNX(ct.cEvent.visible)
				tile.events.Add(ct.cEvent.visible)
			case TileStateStay:
				if !ct.cTnf.IsMoved() {
					continue
				}
				tile.events.AddRange(ct.cTnf.movementEvents.Values()...)
				if !ct.IsTileChanged() {
					continue
				}
				if ct.prevTile.Equal(tile.id) {
					delCTiles = append(delCTiles, ct)
					tile.events.Add(ct.cEvent.invisible)
				} else if ct.currTile.Equal(tile.id) {
					_ = tile.stayInvisibleEvents.AddNX(ct.cEvent.invisible)
					_ = tile.stayVisibleEvents.AddNX(ct.cEvent.visible)
					tile.events.Add(ct.cEvent.visible)
				}
			}
		}
		tile.CacheDelCTile(delCTiles...)
	})
	worker.P[ecs.IComponent](components, func(component ecs.IComponent) {
		ct := component.(*CTile)
		switch ct.state {
		case TileStateEntry:
			for _, id := range ct.interest.Values() {
				tile := s.getTile(id)
				//todo 这里新进入的会重复给了visible事件，如果不希望重复给，把stayVisible放在后面处理
				ct.cEvent.PushEvents(tile.stayVisibleEvents.Values())
				ct.cEvent.PushEvents(tile.events.Values())
			}
			ct.state = TileStateStay
		case TileStateStay:
			if ct.IsTileChanged() {
				for _, id := range ct.entry.Values() {
					tile := s.getTile(id)
					ct.cEvent.PushEvents(tile.stayVisibleEvents.Values())
					ct.cEvent.PushEvents(tile.events.Values())
				}
				for _, id := range ct.stay.Values() {
					tile := s.getTile(id)
					ct.cEvent.PushEvents(tile.events.Values())
				}
				for _, id := range ct.exit.Values() {
					tile := s.getTile(id)
					ct.cEvent.PushEvents(tile.stayInvisibleEvents.Values())
				}
				ct.interest.Clean()
				s.GetInterestTiles(ct.currTile, ct.interest)
				ct.ResetTileChanged()
			} else {
				for _, id := range ct.interest.Values() {
					tile := s.getTile(id)
					//for _, event := range tile.events.Values() {
					//	switch event.Event.(type) {
					//	case *pb.SceneEvent_Visible:
					//		kiwi.Debug("test", nil)
					//	}
					//}
					ct.cEvent.PushEvents(tile.events.Values())
				}
			}
		}
	})
	for _, tile := range s.tiles.Values() {
		tile.Clean()
	}
	//worker.P[*Tile](s.tiles.Values(), func(tile *Tile) {
	//	tile.Clean()
	//})
}

func (s *STile) GetInterestTiles(center util.Vec2Int, arr *ds.Array[util.Vec2Int]) {
	sx := center.X - s.fovLaps
	if sx < 0 {
		sx = 0
	}
	ex := center.X + s.fovLaps
	if ex >= s.tileXCount {
		ex = s.tileXCount - 1
	}
	sy := center.Y - s.fovLaps
	if sy < 0 {
		sy = 0
	}
	ey := center.Y + s.fovLaps
	if ey >= s.tileYCount {
		ey = s.tileYCount - 1
	}
	for x := sx; x <= ex; x++ {
		for y := sy; y <= ey; y++ {
			arr.Add(util.Vec2Int{
				X: x,
				Y: y,
			})
		}
	}
}

func (s *STile) getInterestTileChanged(tile *CTile) {
	tile.exit.Clean()
	tile.entry.Clean()
	tile.stay.Clean()
	ox := tile.currTile.X - tile.prevTile.X
	oy := tile.currTile.Y - tile.prevTile.Y
	n := s.fovLaps << 1
	if ox > n || oy > n {
		s.GetInterestTiles(tile.prevTile, tile.exit)
		s.GetInterestTiles(tile.currTile, tile.entry)
		return
	}
	psx := util.MaxInt(tile.prevTile.X-s.fovLaps, 0)
	psy := util.MaxInt(tile.prevTile.Y-s.fovLaps, 0)
	pex := util.MinInt(tile.prevTile.X+s.fovLaps, s.tileXCount-1)
	pey := util.MinInt(tile.prevTile.Y+s.fovLaps, s.tileYCount-1)
	csx := util.MaxInt(tile.currTile.X-s.fovLaps, 0)
	csy := util.MaxInt(tile.currTile.Y-s.fovLaps, 0)
	cex := util.MinInt(tile.currTile.X+s.fovLaps, s.tileXCount-1)
	cey := util.MinInt(tile.currTile.Y+s.fovLaps, s.tileYCount-1)
	switch {
	case ox < 0:
		//y进入
		for x := csx; x < psx; x++ {
			for y := csy; y <= cey; y++ {
				tile.entry.Add(util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		//y退出
		for x := cex + 1; x <= pex; x++ {
			for y := psy; y <= pey; y++ {
				tile.exit.Add(util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		switch {
		case oy < 0:
			for x := psx; x <= cex; x++ {
				//x进入
				for y := csy; y < psy; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy == 0:
			for x := psx; x <= cex; x++ {
				//不动
				for y := psy; y <= pey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := psx; x <= cex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		}
	case ox == 0:
		switch {
		case oy < 0:
			for x := csx; x <= cex; x++ {
				//x进入
				for y := csy; y < psy; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := csx; x <= cex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		}
	case ox > 0:
		//y退出
		for x := psx; x < csx; x++ {
			for y := psy; y <= pey; y++ {
				tile.exit.Add(util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		//y进入
		for x := pex + 1; x <= cex; x++ {
			for y := csy; y <= cey; y++ {
				tile.entry.Add(util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		switch {
		case oy < 0:
			for x := csx; x <= pex; x++ {
				//x进入
				for y := csy; y < psy; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy == 0:
			for x := csx; x <= pex; x++ {
				//不动
				for y := csy; y <= cey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := csx; x <= pex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tile.exit.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tile.stay.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tile.entry.Add(util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		}
	}
}

func (s *STile) posToTile(pos *pb.Vector2) util.Vec2Int {
	return util.Vec2Int{
		X: int(pos.X) / s.tileSize,
		Y: int(pos.Y) / s.tileSize,
	}
}

func (s *STile) GenRandPos(pos *pb.Vector2) {
	pos.X = rand.Float32() * s.sceneWidth
	pos.Y = rand.Float32() * s.sceneHeight
}

func (s *STile) getFOV(tx, ty int) (minX, maxX, minY, maxY int) {
	minX = util.MaxInt(0, tx-s.fovLaps)
	maxX = util.MaxInt(s.tileXCount, tx+s.fovLaps+1)
	minY = util.MaxInt(0, ty-s.fovLaps)
	maxY = util.MaxInt(s.tileYCount, ty+s.fovLaps+1)
	return
}

func (s *STile) getTile(v util.Vec2Int) *Tile {
	t, _ := s.tiles.Get(getTileTag(v))
	return t
}

func getTileTag(v util.Vec2Int) string {
	return strconv.Itoa(int(v.X)) + "_" + strconv.Itoa(int(v.Y))
}
