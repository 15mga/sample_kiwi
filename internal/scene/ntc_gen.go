// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

func (s *svc) OnPlayerDisconnectNtc(pkt kiwi.IRcvNotice, ntc *pb.PlayerDisconnectNtc) {
	pkt.Err2(util.EcNotImplement, util.M{"ntc": ntc})
}