// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package room

import (
	"game/internal/common"

	"github.com/15mga/kiwi"
)

func (s *svc) bindReqToRes() {
	kiwi.Codec().BindReqToRes(common.Room, RoomNewReq, RoomNewRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomListReq, RoomListRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomEntryReq, RoomEntryRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomExitReq, RoomExitRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomReadyReq, RoomReadyRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomStartReq, RoomStartRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomGetReq, RoomGetRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomModifyReq, RoomModifyRes)
	kiwi.Codec().BindReqToRes(common.Room, RoomReconnectReq, RoomReconnectRes)
}
