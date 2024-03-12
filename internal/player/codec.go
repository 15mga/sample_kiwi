// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package player

import (
	"game/internal/common"
	"game/proto/pb"

	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

const (
	// 客户端请求自己的信息
	PlayerReq           kiwi.TCode = 0
	PlayerRes           kiwi.TCode = 1
	PlayerManyReq       kiwi.TCode = 2
	PlayerManyRes       kiwi.TCode = 3
	PlayerNewReq        kiwi.TCode = 4
	PlayerNewRes        kiwi.TCode = 5
	PlayerChangeNickReq kiwi.TCode = 6
	PlayerChangeNickRes kiwi.TCode = 7
	//请求其他玩家信息
	PlayerIdReq         kiwi.TCode = 8
	PlayerIdRes         kiwi.TCode = 9
	PlayerReconnectReq  kiwi.TCode = 10
	PlayerReconnectRes  kiwi.TCode = 11
	PlayerListReq       kiwi.TCode = 100
	PlayerListRes       kiwi.TCode = 101
	PlayerDisconnectReq kiwi.TCode = 102
	PlayerDisconnectRes kiwi.TCode = 103
	PlayerDisconnectNtc kiwi.TCode = 104
	PlayerOfflineNtc    kiwi.TCode = 105
	PlayerRepeatNtc     kiwi.TCode = 106
)

func (svc *svc) bindCodecFac() {
	kiwi.Codec().BindFac(common.Player, PlayerReq, func() util.IMsg {
		return &pb.PlayerReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerRes, func() util.IMsg {
		return &pb.PlayerRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerManyReq, func() util.IMsg {
		return &pb.PlayerManyReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerManyRes, func() util.IMsg {
		return &pb.PlayerManyRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerNewReq, func() util.IMsg {
		return &pb.PlayerNewReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerNewRes, func() util.IMsg {
		return &pb.PlayerNewRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerChangeNickReq, func() util.IMsg {
		return &pb.PlayerChangeNickReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerChangeNickRes, func() util.IMsg {
		return &pb.PlayerChangeNickRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerIdReq, func() util.IMsg {
		return &pb.PlayerIdReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerIdRes, func() util.IMsg {
		return &pb.PlayerIdRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerReconnectReq, func() util.IMsg {
		return &pb.PlayerReconnectReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerReconnectRes, func() util.IMsg {
		return &pb.PlayerReconnectRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerListReq, func() util.IMsg {
		return &pb.PlayerListReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerListRes, func() util.IMsg {
		return &pb.PlayerListRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerDisconnectReq, func() util.IMsg {
		return &pb.PlayerDisconnectReq{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerDisconnectRes, func() util.IMsg {
		return &pb.PlayerDisconnectRes{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerDisconnectNtc, func() util.IMsg {
		return &pb.PlayerDisconnectNtc{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerOfflineNtc, func() util.IMsg {
		return &pb.PlayerOfflineNtc{}
	})
	kiwi.Codec().BindFac(common.Player, PlayerRepeatNtc, func() util.IMsg {
		return &pb.PlayerRepeatNtc{}
	})
}
