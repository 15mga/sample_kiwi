// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package chat

import (
	"game/internal/common"

	"github.com/15mga/kiwi"
)

func (s *svc) bindReqToRes() {
	kiwi.Codec().BindReqToRes(common.Chat, NewMsgReq, NewMsgRes)
	kiwi.Codec().BindReqToRes(common.Chat, NewChanReq, NewChanRes)
}