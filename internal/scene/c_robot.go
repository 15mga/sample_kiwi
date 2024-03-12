package scene

import (
	"github.com/15mga/kiwi/ecs"
	"math/rand"
)

func NewCRobot() *CRobot {
	return &CRobot{
		Component:   ecs.NewComponent(C_Robot),
		remainingMs: 1000,
	}
}

type CRobot struct {
	ecs.Component
	remainingMs int64
}

func (c *CRobot) Update(ms int64) bool {
	c.remainingMs -= ms
	if c.remainingMs > 0 {
		return false
	}
	c.remainingMs = 5000 + rand.Int63n(5000)
	//return rand.Intn(3) == 0
	return true
}
