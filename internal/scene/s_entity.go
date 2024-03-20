package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewSEntity() *SEntity {
	return &SEntity{
		System: ecs.NewSystem(S_Entity),
	}
}

type SEntity struct {
	ecs.System
}

func (s *SEntity) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindJob(JobEntityAdd, s.onEntityAdd)
	s.BindJob(JobEntityDel, s.onEntityDel)
}

func (s *SEntity) OnUpdate() {
	s.DoJob(JobEntityAdd)
	s.DoJob(JobEntityDel)
	s.Frame().AfterClearTags(
		TagCompSceneEntry,
		TagCompSceneExit,
	)
}

func (s *SEntity) onEntityAdd(data []any) {
	tid, id, visible, handler := util.SplitSlc4[int64, string, *pb.SceneVisible, func(*util.Err)](data)
	if visible.Position == nil {
		handler(util.NewErr(util.EcParamsErr, util.M{
			"error": "not set position",
		}))
	}
	if visible.PawnType == nil {
		handler(util.NewErr(util.EcParamsErr, util.M{
			"error": "not set pawn",
		}))
		return
	}

	sd := GetSceneDataByFrame(s.Frame())
	if _, ok := visible.PawnType.(*pb.SceneVisible_Player); ok {
		if len(sd.Admitted) > 0 {
			for _, player := range sd.Admitted {
				if player.Id != id {
					continue
				}
				if player.Entered {
					handler(util.NewErr(util.EcServiceErr, util.M{
						"player id": id,
						"error":     "entry",
					}))
					return
				}
				player.Entered = true
				s.entry(tid, id, visible)
				return
			}
			handler(util.NewErr(EcSceneEntry_NoEntry, nil))
			return
		}
	}
	_, ok := s.Scene().GetEntity(id)
	if ok {
		handler(util.NewErr(util.EcServiceErr, util.M{
			"player id": id,
			"error":     "entry",
		}))
		return
	}
	s.entry(tid, id, visible)
	handler(nil)
}

func (s *SEntity) entry(tid int64, id string, visible *pb.SceneVisible) {
	e := ecs.NewEntity(id)
	tile := NewCTile()
	e.AddComponents(
		NewCTransform(visible.Position),
		tile,
	)
	var cpawn ICPawn
	switch p := visible.PawnType.(type) {
	case *pb.SceneVisible_Player:
		cpawn = NewCPlayer(visible.GateNodeId, visible.GateAddr, p.Player, s.Frame())
	case *pb.SceneVisible_Monster:
		cpawn = NewCMonster(p.Monster)
	}

	e.AddComponents(
		cpawn,
		NewCEvent(cpawn),
	)
	_ = s.Scene().AddEntity(e)
	s.Scene().TagComponent(tile, TagCompSceneEntry)
}

func (s *SEntity) onEntityDel(data []any) {
	_, eids := util.SplitSlc2[int64, []string](data)
	for _, eid := range eids {
		s.Scene().TagEntityComponent(eid, C_Tile, TagCompSceneExit)
	}
	s.FrameAfter().Push(func() {
		for _, eid := range eids {
			_ = s.Scene().DelEntity(eid)
		}
	})
}
