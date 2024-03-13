// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package scene

import (
	"game/internal/common"
	"game/proto/pb"

	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
)

const (
	SceneEntryReq      kiwi.TCode = 0
	SceneEntryRes      kiwi.TCode = 1
	SceneMovementReq   kiwi.TCode = 2
	SceneMovementRes   kiwi.TCode = 3
	SceneBehaviourReq  kiwi.TCode = 4
	SceneBehaviourRes  kiwi.TCode = 5
	SceneEventPus      kiwi.TCode = 6
	SceneRobotAddReq   kiwi.TCode = 7
	SceneRobotAddRes   kiwi.TCode = 8
	SceneRobotClearReq kiwi.TCode = 9
	SceneRobotClearRes kiwi.TCode = 10
	NewSceneReq        kiwi.TCode = 100
	NewSceneRes        kiwi.TCode = 101
	DisposeSceneReq    kiwi.TCode = 102
	DisposeSceneRes    kiwi.TCode = 103
	SceneGetReq        kiwi.TCode = 104
	SceneGetRes        kiwi.TCode = 105
	SceneHasReq        kiwi.TCode = 106
	SceneHasRes        kiwi.TCode = 107
)

func (svc *svc) bindCodecFac() {
	kiwi.Codec().BindFac(common.Scene, SceneEntryReq, func() util.IMsg {
		return &pb.SceneEntryReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneEntryRes, func() util.IMsg {
		return &pb.SceneEntryRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneMovementReq, func() util.IMsg {
		return &pb.SceneMovementReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneMovementRes, func() util.IMsg {
		return &pb.SceneMovementRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneBehaviourReq, func() util.IMsg {
		return &pb.SceneBehaviourReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneBehaviourRes, func() util.IMsg {
		return &pb.SceneBehaviourRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneEventPus, func() util.IMsg {
		return &pb.SceneEventPus{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneRobotAddReq, func() util.IMsg {
		return &pb.SceneRobotAddReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneRobotAddRes, func() util.IMsg {
		return &pb.SceneRobotAddRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneRobotClearReq, func() util.IMsg {
		return &pb.SceneRobotClearReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneRobotClearRes, func() util.IMsg {
		return &pb.SceneRobotClearRes{}
	})
	kiwi.Codec().BindFac(common.Scene, NewSceneReq, func() util.IMsg {
		return &pb.NewSceneReq{}
	})
	kiwi.Codec().BindFac(common.Scene, NewSceneRes, func() util.IMsg {
		return &pb.NewSceneRes{}
	})
	kiwi.Codec().BindFac(common.Scene, DisposeSceneReq, func() util.IMsg {
		return &pb.DisposeSceneReq{}
	})
	kiwi.Codec().BindFac(common.Scene, DisposeSceneRes, func() util.IMsg {
		return &pb.DisposeSceneRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneGetReq, func() util.IMsg {
		return &pb.SceneGetReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneGetRes, func() util.IMsg {
		return &pb.SceneGetRes{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneHasReq, func() util.IMsg {
		return &pb.SceneHasReq{}
	})
	kiwi.Codec().BindFac(common.Scene, SceneHasRes, func() util.IMsg {
		return &pb.SceneHasRes{}
	})
}