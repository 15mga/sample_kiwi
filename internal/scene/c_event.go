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
		eventPus:  &pb.SceneEventPus{},
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
	movement  *pb.SceneMovementEvt
	//接受别人的
	eventPus     *pb.SceneEventPus
	rcvInvisible *ds.Array[string]
	rcvVisible   *ds.Array[*pb.ScenePawnEvt]
	rcvBehaviour *ds.Array[*pb.SceneBehaviourEvt]
	rcvMovement  *ds.Array[*pb.SceneMovementEvt]
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
	c.movement = tnfComp.movementEvt
	c.invisible = entityId
	c.visible = &pb.ScenePawnEvt{
		PawnId:   entityId,
		Position: tnfComp.position,
	}
	c.rcvInvisible = ds.NewArray[string](16)
	c.rcvVisible = ds.NewArray[*pb.ScenePawnEvt](16)
	c.rcvBehaviour = ds.NewArray[*pb.SceneBehaviourEvt](64)
	c.rcvMovement = ds.NewArray[*pb.SceneMovementEvt](128)
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

func (c *CEvent) PushInvisible(invisible []string) {
	c.rcvInvisible.AddRange(invisible...)
}

func (c *CEvent) PushVisible(visible []*pb.ScenePawnEvt) {
	c.rcvVisible.AddRange(visible...)
}

func (c *CEvent) PushMovement(movement []*pb.SceneMovementEvt) {
	c.rcvMovement.AddRange(movement...)
}

func (c *CEvent) PushBehaviour(behaviour []*pb.SceneBehaviourEvt) {
	c.rcvBehaviour.AddRange(behaviour...)
}

func (c *CEvent) ProcessEvents() {
	if c.rcvVisible.Count() == 0 &&
		c.rcvInvisible.Count() == 0 &&
		c.rcvMovement.Count() == 0 {
		return
	}
	dirty := false
	if c.rcvInvisible.Count() > 0 {
		c.eventPus.Invisible = c.rcvInvisible.Values()
		dirty = true
	} else {
		c.eventPus.Invisible = nil
	}
	if c.rcvVisible.Count() > 0 {
		c.eventPus.Visible = c.rcvVisible.Values()
		dirty = true
	} else {
		c.eventPus.Visible = nil
	}
	if c.rcvMovement.Count() > 0 {
		c.eventPus.Movement = c.rcvMovement.Values()
		dirty = true
	} else {
		c.eventPus.Movement = nil
	}
	if c.rcvBehaviour.Count() > 0 {
		c.eventPus.Behaviour = c.rcvBehaviour.Values()
		dirty = true
	} else {
		c.eventPus.Behaviour = nil
	}
	if !dirty {
		return
	}
	c.pawn.ProcessEvents(c)
	c.rcvInvisible.Reset()
	c.rcvVisible.Reset()
	c.rcvMovement.Reset()
	c.rcvBehaviour.Reset()
}
