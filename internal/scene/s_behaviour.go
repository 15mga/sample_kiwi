package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewSBehaviour() *SBehaviour {
	return &SBehaviour{
		System: ecs.NewSystem(S_Behaviour),
	}
}

type SBehaviour struct {
	ecs.System
	sysTnf *STransform
}

func (s *SBehaviour) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindPJob(JobBehaviour, s.onBehaviour)
}

func (s *SBehaviour) OnAfterStart() {
	s.System.OnAfterStart()
	sysTnf, _ := s.Frame().GetSystem(S_Transform)
	s.sysTnf = sysTnf.(*STransform)
}

func (s *SBehaviour) onBehaviour(data []any) {
	tid, eid, behaviour :=
		util.SplitSlc3[int64, string, *pb.SceneBehaviour](data)
	e, ok := s.Scene().GetEntity(eid)
	if !ok {
		kiwi.TE2(tid, util.EcNotExist, util.M{
			"entity id": eid,
		})
		return
	}

	component, _ := e.GetComponent(C_Behaviour)
	component.(*CBehaviour).Data = behaviour

	component, _ = e.GetComponent(C_Transform)
	tnf := component.(*CTransform)
	var tiles []util.Vec2Int
	s.sysTnf.GetInterestTiles(tnf.CurrTile, &tiles)
	c, _ := component.Entity().GetComponent(C_Event)
	origin := c.(*CEvent)
	for _, t := range tiles {
		comps, ok := s.Scene().GetTagComponents(getTileTag(t))
		if !ok {
			continue
		}
		for _, comp := range comps {
			event := comp.(*CEvent)
			event.AddBehaviourEvent(origin)
		}
	}
}
