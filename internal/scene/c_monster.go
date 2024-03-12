package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/graph"
	"github.com/15mga/kiwi/util"
)

type EMonsterStatus int

const (
	EMonsterStatusIdle EMonsterStatus = iota
	EMonsterStatusWalk
)

func NewCMonster(data *pb.SceneMonster) *CMonster {
	return &CMonster{
		Component: ecs.NewComponent(C_Monster),
		data:      data,
	}
}

type CMonster struct {
	ecs.Component
	data      *pb.SceneMonster
	originPos *pb.Vector2
	moveRange float32 //移动范围
	graph     graph.IGraph
}

func (c *CMonster) Data() util.IMsg {
	return c.data
}

func (c *CMonster) ProcessEvents(event *CEvent) {

}

func (c *CMonster) Update() {

}
