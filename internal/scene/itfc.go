package scene

import "github.com/15mga/kiwi/ecs"

type ICPawn interface {
	ecs.IComponent
	ProcessEvents(*CEvent)
}
