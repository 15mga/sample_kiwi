// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package user

import (
	"game/internal/common"

	"github.com/15mga/kiwi"
)

func (s *svc) bindReqToRes() {
	kiwi.Codec().BindReqToRes(common.User, SignUpReq, SignUpRes)
	kiwi.Codec().BindReqToRes(common.User, SignInReq, SignInRes)
	kiwi.Codec().BindReqToRes(common.User, SignOutReq, SignOutRes)
}