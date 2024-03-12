package scene

import "github.com/15mga/kiwi/ecs"

func NewSEvent() *SEvent {
	return &SEvent{
		System: ecs.NewSystem(S_Event),
	}
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
