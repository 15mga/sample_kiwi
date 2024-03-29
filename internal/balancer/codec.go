// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package balancer

import (
	"game/internal/common"
	"game/proto/pb"

	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

const (
	BalanceRoomReq kiwi.TCode = 0
	BalanceRoomRes kiwi.TCode = 1
)

func (svc *svc) bindCodecFac() {
	kiwi.Codec().BindFac(common.Balancer, BalanceRoomReq, func() util.IMsg {
		return &pb.BalanceRoomReq{}
	})
	kiwi.Codec().BindFac(common.Balancer, BalanceRoomRes, func() util.IMsg {
		return &pb.BalanceRoomRes{}
	})
}
