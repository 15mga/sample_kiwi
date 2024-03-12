package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewCEvent(pawn ICPawn) *CEvent {
	return &CEvent{
		Component: ecs.NewComponent(C_Event),
		pawn:      pawn,
		//transformEvents: ds.NewKSet[string, *pb.SceneTransformEvt](16, func(evt *pb.SceneTransformEvt) string {
		//	return evt.PawnId
		//}),
		//visibleEvents: ds.NewKSet[string, *pb.ScenePawnEvt](16, func(evt *pb.ScenePawnEvt) string {
		//	return evt.PawnId
		//}),
		//invisibleEvents: ds.NewArray[string](16),
		//behaviourEvents: ds.NewKSet[string, *pb.SceneBehaviourEvt](16, func(evt *pb.SceneBehaviourEvt) string {
		//	return evt.PawnId
		//}),
		transformEvents: ds.NewArray[*pb.SceneTransformEvt](16),
		visibleEvents:   ds.NewArray[*pb.ScenePawnEvt](16),
		invisibleEvents: ds.NewArray[string](16),
		behaviourEvents: ds.NewArray[*pb.SceneBehaviourEvt](16),
		Data:            &pb.SceneEventPus{},
	}
}

type IPawn interface {
	Data() util.IMsg
}

type CEvent struct {
	ecs.Component
	//自己给别人的
	pawn      ICPawn
	visible   *pb.ScenePawnEvt
	invisible string
	behaviour *pb.SceneBehaviourEvt
	transform *pb.SceneTransformEvt
	//接受别人的
	Data *pb.SceneEventPus
	//transformEvents *ds.KSet[string, *pb.SceneTransformEvt]
	//visibleEvents   *ds.KSet[string, *pb.ScenePawnEvt]
	//invisibleEvents *ds.Array[string]
	//behaviourEvents *ds.KSet[string, *pb.SceneBehaviourEvt]
	transformEvents *ds.Array[*pb.SceneTransformEvt]
	visibleEvents   *ds.Array[*pb.ScenePawnEvt]
	invisibleEvents *ds.Array[string]
	behaviourEvents *ds.Array[*pb.SceneBehaviourEvt]
}

func (c *CEvent) Start() {
	comp, _ := c.Entity().GetComponent(C_Transform)
	pawn, ok := c.getPawn()
	if !ok {
		kiwi.Fatal2(util.EcNotExist, util.M{
			"component": "pawn",
		})
	}
	tnfComp := comp.(*CTransform)
	entityId := c.Entity().Id()
	c.transform = tnfComp.TransformEventData
	c.invisible = entityId
	c.visible = &pb.ScenePawnEvt{
		PawnId:   entityId,
		Position: tnfComp.Position,
	}
	switch pawn.(type) {
	case *CPlayer:
		c.visible.PawnType = &pb.ScenePawnEvt_Player{
			Player: pawn.Data().(*pb.ScenePlayer),
		}
	case *CMonster:
		c.visible.PawnType = &pb.ScenePawnEvt_Monster{
			Monster: pawn.Data().(*pb.SceneMonster),
		}
	}
}

func (c *CEvent) getPawn() (IPawn, bool) {
	comp, ok := c.Entity().GetComponent(C_Player)
	if ok {
		return comp.(*CPlayer), true
	}
	comp, ok = c.Entity().GetComponent(C_Monster)
	if ok {
		return comp.(*CMonster), true
	}
	return nil, false
}

func (c *CEvent) AddTransformEvent(target *CEvent) {
	if c.pawn == nil {
		return
	} //自己的也发送用于纠正
	//_ = c.transformEvents.Add(target.transform)
	c.transformEvents.Add(target.transform)
}

func (c *CEvent) AddVisibleEvents(target *CEvent) {
	if c.pawn == nil || target == c {
		return
	}
	//_ = c.visibleEvents.Add(target.visible)
	c.visibleEvents.Add(target.visible)
}

func (c *CEvent) AddInvisibleEvent(target *CEvent) {
	if c.pawn == nil || target == c {
		return
	}
	c.invisibleEvents.Add(target.invisible)
}

func (c *CEvent) AddBehaviourEvent(target *CEvent) {
	if c.pawn == nil || target == c {
		return
	}
	//_ = c.behaviourEvents.Add(target.behaviour)
	c.behaviourEvents.Add(target.behaviour)
}

func (c *CEvent) ProcessEvents() {
	if c.pawn == nil {
		return
	}

	c.Data.Invisible = c.invisibleEvents.Values()
	c.Data.Visible = c.visibleEvents.Values()
	c.Data.Transform = c.transformEvents.Values()
	c.Data.Behaviour = c.behaviourEvents.Values()

	if len(c.Data.Invisible) > 0 ||
		len(c.Data.Visible) > 0 ||
		len(c.Data.Transform) > 0 ||
		len(c.Data.Behaviour) > 0 {
		c.pawn.ProcessEvents(c)
		c.invisibleEvents.Reset()
		c.visibleEvents.Reset()
		c.transformEvents.Reset()
		c.behaviourEvents.Reset()
	}
}
