package scene

import (
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/util"
	"time"
)

const (
	_SceneData = "scene_data"
	_MaxCount  = 1024 << 8
)

var (
	_Frames = ds.NewKSet[string, *ecs.Frame](2, func(frame *ecs.Frame) string {
		return frame.Scene().Id()
	})
)

type Conf struct {
	TileSize int
	Width    int
	Height   int
	FovLaps  int
	Pawns    []*pb.SceneVisible
}

func getConf(scene *pb.Scene) *Conf {
	return &Conf{
		TileSize: 64,
		Width:    2048,
		Height:   2048,
		FovLaps:  1,
	}
}

func NewScene(tid int64, scn *pb.Scene) *util.Err {
	_, ok := _Frames.Get(scn.Id)
	if ok {
		return util.NewErr(util.EcExist, util.M{
			"scene id": scn.Id,
		})
	}
	conf := getConf(scn)
	scene := ecs.NewScene(scn.Id, ecs.TScene(scn.TplId))
	scene.Data().Set(_SceneData, scn)
	frame := ecs.NewFrame(scene, ecs.FrameSystems(
		NewSEntity(),
		NewSRobot(common.Conf().Test.MaxRobot),
		NewSTransform(conf.Width, conf.Height),
		NewSTile(conf.TileSize, conf.Width, conf.Height, conf.FovLaps),
		NewSEvent(),
		//NewSNtcSender(),
	), ecs.FrameTickDur(time.Millisecond*50))
	frame.Start()
	_ = _Frames.Add(frame)
	kiwi.TI(tid, "new scene", util.M{
		"data": scn,
	})
	return nil
}

func DisposeScene(sceneId string) {
	frame, ok := _Frames.Del(sceneId)
	if !ok {
		return
	}
	frame.Stop()
}

func GetSceneDataById(sceneId string) (*pb.Scene, bool) {
	frame, ok := _Frames.Get(sceneId)
	if !ok {
		return nil, false
	}
	return GetSceneDataByFrame(frame), true
}

func HasScene(sceneId string) bool {
	return _Frames.Has(sceneId)
}

func PushJob(sceneId, job string, params ...any) bool {
	frame, ok := _Frames.Get(sceneId)
	if !ok {
		return false
	}
	frame.PushJob(job, params...)
	return true
}

func GetSceneDataByFrame(frame *ecs.Frame) *pb.Scene {
	scene, _ := util.MGet[*pb.Scene](frame.Scene().Data(), _SceneData)
	return scene
}
