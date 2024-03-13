package common

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/util"
)

func ReqGateUpdate(tid, nodeId int64, id string, head, cache util.M) {
	var (
		headBytes, cacheBytes []byte
		err                   *util.Err
	)
	if head != nil {
		headBytes, err = kiwi.Packer().PackM(head)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
	}
	if cache != nil {
		cacheBytes, err = kiwi.Packer().PackM(cache)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
	}
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateUpdateReq{
		Id:    id,
		Head:  headBytes,
		Cache: cacheBytes,
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateAddrUpdate(tid, nodeId int64, addr string, head, cache util.M) {
	headBytes, err := kiwi.Packer().PackM(head)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	cacheBytes, err := kiwi.Packer().PackM(cache)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateAddrUpdateReq{
		Addr:  addr,
		Head:  headBytes,
		Cache: cacheBytes,
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateRemove(tid, nodeId int64, id string, head, cache []string) {
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateRemoveReq{
		Id:    id,
		Head:  head,
		Cache: cache,
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateRemoveAddr(tid, nodeId int64, addr string, head, cache []string) {
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateAddrRemoveReq{
		Addr:  addr,
		Head:  head,
		Cache: cache,
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateToId(tid int64, id string, msg util.IMsg) {
	bytes, err := kiwi.Codec().PbMarshal(msg)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	sc := int32(kiwi.MergeSvcCode(kiwi.Codec().MsgToSvcCode(msg)))
	core.AsyncReq(tid, nil, &pb.GateSendToIdReq{
		Id:      id,
		SvcCode: sc,
		Payload: bytes,
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateNodeToId(tid, nodeId int64, id string, msg util.IMsg,
	onFail util.FnInt64MUint16, onOk util.FnInt64MMsg) {
	bytes, err := kiwi.Codec().PbMarshal(msg)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	sc := int32(kiwi.MergeSvcCode(kiwi.Codec().MsgToSvcCode(msg)))
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateSendToIdReq{
		Id:      id,
		SvcCode: sc,
		Payload: bytes,
	}, onFail, onOk)
}

func ReqGateNodeToAddr(tid, nodeId int64, addr string, msg util.IMsg,
	onFail util.FnInt64MUint16, onOk util.FnInt64MMsg) {
	bytes, err := kiwi.Codec().PbMarshal(msg)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	sc := int32(kiwi.MergeSvcCode(kiwi.Codec().MsgToSvcCode(msg)))
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateSendToAddrReq{
		Addr:    addr,
		SvcCode: sc,
		Payload: bytes,
	}, onFail, onOk)
}

func ReqGateToMultiId(tid int64, id []string, msg util.IMsg) {
	if len(id) == 0 {
		return
	}
	bytes, err := kiwi.Codec().PbMarshal(msg)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	sc := int32(kiwi.MergeSvcCode(kiwi.Codec().MsgToSvcCode(msg)))
	for _, i := range id {
		core.AsyncReq(tid, nil, &pb.GateSendToIdReq{
			Id:      i,
			SvcCode: sc,
			Payload: bytes,
		}, func(tid int64, m util.M, code uint16) {
			kiwi.TE2(tid, code, nil)
		}, func(tid int64, m util.M, msg util.IMsg) {

		})
	}
}

func ReqGateToMultiAddr(tid int64, addr []string, msg util.IMsg) {
	if len(addr) == 0 {
		return
	}
	bytes, err := kiwi.Codec().PbMarshal(msg)
	if err != nil {
		kiwi.TE(tid, err)
		return
	}
	sc := int32(kiwi.MergeSvcCode(kiwi.Codec().MsgToSvcCode(msg)))
	for _, a := range addr {
		core.AsyncReq(tid, nil, &pb.GateSendToAddrReq{
			Addr:    a,
			SvcCode: sc,
			Payload: bytes,
		}, func(tid int64, m util.M, code uint16) {
			kiwi.TE2(tid, code, nil)
		}, func(tid int64, m util.M, msg util.IMsg) {

		})
	}
}

func ReqGateToMultiIdMap(tid, nodeId int64, idToMsg map[string]util.IMsg) {
	l := len(idToMsg)
	if l == 0 {
		return
	}
	var buffer util.ByteBuffer
	buffer.InitCap(1024)
	buffer.WUint16(uint16(l))
	for id, msg := range idToMsg {
		bytes, err := kiwi.Codec().PbMarshal(msg)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
		buffer.WString(id)
		svc, code := kiwi.Codec().MsgToSvcCode(msg)
		payload, err := PackUserPus(svc, code, bytes)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
		buffer.WBytes(payload)
	}
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateSendToMultiIdReq{
		Payload: buffer.All(),
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func ReqGateToMultiAddrMap(tid, nodeId int64, addrToMsg map[string]util.IMsg) {
	l := len(addrToMsg)
	if l == 0 {
		return
	}
	var buffer util.ByteBuffer
	buffer.InitCap(1024)
	buffer.WUint16(uint16(l))
	for addr, msg := range addrToMsg {
		bytes, err := kiwi.Codec().PbMarshal(msg)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
		buffer.WString(addr)
		svc, code := kiwi.Codec().MsgToSvcCode(msg)
		payload, err := PackUserPus(svc, code, bytes)
		if err != nil {
			kiwi.TE(tid, err)
			return
		}
		buffer.WBytes(payload)
	}
	core.AsyncReqNode(tid, nodeId, nil, &pb.GateSendToMultiIdReq{
		Payload: buffer.All(),
	}, func(tid int64, m util.M, code uint16) {
		kiwi.TE2(tid, code, nil)
	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}
