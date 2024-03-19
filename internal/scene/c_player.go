package scene

import (
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewCPlayer(gateNodeId int64, gateAgentAddr string, data *pb.ScenePlayer, frame *ecs.Frame) *CPlayer {
	return &CPlayer{
		Component:     ecs.NewComponent(C_Player),
		gateNodeId:    gateNodeId,
		gateAgentAddr: gateAgentAddr,
		data:          data,
		frame:         frame,
	}
}

type CPlayer struct {
	ecs.Component
	gateNodeId    int64
	gateAgentAddr string
	data          *pb.ScenePlayer
	frame         *ecs.Frame
}

func (c *CPlayer) Data() util.IMsg {
	return c.data
}

func (c *CPlayer) SetAgentAddr(addr string) {
	c.gateAgentAddr = addr
}

func (c *CPlayer) SetGateNodeId(gateNodeId int64) {
	c.gateNodeId = gateNodeId
}

func (c *CPlayer) ProcessEvents(event *CEvent) {
	//tile := c.Entity().MGetComponent(C_Tile).(*CTile)
	//kiwi.Debug("player process", util.M{
	//	"player":    c.data,
	//	"curr tile": tile.currTile,
	//	"events":    event.eventPus,
	//})
	common.ReqGateNodeToAddr(0, c.gateNodeId, c.gateAgentAddr, event.eventPus,
		func(tid int64, m util.M, code uint16) {

		}, func(tid int64, m util.M, msg util.IMsg) {

		})
	//c.frame.PutJob(JobSendNtc, c.gateNodeId, c.gateAgentAddr, event.Data)
}
