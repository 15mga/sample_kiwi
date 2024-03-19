package scene

import (
	"github.com/15mga/kiwi/ecs"
)

func NewSEvent() *SEvent {
	s := &SEvent{
		System: ecs.NewSystem(S_Event),
	}
	return s
}

type SEvent struct {
	ecs.System
}

func (s *SEvent) OnUpdate() {
	s.PTagComponents(string(C_Event), s.prcEvent)
}

func (s *SEvent) prcEvent(component ecs.IComponent) {
	component.(*CEvent).ProcessEvents()
}
