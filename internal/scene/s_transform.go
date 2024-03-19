package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewSTransform(sceneWidth, sceneHeight int) *STransform {
	s := &STransform{
		System:      ecs.NewSystem(S_Transform),
		sceneWidth:  float32(sceneWidth) - 0.1,
		sceneHeight: float32(sceneHeight) - 0.1,
	}
	return s
}

type STransform struct {
	ecs.System
	sceneWidth, sceneHeight float32
}

func (s *STransform) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindPFnJob(JobMovement, s.onMovement)
}

func (s *STransform) OnUpdate() {
	s.DoJob(JobMovement)
	s.processMove()

	//清理工作
	s.FrameAfter().Push(func() {
		s.PTagComponents(TagCompMove, func(component ecs.IComponent) {
			component.(*CTransform).Clean()
		})
	})
}

func (s *STransform) onMovement(link *ds.FnLink, data []any) {
	tid, eid, movement := util.SplitSlc3[int64, string, *pb.SceneMovementReq](data)
	e, ok := s.Scene().GetEntity(eid)
	if !ok {
		kiwi.TE2(tid, util.EcNotExist, util.M{
			"entity id": eid,
		})
		return
	}
	tnf := e.MGetComponent(C_Transform).(*CTransform)
	tnf.PushMovement(movement)
	link.Push(func() {
		s.Scene().TagComponent(tnf, TagCompMove)
	})
}

func (s *STransform) processMove() {
	_, _ = s.PTagComponents(TagCompMove, func(component ecs.IComponent) {
		tnf := component.(*CTransform)
		tnf.ProcessMovement(s.Frame().NowMillSecs(), s.sceneWidth, s.sceneHeight)
	})
}
