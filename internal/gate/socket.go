package gate

import (
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/util"
	"strconv"
)

// CheckIp 检查是否黑名单
func CheckIp(ip string) bool {
	return true
}

// RecordSocketErr 记录错误
func RecordSocketErr(agent kiwi.IAgent, err *util.Err) {
	head := util.M{}
	cache := util.M{}
	agent.CopyHead(head)
	agent.CopyCache(cache)
	err.AddParam("head", head)
	err.AddParam("cache", cache)
	kiwi.Error(err)
}

func Disconnected(agent kiwi.IAgent, _ *util.Err) {
	cache := util.M{}
	agent.CopyCache(cache)
	_, ok := util.MGet[string](cache, common.CcToken)
	if !ok { //没有token,没登录或者服务端主动断开的链接
		return
	}
	cacheBytes, err := cache.ToBytes()
	if err != nil {
		kiwi.Error(err)
		return
	}
	head := util.M{}
	agent.CopyHead(head)
	core.AsyncReq(0, head, &pb.PlayerDisconnectReq{
		Cache: cacheBytes,
	}, func(tid int64, m util.M, code uint16) {

	}, func(tid int64, m util.M, msg util.IMsg) {

	})
}

func InitAgent(agent kiwi.IAgent) {
	agent.SetHead(common.HdMask, util.GenMask(common.RGuest))
}

func SocketReceiver(agent kiwi.IAgent, bytes []byte) {
	svc, code, seqId, payload, err := common.UnpackUserReq(bytes)
	if err != nil {
		RecordSocketErr(agent, err)
		return
	}
	roleMask, _ := agent.GetHead(common.HdMask)
	mask := roleMask.(int64)
	ok := kiwi.Gate().Authenticate(mask, svc, code)
	if !ok {
		RecordSocketErr(agent, util.NewErr(util.EcNoAuth, util.M{
			"service": svc,
			"code":    code,
		}))
		return
	}

	head := util.M{
		"seq_id": seqId,
	}
	agent.CopyHead(head)
	nodeId, ok := agent.GetCache(strconv.Itoa(int(svc)))
	if ok {
		core.AsyncReqNodeBytes(0, nodeId.(int64), svc, code, head, false, payload, onResErr, onResOk)
	} else {
		core.AsyncReqBytes(0, svc, code, head, false, payload, onResErr, onResOk)
	}
}

func onResOk(tid int64, head util.M, payload []byte) {
	addr, _ := util.MGet[string](head, "addr")
	seqId, _ := util.MGet[uint32](head, "seq_id")
	pkt, e := common.PackUserOk(seqId, payload)
	if e != nil {
		kiwi.Error(e)
		return
	}
	kiwi.Gate().AddrSend(tid, addr, pkt, func(ok bool) {
		if !ok {
			kiwi.TD(tid, "not exist", util.M{
				"addr": addr,
			})
		}
	})
}

func onResErr(tid int64, head util.M, resCode uint16) {
	addr, _ := util.MGet[string](head, "addr")
	seqId, _ := util.MGet[uint32](head, "seq_id")
	pkt, e := common.PackUserFail(seqId, resCode)
	if e != nil {
		kiwi.Error(e)
		return
	}
	kiwi.Gate().AddrSend(tid, addr, pkt, func(ok bool) {
		if !ok {
			kiwi.TD(tid, "not exist", util.M{
				"addr": addr,
			})
		}
	})
}
