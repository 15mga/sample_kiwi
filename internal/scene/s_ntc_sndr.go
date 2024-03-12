package scene

import (
	"game/internal/common"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
)

func NewSNtcSender() *SNtcSender {
	return &SNtcSender{
		System: ecs.NewSystem(S_Ntc_Sender),
	}
}

type SNtcSender struct {
	ecs.System
	nodeIdToNtc map[int64]map[string]util.IMsg
}

func (s *SNtcSender) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindJob(JobSendNtc, s.onSendNtc)
}

func (s *SNtcSender) OnUpdate() {
	s.nodeIdToNtc = make(map[int64]map[string]util.IMsg)
	s.DoJob(JobSendNtc)
	for nodeId, m := range s.nodeIdToNtc {
		common.ReqGateToMultiAddrMap(0, nodeId, m)
	}
}

func (s *SNtcSender) onSendNtc(data []any) {
	nodeId, addr, msg := util.SplitSlc3[int64, string, util.IMsg](data)
	m, ok := s.nodeIdToNtc[nodeId]
	if !ok {
		m = make(map[string]util.IMsg, 32)
		s.nodeIdToNtc[nodeId] = m
	}
	m[addr] = msg
}
