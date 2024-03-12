package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"time"
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
	tid, pawn, handler := util.SplitSlc3[int64, *pb.ScenePawn, func(*util.Err)](data)
	if pawn.Pawn.Position == nil {
		handler(util.NewErr(util.EcParamsErr, util.M{
			"error": "not set position",
		}))
	}
	if pawn.Pawn.PawnType == nil {
		handler(util.NewErr(util.EcParamsErr, util.M{
			"error": "not set pawn",
		}))
		return
	}

	pawnId := pawn.Pawn.PawnId
	sd := GetSceneDataByFrame(s.Frame())
	if _, ok := pawn.Pawn.PawnType.(*pb.ScenePawnEvt_Player); ok {
		if len(sd.Admitted) > 0 {
			for _, player := range sd.Admitted {
				if player.Id != pawnId {
					continue
				}
				if player.Entered {
					handler(util.NewErr(util.EcServiceErr, util.M{
						"player id": pawnId,
						"error":     "entry",
					}))
					return
				}
				player.Entered = true
				s.entry(tid, pawn)
				return
			}
			handler(util.NewErr(EcSceneEntry_NoEntry, nil))
			return
		}
	}
	_, ok := s.Scene().GetEntity(pawnId)
	if ok {
		handler(util.NewErr(util.EcServiceErr, util.M{
			"player id": pawnId,
			"error":     "entry",
		}))
		return
	}
	s.entry(tid, pawn)
	handler(nil)
}

func (s *SEntity) entry(tid int64, pawn *pb.ScenePawn) {
	pawnId := pawn.Pawn.PawnId
	e := ecs.NewEntity(pawnId)
	tnf := NewCTransform(pawn.Pawn.Position)
	e.AddComponents(
		NewCBehaviour(&pb.SceneBehaviour{
			Timestamp: time.Now().UnixMilli(),
			BehaviourType: &pb.SceneBehaviour_Idle{
				Idle: &pb.BehaviourIdle{},
			},
		}),
		tnf,
	)
	var cpawn ICPawn
	switch p := pawn.Pawn.PawnType.(type) {
	case *pb.ScenePawnEvt_Player:
		cpawn = NewCPlayer(pawn.PlayerGateNodeId, pawn.PlayerGateAddr, p.Player, s.Frame())
	case *pb.ScenePawnEvt_Monster:
		cpawn = NewCMonster(p.Monster)
	}

	e.AddComponents(
		cpawn,
		NewCEvent(cpawn),
	)
	_ = s.Scene().AddEntity(e)
	s.Scene().TagComponent(tnf, TagCompSceneEntry)
}

func (s *SEntity) onEntityDel(data []any) {
	_, eids := util.SplitSlc2[int64, []string](data)
	for _, eid := range eids {
		s.Scene().TagEntityComponent(eid, C_Transform, TagCompSceneExit)
	}
	s.FrameAfter().Push(func() {
		for _, eid := range eids {
			_ = s.Scene().DelEntity(eid)
		}
	})
}
