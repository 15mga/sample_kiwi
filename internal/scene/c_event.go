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
		rcvEvents: ds.NewArray[*pb.SceneEvent](64),
	}
}

type IPawn interface {
	Data() util.IMsg
}

type CEvent struct {
	ecs.Component
	//自己给别人的
	pawn      ICPawn
	visible   *pb.SceneEvent
	invisible *pb.SceneEvent
	//接受别人的
	eventPus  *pb.SceneEventPus
	rcvEvents *ds.Array[*pb.SceneEvent]
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
	c.invisible = &pb.SceneEvent{
		Id: entityId,
		Event: &pb.SceneEvent_Invisible{
			Invisible: &pb.SceneInvisible{},
		},
	}
	var visible *pb.SceneVisible
	switch pawn.(type) {
	case *CPlayer:
		visible = &pb.SceneVisible{
			PawnType: &pb.SceneVisible_Player{
				Player: pawn.Data().(*pb.ScenePlayer),
			},
		}
	case *CMonster:
		visible = &pb.SceneVisible{
			PawnType: &pb.SceneVisible_Monster{
				Monster: pawn.Data().(*pb.SceneMonster),
			},
		}
	}
	visible.Position = tnfComp.position
	c.visible = &pb.SceneEvent{
		Id: entityId,
		Event: &pb.SceneEvent_Visible{
			Visible: visible,
		},
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

func (c *CEvent) PushEvents(events []*pb.SceneEvent) {
	c.rcvEvents.AddRange(events...)
}

func (c *CEvent) ProcessEvents() {
	if c.rcvEvents.Count() == 0 {
		return
	}
	c.eventPus.Events = c.rcvEvents.Values()
	c.pawn.ProcessEvents(c)
	c.rcvEvents.Reset()
	c.eventPus.Events = nil
}
