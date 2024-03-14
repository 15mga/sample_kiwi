package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/worker"
	"math/rand"
	"strconv"
)

func NewSTransform(tileSize, sceneWidth, sceneHeight, fovLaps int32) *STransform {
	s := &STransform{
		System:      ecs.NewSystem(S_Transform),
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

type STransform struct {
	ecs.System
	tileSize                int32
	sceneWidth, sceneHeight float32
	tileXCount, tileYCount  int32
	fovLaps                 int32
}

func (s *STransform) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindPFnJob(JobMovement, s.onMovement)
	s.BindPJob(JobBehaviour, s.onBehaviour)
}

func (s *STransform) FovLaps() int32 {
	return s.fovLaps
}

func (s *STransform) GenRandPos(pos *pb.Vector2) {
	pos.X = rand.Float32() * s.sceneWidth
	pos.Y = rand.Float32() * s.sceneHeight
}

func (s *STransform) getFOV(tx, ty int32) (minX, maxX, minY, maxY int32) {
	minX = util.MaxInt32(0, tx-s.fovLaps)
	maxX = util.MinInt32(s.tileXCount, tx+s.fovLaps+1)
	minY = util.MaxInt32(0, ty-s.fovLaps)
	maxY = util.MinInt32(s.tileYCount, ty+s.fovLaps+1)
	return
}

func (s *STransform) OnUpdate() {
	s.DoJob(JobMovement)
	s.processSceneExit()
	s.processMoving()
	s.processTileChange()
	s.processSceneEntry()
	s.DoJob(JobBehaviour)
	s.FrameAfter().Push(func() {
		s.Scene().ClearTag(TagCompTileChange)
		s.PTagComponents(string(C_Transform), func(component ecs.IComponent) {
			component.(*CTransform).ClearMovement()
		})
	})
}

func (s *STransform) onMovement(link *ds.FnLink, data []any) {
	tid, eid, movement := util.SplitSlc3[int64, string, *pb.SceneMovement](data)
	e, ok := s.Scene().GetEntity(eid)
	if !ok {
		kiwi.TE2(tid, util.EcNotExist, util.M{
			"entity id": eid,
		})
		return
	}
	c, ok := e.GetComponent(C_Transform)
	if !ok {
		kiwi.TE2(tid, util.EcNotExist, util.M{
			"entity id": eid,
			"component": C_Transform,
		})
		return
	}
	tnf := c.(*CTransform)
	tnf.PushMovement(movement)
}

func (s *STransform) processMoving() {
	components, ok := s.PTagComponents(string(C_Transform), func(component ecs.IComponent) {
		tnf := component.(*CTransform)
		tnf.ProcessMovement(s.Frame().NowMillSecs(), s.sceneWidth, s.sceneHeight)
		if tnf.moved {
			tnf.UpdateTile(s.posToTile(tnf.Position))
		}
	})
	if !ok || len(components) == 0 {
		return
	}
	//这里有并发写入，只能在frame协程
	for _, component := range components {
		tnf := component.(*CTransform)
		if !tnf.moved {
			continue
		}
		if tnf.IsTileChanged() {
			s.Scene().UntagComponent(tnf, getTileTag(tnf.PrevTile))
			s.Scene().TagComponent(tnf, TagCompTileChange, getTileTag(tnf.CurrTile))
		}
		// 这里移动过程中，每一帧都会更新，
		// 如果没有Movement，不发送则需要在客户端做预判与补偿
		//if tnf.IsMovementDirty() {
		var tiles []util.Vec2Int
		s.GetInterestTiles(tnf.CurrTile, &tiles)
		evtComp, _ := tnf.Entity().GetComponent(C_Event)
		origin := evtComp.(*CEvent)
		for _, t := range tiles {
			comps, ok := s.Scene().GetTagComponents(getTileTag(t))
			if !ok {
				continue
			}
			for _, comp := range comps {
				targetComp, _ := comp.Entity().GetComponent(C_Event)
				targetEvent := targetComp.(*CEvent)
				targetEvent.AddTransformEvent(origin)
			}
		}
		//}
	}
}

func (s *STransform) processTileChange() {
	components, ok := s.Scene().GetTagComponents(TagCompTileChange)
	if !ok {
		return
	}
	worker.PToLink[ecs.IComponent, *TileChangeData](components, func(component ecs.IComponent, d *ds.Link[*TileChangeData]) {
		tnf := component.(*CTransform)
		var tcd TileChangeData
		s.getInterestTileChanged(tnf, &tcd)
		d.Push(&tcd)
	}, s.pushInterestEvents)
}

func (s *STransform) pushInterestEvents(d *ds.Link[*TileChangeData]) {
	d.Iter(func(tcd *TileChangeData) {
		origin := tcd.Origin
		for _, t := range tcd.Exit {
			components, ok := s.Scene().GetTagComponents(getTileTag(t))
			if !ok {
				continue
			}
			for _, comp := range components {
				targetComp, _ := comp.Entity().GetComponent(C_Event)
				targetEvent := targetComp.(*CEvent)
				if targetEvent == origin {
					continue
				}
				targetEvent.AddInvisibleEvent(origin)
				origin.AddInvisibleEvent(targetEvent)
			}
		}
		for _, t := range tcd.Entry {
			components, ok := s.Scene().GetTagComponents(getTileTag(t))
			if !ok {
				continue
			}
			for _, component := range components {
				targetComp, _ := component.Entity().GetComponent(C_Event)
				targetEvent := targetComp.(*CEvent)
				if targetEvent == origin {
					continue
				}
				targetEvent.AddVisibleEvents(origin)
				origin.AddVisibleEvents(targetEvent)
			}
		}
	})
}

func (s *STransform) processSceneExit() {
	components, ok := s.Scene().GetTagComponents(TagCompSceneExit)
	if !ok {
		return
	}

	//广播事件
	for _, component := range components {
		c, _ := component.Entity().GetComponent(C_Event)
		origin := c.(*CEvent)
		var tiles []util.Vec2Int
		tnf := component.(*CTransform)
		s.GetInterestTiles(tnf.CurrTile, &tiles)
		for _, t := range tiles {
			comps, ok := s.Scene().GetTagComponents(getTileTag(t))
			if !ok {
				continue
			}
			for _, comp := range comps {
				targetComp, _ := comp.Entity().GetComponent(C_Event)
				targetEvent := targetComp.(*CEvent)
				targetEvent.AddInvisibleEvent(origin)
			}
		}
	}
}

func (s *STransform) processSceneEntry() {
	entries, ok := s.Scene().GetTagComponents(TagCompSceneEntry)
	if !ok {
		return
	}

	//初始化格子
	entryTagMap := make(map[string]string, 32)
	for _, entry := range entries {
		tnf := entry.(*CTransform)
		tnf.InitTile(s.posToTile(tnf.Position))
		//新进入格子
		entryTag := getTileEntryTag(tnf.CurrTile)
		s.Scene().TagComponent(tnf, entryTag)
		entryTagMap[entryTag] = getTileTag(tnf.CurrTile)
	}

	for _, entry := range entries {
		tnf := entry.(*CTransform)
		var tiles []util.Vec2Int
		s.GetInterestTiles(tnf.CurrTile, &tiles)
		c, _ := entry.Entity().GetComponent(C_Event)
		origin := c.(*CEvent)
		for _, t := range tiles {
			targets, ok := s.Scene().GetTagComponents(getTileTag(t))
			if ok {
				worker.PTo[ecs.IComponent, *CEvent](targets, func(target ecs.IComponent) *CEvent {
					targetComp, _ := target.Entity().GetComponent(C_Event)
					targetEvent := targetComp.(*CEvent)
					targetEvent.AddVisibleEvents(origin)
					return targetEvent
				}, func(targetEvents []*CEvent) {
					for _, target := range targetEvents {
						origin.AddVisibleEvents(target)
					}
				})
			}
			entryTargets, ok := s.Scene().GetTagComponents(getTileEntryTag(t))
			if ok {
				worker.P[ecs.IComponent](entryTargets, func(target ecs.IComponent) {
					targetComp, _ := target.Entity().GetComponent(C_Event)
					targetEvent := targetComp.(*CEvent)
					targetEvent.AddVisibleEvents(origin)
				})
			}
		}
	}
	s.FrameAfter().Push(func() {
		//清理有新实体进入格子标签
		for entryTileTag, tileTag := range entryTagMap {
			s.Scene().TransferComponentTag(tileTag, entryTileTag)
		}
	})
}

func (s *STransform) onBehaviour(data []any) {
	//todo
}

func (s *STransform) GetInterestTiles(center util.Vec2Int, slc *[]util.Vec2Int) {
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
			*slc = append(*slc, util.Vec2Int{
				X: x,
				Y: y,
			})
		}
	}
}

