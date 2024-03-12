package scene

import (
	"game/internal/common"
	"game/internal/player"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/util"
	"strconv"
)

func (s *Svc) OnSceneEntry(pkt kiwi.IRcvRequest, req *pb.SceneEntryReq, res *pb.SceneEntryRes) {
	playerId := pkt.HeadId()
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdRoomId) //场景id使用的房间id
	if !ok {
		pkt.Fail(EcSceneBehaviour_NotEntry)
		return
	}
	core.AsyncSubReq[*pb.PlayerIdRes](pkt, pkt.Head(), &pb.PlayerIdReq{
		Id: playerId,
		Projection: []string{
			player.Id,
			player.TeamId,
			player.Nick,
			player.Hero,
			player.Avatar,
			player.Gender,
			player.LastAddr,
			player.LastGateNode,
		},
	}, func(tid int64, m util.M, code uint16) {
		pkt.Fail(code)
	}, func(tid int64, m util.M, playerRes *pb.PlayerIdRes) {
		p := playerRes.Player
		addr, _ := util.MGet[string](pkt.Head(), common.HdGateAddr)
		nodeId, _ := util.MGet[int64](pkt.Head(), common.HdGateNodeId)
		position := s.getPlayerPosition(p)
		pawn := &pb.ScenePawn{
			Pawn: &pb.ScenePawnEvt{PawnId: p.Id,
				Position: position,
				PawnType: &pb.ScenePawnEvt_Player{
					Player: &pb.ScenePlayer{
						Nick:   p.Nick,
						Avatar: p.Avatar,
						TeamId: p.TeamId,
						Gender: p.Gender,
					},
				},
			},
			PlayerGateAddr:   addr,
			PlayerGateNodeId: nodeId,
		}
		ok := PushJob(sceneId, JobEntityAdd, tid, pawn,
			func(err *util.Err) {
				if err != nil {
					pkt.Err(err)
					return
				}
				res.Position = position
				pkt.Ok(res)
				pr := playerRes.Player
				common.ReqGateAddrUpdate(pkt.Tid(), pr.LastGateNode,
					pr.LastAddr, util.M{
						common.HdSceneId: sceneId, //场景id
					}, util.M{
						strconv.Itoa(int(common.Scene)): kiwi.GetNodeMeta().NodeId, //场景服务的节点id
					})
			})
		if !ok {
			pkt.Fail(EcSceneEntry_NotExistScene)
		}
	})
}

func (s *Svc) getPlayerPosition(player *pb.Player) *pb.Vector2 {
	return &pb.Vector2{
		X: 100,
		Y: 100,
	}
}

func (s *Svc) OnSceneMovement(pkt kiwi.IRcvRequest, req *pb.SceneMovementReq, res *pb.SceneMovementRes) {
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdSceneId)
	if !ok {
		pkt.Fail(EcSceneMovement_NotEntry)
		return
	}
	playerId := pkt.HeadId()
	PushJob(sceneId, JobMovement, pkt.Tid(), playerId, req.Movement)
	pkt.Ok(res)
}

func (s *Svc) OnSceneGet(pkt kiwi.IRcvRequest, req *pb.SceneGetReq, res *pb.SceneGetRes) {
	scene, ok := GetSceneDataById(req.Id)
	if !ok {
		pkt.Fail(EcSceneGet_NotExistScene)
		return
	}
	res.Scene = scene
	pkt.Ok(res)
}

func (s *Svc) OnSceneHas(pkt kiwi.IRcvRequest, req *pb.SceneHasReq, res *pb.SceneHasRes) {
	res.Exist = HasScene(req.Id)
	pkt.Ok(res)
}

func (s *Svc) OnSceneBehaviour(pkt kiwi.IRcvRequest, req *pb.SceneBehaviourReq, res *pb.SceneBehaviourRes) {
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdSceneId)
	if !ok {
		pkt.Fail(EcSceneBehaviour_NotEntry)
		return
	}
	playerId := pkt.HeadId()
	PushJob(sceneId, JobBehaviour, playerId, req.Behaviour, func(code uint16) {
		if code > 0 {
			pkt.Fail(code)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnSceneRobotAdd(pkt kiwi.IRcvRequest, req *pb.SceneRobotAddReq, res *pb.SceneRobotAddRes) {
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdSceneId)
	if !ok {
		pkt.Fail(EcSceneRobotAdd_NotEntry)
		return
	}
	PushJob(sceneId, JobRobotAdd, pkt.Tid(), req.Count, func(code uint16) {
		if code > 0 {
			pkt.Fail(code)
			return
		}
		pkt.Ok(res)
	})
}

func (s *Svc) OnSceneRobotClear(pkt kiwi.IRcvRequest, req *pb.SceneRobotClearReq, res *pb.SceneRobotClearRes) {
	sceneId, ok := util.MGet[string](pkt.Head(), common.HdSceneId)
	if !ok {
		pkt.Fail(EcSceneRobotClear_NotEntry)
		return
	}
	PushJob(sceneId, JobRobotClear, pkt.Tid())
	pkt.Ok(res)
}

func (s *Svc) OnNewScene(pkt kiwi.IRcvRequest, req *pb.NewSceneReq, res *pb.NewSceneRes) {
	err := NewScene(pkt.Tid(), req.Scene)
	if err != nil {
		pkt.Err(err)
		return
	}
	pkt.Ok(res)
}

func (s *Svc) OnDisposeScene(pkt kiwi.IRcvRequest, req *pb.DisposeSceneReq, res *pb.DisposeSceneRes) {
	err := DisposeScene(req.SceneId)
	if err != nil {
		pkt.Err(err)
		return
	}
	pkt.Ok(res)
}
