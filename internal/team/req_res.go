// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package team

import (
	"game/internal/common"

	"github.com/15mga/kiwi"
)

func (s *svc) bindReqToRes() {
	kiwi.Codec().BindReqToRes(common.Team, NewTeamReq, NewTeamRes)
}
