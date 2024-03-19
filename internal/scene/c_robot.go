package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"math/rand"
)

func NewCRobot() *CRobot {
	return &CRobot{
		Component:   ecs.NewComponent(C_Robot),
		remainingMs: 1000,
		dir:         &pb.Vector2{},
	}
}

type CRobot struct {
	ecs.Component
	remainingMs int64
	status      int
	dir         *pb.Vector2
	speed       float32
	rand        *rand.Rand
}

func (c *CRobot) Init() {
	c.rand = rand.New(rand.NewSource(rand.Int63()))
}

func (c *CRobot) Update(ms int64) bool {
	c.remainingMs -= ms
	if c.remainingMs > 0 {
		return false
	}
	n := c.rand.Int63n(5000)
	c.remainingMs = 5000 + n
	if n%3 > 0 {
		c.speed = 0
	} else {
		c.speed = 4
		Vec2ToPbVec(util.RandDir(), c.dir)
	}
	return true
}
