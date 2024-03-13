package gate

import (
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
	"time"
)

func (s *Svc) OnHeartbeat(pkt kiwi.IRcvRequest, req *pb.HeartbeatReq, res *pb.HeartbeatRes) {
	res.ReqTs = req.ReqTs
	res.ResTs = time.Now().UnixMilli()
	pkt.Ok(res)
}

func (s *Svc) OnGateSendToId(pkt kiwi.IRcvRequest, req *pb.GateSendToIdReq, res *pb.GateSendToIdRes) {
	svc, code := kiwi.SplitSvcCode(uint16(req.SvcCode))
	bytes, err := common.PackUserPus(svc, code, req.Payload)
	if err != nil {
		kiwi.TE(pkt.Tid(), err)
		return
	}
	kiwi.Gate().Send(pkt.Tid(), req.Id, bytes, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateSendToId_NotExistId)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateSendToAddr(pkt kiwi.IRcvRequest, req *pb.GateSendToAddrReq, res *pb.GateSendToAddrRes) {
	svc, code := kiwi.SplitSvcCode(uint16(req.SvcCode))
	bytes, err := common.PackUserPus(svc, code, req.Payload)
	if err != nil {
		kiwi.TE(pkt.Tid(), err)
		return
	}
	kiwi.Gate().AddrSend(pkt.Tid(), req.Addr, bytes, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateSendToAddr_NotExistAddr)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateSendToMultiId(pkt kiwi.IRcvRequest, req *pb.GateSendToMultiIdReq, res *pb.GateSendToMultiIdRes) {
	payload := req.Payload
	var buffer util.ByteBuffer
	buffer.InitBytes(payload)
	count, err := buffer.RUint16()
	if err != nil {
		kiwi.TE(pkt.Tid(), err)
		return
	}
	idToMsg := make(map[string][]byte, count)
	for i := uint16(0); i < count; i++ {
		id, err := buffer.RString()
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
		p, err := buffer.RBytes()
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
		idToMsg[id] = p
	}
	kiwi.Gate().MultiSend(pkt.Tid(), idToMsg, func(m map[string]bool) {
		res.Result = m
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateSendToMultiAddr(pkt kiwi.IRcvRequest, req *pb.GateSendToMultiAddrReq, res *pb.GateSendToMultiAddrRes) {
	payload := req.Payload
	var buffer util.ByteBuffer
	buffer.InitBytes(payload)
	count, err := buffer.RUint16()
	if err != nil {
		kiwi.TE(pkt.Tid(), err)
		return
	}
	addrToMsg := make(map[string][]byte, count)
	for i := uint16(0); i < count; i++ {
		id, err := buffer.RString()
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
		p, err := buffer.RBytes()
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
		addrToMsg[id] = p
	}
	kiwi.Gate().MultiAddrSend(pkt.Tid(), addrToMsg, func(m map[string]bool) {
		res.Result = m
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateSendToAll(pkt kiwi.IRcvRequest, req *pb.GateSendToAllReq, res *pb.GateSendToAllRes) {
	svc, code := kiwi.SplitSvcCode(uint16(req.SvcCode))
	bytes, err := common.PackUserPus(svc, code, req.Payload)
	if err != nil {
		kiwi.TE(pkt.Tid(), err)
		return
	}
	kiwi.Gate().AllSend(pkt.Tid(), bytes)
}

func (s *Svc) OnGateCloseId(pkt kiwi.IRcvRequest, req *pb.GateCloseIdReq, res *pb.GateCloseIdRes) {
	kiwi.Gate().CloseWithId(pkt.Tid(), req.Id, nil, []string{"token"})
	pkt.Ok(res)
}

func (s *Svc) OnGateCloseAddr(pkt kiwi.IRcvRequest, req *pb.GateCloseAddrReq, res *pb.GateCloseAddrRes) {
	kiwi.Gate().CloseWithAddr(pkt.Tid(), req.Addr, nil, []string{"token"})
	pkt.Ok(res)
}

func (s *Svc) OnGateUpdate(pkt kiwi.IRcvRequest, req *pb.GateUpdateReq, res *pb.GateUpdateRes) {
	var (
		head, cache util.M
	)
	if req.Head != nil {
		head = make(util.M)
		err := kiwi.Packer().UnpackM(req.Head, head)
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
	}
	if req.Cache != nil {
		cache = make(util.M)
		err := kiwi.Packer().UnpackM(req.Cache, cache)
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
	}
	kiwi.Gate().UpdateHeadCache(pkt.Tid(), req.Id, head, cache, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateUpdate_NotExistId)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateAddrUpdate(pkt kiwi.IRcvRequest, req *pb.GateAddrUpdateReq, res *pb.GateAddrUpdateRes) {
	var (
		head, cache util.M
	)
	if req.Head != nil {
		head = make(util.M)
		err := kiwi.Packer().UnpackM(req.Head, head)
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
	}
	if req.Cache != nil {
		cache = make(util.M)
		err := kiwi.Packer().UnpackM(req.Cache, cache)
		if err != nil {
			kiwi.TE(pkt.Tid(), err)
			return
		}
	}
	kiwi.Gate().UpdateAddrHeadCache(pkt.Tid(), req.Addr, head, cache, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateAddrUpdate_NotExistAddr)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateRemove(pkt kiwi.IRcvRequest, req *pb.GateRemoveReq, res *pb.GateRemoveRes) {
	kiwi.Gate().RemoveHeadCache(pkt.Tid(), req.Id, req.Head, req.Cache, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateRemove_NotExistId)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateAddrRemove(pkt kiwi.IRcvRequest, req *pb.GateAddrRemoveReq, res *pb.GateAddrRemoveRes) {
	kiwi.Gate().RemoveAddrHeadCache(pkt.Tid(), req.Addr, req.Head, req.Cache, func(ok bool) {
		if !ok {
			pkt.Fail(EcGateAddrRemove_NotExistAddr)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateGet(pkt kiwi.IRcvRequest, req *pb.GateGetReq, res *pb.GateGetRes) {
	kiwi.Gate().GetHeadCache(pkt.Tid(), req.Id, func(head util.M, cache util.M, ok bool) {
		if !ok {
			pkt.Fail(EcGateGet_NotExistId)
			return
		}
		res.Head, _ = kiwi.Packer().PackM(head)
		res.Cache, _ = kiwi.Packer().PackM(cache)
		pkt.Ok(res)
	})
}

func (s *Svc) OnGateAddrGet(pkt kiwi.IRcvRequest, req *pb.GateAddrGetReq, res *pb.GateAddrGetRes) {
	kiwi.Gate().GetAddrHeadCache(pkt.Tid(), req.Addr, func(head util.M, cache util.M, ok bool) {
		if !ok {
			pkt.Fail(EcGateAddrGet_NotExistAddr)
			return
		}
		res.Head, _ = kiwi.Packer().PackM(head)
		res.Cache, _ = kiwi.Packer().PackM(cache)
		pkt.Ok(res)
	})
}
