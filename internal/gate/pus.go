package gate

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
)

func (s *Svc) OnGateCloseId(pkt kiwi.IRcvPush, pus *pb.GateCloseIdPus) {
	kiwi.Gate().CloseWithId(pkt.Tid(), pus.Id, nil, []string{"token"})
}

func (s *Svc) OnGateCloseAddr(pkt kiwi.IRcvPush, pus *pb.GateCloseAddrPus) {
	kiwi.Gate().CloseWithAddr(pkt.Tid(), pus.Addr, nil, []string{"token"})
}
