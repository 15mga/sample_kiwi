// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
)

func (svc *svc) watchNtc() {
	kiwi.Router().WatchNotice(&pb.PlayerDisconnectNtc{}, func(ntc kiwi.IRcvNotice) {
		core.SelfPrcNtc[*pb.PlayerDisconnectNtc](ntc, _svc.OnPlayerDisconnectNtc)
	})
}
