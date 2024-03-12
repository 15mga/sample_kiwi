package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
)

func NewCBehaviour(data *pb.SceneBehaviour) *CBehaviour {
	return &CBehaviour{
		Component: ecs.NewComponent(C_Behaviour),
		Data:      data,
	}
}

type CBehaviour struct {
	ecs.Component
	Data *pb.SceneBehaviour
}
