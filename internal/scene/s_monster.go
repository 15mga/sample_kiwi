package scene

import "github.com/15mga/kiwi/ecs"

func NewSMonster() *SMonster {
	s := &SMonster{
		System: ecs.NewSystem(S_Monster),
	}
	return s
}

type SMonster struct {
	ecs.System
}

func (s *SMonster) OnBeforeStart() {
	s.System.OnBeforeStart()
}

func (s *SMonster) OnUpdate() {
	s.PTagComponents(string(C_Monster), 128, func(component ecs.IComponent) {
		component.(*CMonster).Update()
	})
}
