// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package gate

import (
	"game/internal/common"

	"github.com/15mga/kiwi"
)

func (s *svc) bindReqToRes() {
	kiwi.Codec().BindReqToRes(common.Gate, HeartbeatReq, HeartbeatRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateBanAddrReq, GateBanAddrRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateSendToIdReq, GateSendToIdRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateSendToAddrReq, GateSendToAddrRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateSendToMultiIdReq, GateSendToMultiIdRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateSendToMultiAddrReq, GateSendToMultiAddrRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateSendToAllReq, GateSendToAllRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateUpdateReq, GateUpdateRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateAddrUpdateReq, GateAddrUpdateRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateRemoveReq, GateRemoveRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateAddrRemoveReq, GateAddrRemoveRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateGetReq, GateGetRes)
	kiwi.Codec().BindReqToRes(common.Gate, GateAddrGetReq, GateAddrGetRes)
}
