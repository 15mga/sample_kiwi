// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package user

import (
	"game/internal/common"
	"github.com/15mga/kiwi"
)

var (
	_svc = &Svc{}
)

func Service() *Svc {
	return _svc
}

type Svc struct {
	svc
}

type svc struct {
}

func (s *svc) Svc() kiwi.TSvc {
	return common.User
}

func (s *svc) Start() {
	s.bindCodecFac()
	s.initColl()
	s.registerReq()
	s.bindReqToRes()
}

func (s *svc) Shutdown() {
}

func (s *svc) Dispose() {
}
