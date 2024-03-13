package player

import (
	"errors"
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/sid"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/util/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *Svc) OnPlayer(pkt kiwi.IRcvRequest, req *pb.PlayerReq, res *pb.PlayerRes) {
	addr, _ := util.MGet[string](pkt.Head(), common.HdGateAddr)
	opt := options.FindOneAndUpdate().SetProjection(bson.D{
		{Status, 1},
		{LastGateNode, 1},
		{LastAddr, 1},
	})
	uid, _ := util.MGet[string](pkt.Head(), common.HdUserId)
	var player pb.Player
	err := mgo.FindOneAndUpdate(SchemaPlayer, bson.D{
		{UserId, uid},
	}, bson.M{
		"$set": bson.D{
			{Status, pb.PlayerStatus_Online},
			{LastGateNode, pkt.SenderId()},
			{LastAddr, addr},
			{LastSignIn, time.Now().Unix()},
		},
	}, &player, opt)
	if err != nil {
		pkt.Fail(EcPlayer_NotExistId)
		return
	}

	if player.Status == pb.PlayerStatus_Online {
		//给老连接推送重复登录
		common.ReqGateNodeToAddr(pkt.Tid(), player.LastGateNode, player.LastAddr, &pb.PlayerRepeatNtc{},
			func(tid int64, m util.M, code uint16) {

			}, func(tid int64, m util.M, msg util.IMsg) {

			})
		//关闭老链接
		core.AsyncReqNode(pkt.Tid(), player.LastGateNode, nil, &pb.GateCloseAddrReq{
			Addr: player.LastAddr,
		}, func(tid int64, m util.M, code uint16) {

		}, func(tid int64, m util.M, msg util.IMsg) {

		})
	}

	common.ReqGateAddrUpdate(pkt.Tid(), pkt.SenderId(), addr, util.M{
		common.HdPlayerId: player.Id,
	}, nil)

	player.LastAddr = ""
	player.LastGateNode = 0
	res.Player = &player
	pkt.Ok(res)
}

func (s *Svc) OnPlayerMany(pkt kiwi.IRcvRequest, req *pb.PlayerManyReq, res *pb.PlayerManyRes) {
	opt := options.Find()
	if len(req.Projection) > 0 {
		d := bson.D{}
		mgo.BuildProjectionD(&d, nil, req.Projection...)
		opt.SetProjection(d)
	} else {
		opt.SetProjection(bson.D{
			{UserId, -1},
			{Status, -1},
		})
	}
	var players []*pb.Player
	e := mgo.Find(SchemaPlayer, bson.M{
		"$in": bson.D{{Id, req.Id}},
	}, &players, opt)
	if e != nil {
		err := util.WrapErr(util.EcDbErr, e)
		pkt.Err(err)
		kiwi.TE(pkt.Tid(), err)
		return
	}
	res.Players = players
	pkt.Ok(res)
}

func (s *Svc) OnPlayerList(pkt kiwi.IRcvRequest, req *pb.PlayerListReq, res *pb.PlayerListRes) {
	var filter []bson.M
	if req.NameFilter != "" {
		filter = append(filter, bson.M{
			"$text": bson.M{
				"$search": req.NameFilter,
			},
		})
	}
	if req.HeroFilter > 0 {
		filter = append(filter, bson.M{
			Hero: req.HeroFilter,
		})
	}
	opt := options.Find()
	if len(req.Projection) > 0 {
		d := bson.D{}
		mgo.BuildProjectionD(&d, nil, req.Projection...)
		opt.SetProjection(d)
	} else {
		opt.SetProjection(bson.D{
			{UserId, -1},
			{Status, -1},
		})
	}
	mgo.FindWithTotal[pb.Player](SchemaPlayer, bson.M{
		"$and": filter,
	}, func(i int64, players []*pb.Player, e error) {
		if e != nil {
			pkt.Err(util.WrapErr(util.EcDbErr, e))
			return
		}
		res.List = players
		pkt.Ok(res)
	}, opt)
}

func (s *Svc) OnPlayerNew(pkt kiwi.IRcvRequest, req *pb.PlayerNewReq, res *pb.PlayerNewRes) {
	uid, _ := util.MGet[string](pkt.Head(), common.HdUserId)
	player := &pb.Player{
		Id:     sid.GetStrId(),
		UserId: uid,
		Nick:   req.Nick,
		Hero:   req.Hero,
		Status: pb.PlayerStatus_Online,
	}
	_, err := mgo.InsertOne(SchemaPlayer, player)
	if err != nil {
		pkt.Fail(EcPlayerNew_ExistNick)
		return
	}
	res.Player = player

	addr, _ := util.MGet[string](pkt.Head(), common.HdGateAddr)
	common.ReqGateAddrUpdate(pkt.Tid(), pkt.SenderId(), addr, util.M{
		common.HdPlayerId: player.Id,
	}, nil)

	pkt.Ok(res)
}

