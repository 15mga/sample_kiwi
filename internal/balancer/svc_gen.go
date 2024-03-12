// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package balancer

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
	return common.Balancer
}

func (s *svc) Start() {
	s.bindCodecFac()
	s.registerPusAndReq()
	s.bindReqToRes()
}

func (s *svc) Shutdown() {
}

func (s *svc) Dispose() {
}
