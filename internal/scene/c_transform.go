package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"time"
)

func NewCTransform(pos *pb.Vector2) *CTransform {
	c := &CTransform{
		Component: ecs.NewComponent(C_Transform),
		position:  pos,
		movementReqs: ds.NewRing[*pb.SceneMovementReq](
			ds.RingMinCap[*pb.SceneMovementReq](4),
			ds.RingMaxCap[*pb.SceneMovementReq](64),
		),
	}
	return c
}

type CTransform struct {
	ecs.Component
	position       *pb.Vector2
	movementReqs   *ds.Ring[*pb.SceneMovementReq]
	lastMovement   *pb.SceneMovementReq
	movementEvents *ds.Array[*pb.SceneEvent]
	moved          bool
}

func (c *CTransform) Start() {
	c.movementEvents = ds.NewArray[*pb.SceneEvent](8)
	c.lastMovement = &pb.SceneMovementReq{
		Direction: &pb.Vector2{
			X: 1,
			Y: 0,
		},
		MoveSpeed: 0,
		Timestamp: 0,
	}
	c.PushMovement(c.lastMovement)
}

func (c *CTransform) IsMoved() bool {
	return c.moved
}

func (c *CTransform) PushMovement(movement *pb.SceneMovementReq) {
	movement.Timestamp = time.Now().UnixMilli()
	_ = c.movementReqs.Put(movement)
}

func (c *CTransform) ProcessMovement(nowMs int64, maxX, maxY float32) {
	if c.movementReqs.Available() > 0 {
		for c.movementReqs.Available() > 0 {
			movement, _ := c.movementReqs.Pop()
			if !isStopMove(c.lastMovement) {
				c.updatePosition(movement.Timestamp-c.lastMovement.Timestamp, c.lastMovement, maxX, maxY)
				c.moved = true
			}
			c.lastMovement = movement
		}
	} else {
		if isStopMove(c.lastMovement) {
			c.moved = false
		} else {
			c.updatePosition(nowMs-c.lastMovement.Timestamp, c.lastMovement, maxX, maxY)
			c.moved = true
		}
		c.lastMovement.Timestamp = nowMs
	}
}

func (c *CTransform) Clean() {
	c.movementEvents.Reset()
}

func isStopMove(movement *pb.SceneMovementReq) bool {
	return movement.MoveSpeed == 0 || movement.Direction == nil
}

func (c *CTransform) updatePosition(durMs int64, movement *pb.SceneMovementReq, maxX, maxY float32) {
	secs := float32(durMs) / 1000
	offset := util.Vec2Mul(PbVecToVec2(movement.Direction), movement.MoveSpeed*secs)
	pos := util.Vec2Add(PbVecToVec2(c.position), offset)
	pos.X = util.ClampFloat(pos.X, 0, maxX)
	pos.Y = util.ClampFloat(pos.Y, 0, maxY)
	Vec2ToPbVec(pos, c.position)
	c.movementEvents.Add(&pb.SceneEvent{
		Id: c.Entity().Id(),
		Event: &pb.SceneEvent_Movement{
			Movement: &pb.SceneMovement{
				Position:  c.position,
				Timestamp: movement.Timestamp,
				Duration:  durMs,
			},
		},
	})
}
