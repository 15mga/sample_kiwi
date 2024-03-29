// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package chat

import (
	"game/internal/common"
	"game/proto/pb"

	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

const (
	NewMsgReq  kiwi.TCode = 0
	NewMsgRes  kiwi.TCode = 1
	MsgPus     kiwi.TCode = 2
	NewChanReq kiwi.TCode = 100
	NewChanRes kiwi.TCode = 101
)

func (svc *svc) bindCodecFac() {
	kiwi.Codec().BindFac(common.Chat, NewMsgReq, func() util.IMsg {
		return &pb.NewMsgReq{}
	})
	kiwi.Codec().BindFac(common.Chat, NewMsgRes, func() util.IMsg {
		return &pb.NewMsgRes{}
	})
	kiwi.Codec().BindFac(common.Chat, MsgPus, func() util.IMsg {
		return &pb.MsgPus{}
	})
	kiwi.Codec().BindFac(common.Chat, NewChanReq, func() util.IMsg {
		return &pb.NewChanReq{}
	})
	kiwi.Codec().BindFac(common.Chat, NewChanRes, func() util.IMsg {
		return &pb.NewChanRes{}
	})
}
