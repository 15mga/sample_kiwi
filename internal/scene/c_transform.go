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
		Position:  pos,
		movements: ds.NewRing[*pb.SceneMovement](
			ds.RingMinCap[*pb.SceneMovement](4),
			ds.RingMaxCap[*pb.SceneMovement](64),
		),
	}
	return c
}

type CTransform struct {
	ecs.Component
	Position           *pb.Vector2
	movements          *ds.Ring[*pb.SceneMovement]
	lastMovement       *pb.SceneMovement
	TransformEventData *pb.SceneTransformEvt
	PrevTile           util.Vec2Int
	CurrTile           util.Vec2Int
	cEvent             *CEvent
	tileChanged        bool
	moving             bool
}

func (c *CTransform) Start() {
	comp, _ := c.Entity().GetComponent(C_Event)
	c.cEvent = comp.(*CEvent)
	c.TransformEventData = &pb.SceneTransformEvt{
		PawnId: c.Entity().Id(),
	}
	c.lastMovement = &pb.SceneMovement{
		Direction: &pb.Vector2{
			X: 1,
			Y: 0,
		},
		MoveSpeed: 0,
		Timestamp: 0,
	}
	c.PushMovement(c.lastMovement)
}

func (c *CTransform) InitTile(tile util.Vec2Int) {
	c.CurrTile = tile
}

func (c *CTransform) PushMovement(movement *pb.SceneMovement) {
	_ = c.movements.Put(&pb.SceneMovement{
		Direction: movement.Direction,
		MoveSpeed: movement.MoveSpeed,
		Timestamp: time.Now().UnixMilli(),
	})
}

func (c *CTransform) ProcessMovement(nowMs int64, maxX, maxY float32) {
	if c.movements.Available() > 0 {
		for c.movements.Available() > 0 {
			movement, _ := c.movements.Pop()
			c.updatePosition(movement.Timestamp-c.lastMovement.Timestamp, c.lastMovement, maxX, maxY)
			c.lastMovement = movement
			c.moving = c.lastMovement.MoveSpeed > 0
		}
	} else {
		c.updatePosition(nowMs-c.lastMovement.Timestamp, c.lastMovement, maxX, maxY)
		c.lastMovement.Timestamp = nowMs
	}
}

func (c *CTransform) ClearMovement() {
	c.tileChanged = false
	c.TransformEventData.Movement = nil
}

func (c *CTransform) IsMoving() bool {
	return c.moving
}

func (c *CTransform) IsTileChanged() bool {
	return c.tileChanged
}

func isStopMove(movement *pb.SceneMovement) bool {
	return movement.MoveSpeed == 0 || movement.Direction == nil
}

func (c *CTransform) updatePosition(durMs int64, movement *pb.SceneMovement, maxX, maxY float32) {
	if isStopMove(movement) {
		return
	}
	secs := float32(durMs) / 1000
	offset := util.Vec2Mul(PbVecToVec2(movement.Direction), movement.MoveSpeed*secs)
	pos := util.Vec2Add(PbVecToVec2(c.Position), offset)
	pos.X = util.ClampFloat(pos.X, 0, maxX)
	pos.Y = util.ClampFloat(pos.Y, 0, maxY)
	Vec2ToPbVec(pos, c.Position)
	c.TransformEventData.Movement = append(c.TransformEventData.Movement, &pb.SceneMovementEvt{
		Position: c.Position,
		Duration: durMs,
	})
}

func (c *CTransform) UpdateTile(tile util.Vec2Int) {
	if c.CurrTile.Equal(tile) {
		return
	}
	c.PrevTile, c.CurrTile = c.CurrTile, tile
	c.tileChanged = true
}

type TileChangeData struct {
	Entry, Exit, Stay []util.Vec2Int
	Origin            *CEvent
}
