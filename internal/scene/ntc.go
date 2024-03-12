package scene

import (
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

func (s *Svc) OnPlayerDisconnectNtc(pkt kiwi.IRcvNotice, ntc *pb.PlayerDisconnectNtc) {
	defer pkt.Complete()
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdSceneId)
	if !ok {
		return
	}
	PushJob(sceneId, JobEntityDel, pkt.Tid(), []string{pkt.HeadId()})
}