func (s *STransform) getInterestTileChanged(tnf *CTransform, tcd *TileChangeData) {
	tcd.Origin = tnf.cEvent
	ox := tnf.CurrTile.X - tnf.PrevTile.X
	oy := tnf.CurrTile.Y - tnf.PrevTile.Y
	n := s.fovLaps << 1
	if ox > n || oy > n {
		s.GetInterestTiles(tnf.PrevTile, &tcd.Exit)
		s.GetInterestTiles(tnf.CurrTile, &tcd.Entry)
		return
	}
	psx := util.MaxInt32(tnf.PrevTile.X-s.fovLaps, 0)
	psy := util.MaxInt32(tnf.PrevTile.Y-s.fovLaps, 0)
	pex := util.MinInt32(tnf.PrevTile.X+s.fovLaps, s.tileXCount-1)
	pey := util.MinInt32(tnf.PrevTile.Y+s.fovLaps, s.tileYCount-1)
	csx := util.MaxInt32(tnf.CurrTile.X-s.fovLaps, 0)
	csy := util.MaxInt32(tnf.CurrTile.Y-s.fovLaps, 0)
	cex := util.MinInt32(tnf.CurrTile.X+s.fovLaps, s.tileXCount-1)
	cey := util.MinInt32(tnf.CurrTile.Y+s.fovLaps, s.tileYCount-1)
	switch {
	case ox < 0:
		//y进入
		for x := csx; x < psx; x++ {
			for y := csy; y <= cey; y++ {
				tcd.Entry = append(tcd.Entry, util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		//y退出
		for x := cex + 1; x <= pex; x++ {
			for y := psy; y <= pey; y++ {
				tcd.Exit = append(tcd.Exit, util.Vec2Int{
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
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy == 0:
			for x := psx; x <= cex; x++ {
				//不动
				for y := psy; y <= pey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := psx; x <= cex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
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
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := csx; x <= cex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
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
				tcd.Exit = append(tcd.Exit, util.Vec2Int{
					X: x,
					Y: y,
				})
			}
		}
		//y进入
		for x := pex + 1; x <= cex; x++ {
			for y := csy; y <= cey; y++ {
				tcd.Entry = append(tcd.Entry, util.Vec2Int{
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
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := psy; y <= cey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x退出
				for y := cey + 1; y <= pey; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy == 0:
			for x := csx; x <= pex; x++ {
				//不动
				for y := csy; y <= cey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		case oy > 0:
			for x := csx; x <= pex; x++ {
				//x退出
				for y := psy; y < csy; y++ {
					tcd.Exit = append(tcd.Exit, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//不动
				for y := csy; y <= pey; y++ {
					tcd.Stay = append(tcd.Stay, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
				//x进入
				for y := pey + 1; y <= cey; y++ {
					tcd.Entry = append(tcd.Entry, util.Vec2Int{
						X: x,
						Y: y,
					})
				}
			}
		}
	}
}

func (s *STransform) posToTile(pos *pb.Vector2) util.Vec2Int {
	return util.Vec2Int{
		X: int32(pos.X) / s.tileSize,
		Y: int32(pos.Y) / s.tileSize,
	}
}

func getTileTag(v util.Vec2Int) string {
	return strconv.Itoa(int(v.X)) + "_" + strconv.Itoa(int(v.Y))
}
func getTileEntryTag(v util.Vec2Int) string {
	return strconv.Itoa(int(v.X)) + "_" + strconv.Itoa(int(v.Y)) + "_entry"
}