func (s *Svc) OnPlayerId(pkt kiwi.IRcvRequest, req *pb.PlayerIdReq, res *pb.PlayerIdRes) {
	opt := options.FindOne()
	if len(req.Projection) > 0 {
		d := bson.D{}
		mgo.BuildProjectionD(&d, nil, req.Projection...)
		opt.SetProjection(d)
	} else {
		opt.SetProjection(bson.D{
			{UserId, 0},
			{Status, 0},
		})
	}
	var player pb.Player
	err := mgo.FindOne(SchemaPlayer, bson.D{
		{Id, req.Id},
	}, &player, opt)
	if err != nil {
		pkt.Fail(EcPlayer_NotExistId)
		return
	}
	res.Player = &player
	pkt.Ok(res)
}

func (s *Svc) OnPlayerReconnect(pkt kiwi.IRcvRequest, req *pb.PlayerReconnectReq, res *pb.PlayerReconnectRes) {
	filter := bson.M{}
	if req.Token != "" {
		filter[ReconnectToken] = req.Token
	} else {
		filter[ReconnectId] = pkt.HeadId()
	}
	var reconnect pb.Reconnect
	e := mgo.FindOneAndDel(SchemaReconnect, filter, &reconnect)
	if errors.Is(e, mongo.ErrNoDocuments) {
		pkt.Fail(EcPlayerReconnect_NotExist)
		return
	}
	addr, _ := util.MGet[string](pkt.Head(), common.HdGateAddr)
	var (
		head  = util.M{}
		cache = util.M{}
	)
	err := head.FromBytes(reconnect.Head)
	if err != nil {
		pkt.Err(err)
		return
	}
	err = cache.FromBytes(reconnect.Cache)
	if err != nil {
		pkt.Err(err)
		return
	}
	common.ReqGateAddrUpdate(pkt.Tid(), pkt.SenderId(), addr, head, cache)
	res.RoomId, _ = util.MGet[string](head, common.HdRoomId)
	sceneId, sceneExist := util.MGet[string](head, common.HdSceneId)
	if !sceneExist {
		pkt.Ok(res)
		return
	}
	//请求场景是否存在
	core.AsyncSubReq[*pb.SceneHasRes](pkt, nil,
		&pb.SceneHasReq{
			Id: sceneId,
		}, func(tid int64, m util.M, code uint16) {
			pkt.Ok(res)
		}, func(tid int64, m util.M, hasRes *pb.SceneHasRes) {
			if hasRes.Exist {
				res.SceneId = sceneId
			}
			pkt.Ok(res)
		})
}

func (s *Svc) OnPlayerDisconnect(pkt kiwi.IRcvRequest, req *pb.PlayerDisconnectReq, res *pb.PlayerDisconnectRes) {
	id, ok := util.MGet[string](pkt.Head(), common.HdPlayerId)
	if !ok {
		pkt.Ok(res)
		return
	}
	cacheBytes := req.Cache
	cache := util.M{}
	err := cache.FromBytes(cacheBytes)
	if err != nil {
		pkt.Err(err)
		return
	}
	token, ok := util.MGet[string](cache, common.CcToken)
	if !ok {
		pkt.Ok(res)
		return
	}
	_, e := mgo.UpdateOne(SchemaPlayer, bson.M{
		Id:     id,
		Status: pb.PlayerStatus_Online,
	}, bson.M{
		"$set": bson.M{
			Status: pb.PlayerStatus_Disconnect,
		},
	})
	if e != nil {
		pkt.Err3(util.EcDbErr, e)
		return
	}
	headBytes, _ := pkt.Head().ToBytes()
	_, e = mgo.InsertOne(SchemaReconnect, bson.M{
		ReconnectId:        id,
		ReconnectToken:     token,
		ReconnectTimestamp: time.Now().Format(time.RFC3339),
		ReconnectHead:      headBytes,
		ReconnectCache:     cacheBytes,
	})
	if e != nil {
		pkt.Err3(util.EcDbErr, e)
		return
	}

	_ = core.Ntf(pkt.Tid(), pkt.Head(), &pb.PlayerDisconnectNtc{
		Cache: cacheBytes,
	})
	pkt.Ok(res)

	//时间注意和索引一致
	_IdToOfflineTimer.Set(id, time.AfterFunc(time.Second*5, func() {
		core.ActivePrcReq[*pb.PlayerDisconnectReq, *pb.PlayerDisconnectRes](pkt, id,
			func(pkt kiwi.IRcvRequest, req *pb.PlayerDisconnectReq, res *pb.PlayerDisconnectRes) {
				r, e := mgo.DelOne(SchemaReconnect, bson.M{
					ReconnectId: id,
				})
				if e != nil || r.DeletedCount == 0 {
					return
				}
				updateResult, e := mgo.UpdateOne(SchemaPlayer, bson.M{
					Id:     id,
					Status: pb.PlayerStatus_Disconnect,
				}, bson.M{
					"$set": bson.M{
						Status: pb.PlayerStatus_Offline,
					},
				})
				if e != nil {
					kiwi.TE3(pkt.Tid(), util.EcDbErr, e)
					return
				}
				if updateResult.ModifiedCount == 0 {
					return
				}
				_ = core.Ntf(pkt.Tid(), pkt.Head(), &pb.PlayerOfflineNtc{
					Id: id,
				})
			})
	}))
}
