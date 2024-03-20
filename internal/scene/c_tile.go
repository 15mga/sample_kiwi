package scene

import (
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewCTile() *CTile {
	c := &CTile{
		Component: ecs.NewComponent(C_Tile),
		interest:  ds.NewArray[util.Vec2Int](9),
		entry:     ds.NewArray[util.Vec2Int](9),
		exit:      ds.NewArray[util.Vec2Int](9),
		stay:      ds.NewArray[util.Vec2Int](9),
	}
	return c
}

type CTile struct {
	ecs.Component
	prevTile          util.Vec2Int
	currTile          util.Vec2Int
	cEvent            *CEvent
	cTnf              *CTransform
	tileChanged       bool
	interest          *ds.Array[util.Vec2Int]
	state             TileState
	entry, exit, stay *ds.Array[util.Vec2Int]
}

func (c *CTile) Start() {
	c.cEvent = c.Entity().MGetComponent(C_Event).(*CEvent)
	c.cTnf = c.Entity().MGetComponent(C_Transform).(*CTransform)
}

func (c *CTile) InitTile(tile util.Vec2Int) {
	c.currTile = tile
	c.prevTile = tile
}

func (c *CTile) ResetTileChanged() {
	c.tileChanged = false
}

func (c *CTile) IsTileChanged() bool {
	return c.tileChanged
}

func (c *CTile) UpdateTile(tile util.Vec2Int) {
	if c.currTile.Equal(tile) {
		return
	}
	c.prevTile, c.currTile = c.currTile, tile
	c.tileChanged = true
}
