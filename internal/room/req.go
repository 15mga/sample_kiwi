package room

import (
	"errors"
	"fmt"
	"game/internal/common"
	"game/internal/player"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/sid"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/util/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"strconv"
	"time"
)

func (s *Svc) OnRoomNew(pkt kiwi.IRcvRequest, req *pb.RoomNewReq, res *pb.RoomNewRes) {
	roomId := sid.GetStrId()
	room := &pb.Room{
		Id:          roomId,
		Name:        req.Name,
		Mode:        req.Mode,
		SceneTplId:  req.SceneTplId,
		Password:    req.Password,
		OwnerId:     pkt.HeadId(),
		Status:      pb.RoomStatus_InLobby,
		MaxPlayers:  req.MaxPlayers,
		CurrPlayers: 1,
		CreateTime:  time.Now().Unix(),
	}
	players := make(map[string]*pb.RoomPlayer, 8)
	players[pkt.HeadId()] = &pb.RoomPlayer{
		Ready:      false,
		Disconnect: false,
		EnterTs:    time.Now().Unix(),
		Seat:       0,
	}
	room.Players = players
	_, err := mgo.InsertOne(SchemaRoom, room)
	if err != nil {
		pkt.Fail(EcRoomNew_ExistName)
		return
	}
	res.Room = room
	pkt.Ok(res)

	common.ReqGateUpdate(pkt.Tid(), pkt.SenderId(), pkt.HeadId(), util.M{
		common.HdRoomId: roomId,
	}, nil)
}

func getNewOwnerId(room *pb.Room) string {
	ts := int64(math.MaxInt64)
	ownerId := ""
	for playerId, p := range room.Players {
		if p.EnterTs < ts {
			ts = p.EnterTs
			ownerId = playerId
		}
	}
	return ownerId
}

func ntcRoomToPlayers(tid int64, room *pb.Room, ntc util.IMsg) {
	idSlc := make([]string, 0, len(room.Players))
	for playerId := range room.Players {
		idSlc = append(idSlc, playerId)
	}
	common.ReqGateToMultiId(tid, idSlc, ntc)
}

func (s *Svc) OnRoomList(pkt kiwi.IRcvRequest, req *pb.RoomListReq, res *pb.RoomListRes) {
	filter := make([]bson.M, 0, 3)
	if req.NameFilter != "" {
		filter = append(filter, bson.M{
			"$text": bson.M{
				"$search": req.NameFilter,
			},
		})
	}
	if req.ModeFilter > 0 {
		filter = append(filter, bson.M{
			Mode: req.ModeFilter,
		})
	}
	if req.SceneTplFilter > 0 {
		filter = append(filter, bson.M{
			SceneTplId: req.SceneTplFilter,
		})
	}
	opt := options.Find()
	if len(req.Projection) > 0 {
		d := bson.D{}
		mgo.BuildProjectionD(&d, []string{Password}, req.Projection...)
		opt.SetProjection(d)
	}
	var mFilter bson.M
	switch len(filter) {
	case 0:
		mFilter = bson.M{}
	case 1:
		mFilter = filter[0]
	default:
		mFilter = bson.M{
			"$and": filter,
		}
	}
	mgo.FindWithTotal[pb.Room](SchemaRoom, mFilter,
		func(i int64, rooms []*pb.Room, e error) {
			if e != nil {
				pkt.Err(util.WrapErr(util.EcDbErr, e))
				return
			}
			res.List = rooms
			pkt.Ok(res)
		}, opt)
}

func (s *Svc) OnRoomEntry(pkt kiwi.IRcvRequest, req *pb.RoomEntryReq, res *pb.RoomEntryRes) {
	_, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if ok {
		pkt.Fail(EcRoomEntry_HadEntered)
		return
	}
	core.AsyncSubReq[*pb.PlayerIdRes](pkt, util.M{
		"id": pkt.HeadId(),
	}, &pb.PlayerIdReq{
		Id: pkt.HeadId(),
		Projection: []string{
			player.Nick,
			player.Hero,
		},
	}, func(tid int64, m util.M, u uint16) {
		pkt.Fail(util.EcServiceErr)
	}, func(tid int64, m util.M, playerRes *pb.PlayerIdRes) {
		var room pb.Room
		seat := int32(-1)
		e := mgo.FindOne(SchemaRoom, bson.D{
			{Id, req.RoomId},
		}, &room, options.FindOne().SetProjection(bson.D{
			{CreateTime, 0},
		}))
		if e != nil {
			if errors.Is(e, mongo.ErrNoDocuments) {
				pkt.Fail(EcRoomEntry_NotExistRoomId)
				return
			}
			pkt.Fail(util.EcDbErr)
			return
		}
		if room.Status != pb.RoomStatus_InLobby {
			pkt.Fail(EcRoomEntry_CurrStatusCanNotEntry)
			return
		}
		if room.CurrPlayers == room.MaxPlayers {
			pkt.Fail(EcRoomEntry_RoomFull)
			return
		}
		seatMap := make(map[int32]struct{})
		for _, p := range room.Players {
			seatMap[p.Seat] = struct{}{}
		}

		for i := int32(0); i < room.MaxPlayers; i++ {
			if _, ok := seatMap[i]; !ok {
				seat = i
				break
			}
		}
		_, e = mgo.UpdateOne(SchemaRoom, bson.D{
			{Id, req.RoomId},
		}, bson.M{
			"$inc": bson.M{CurrPlayers: 1},
			"$set": bson.M{
				fmt.Sprintf("%s.%s", Players, pkt.HeadId()): &pb.RoomPlayer{
					Ready:      false,
					Disconnect: false,
					EnterTs:    time.Now().Unix(),
					Seat:       seat,
				},
			},
		})
		if e != nil {
			if errors.Is(e, mongo.ErrNoDocuments) {
				pkt.Fail(EcRoomEntry_NotExistRoomId)
				return
			}
			pkt.Fail(util.EcDbErr)
			return
		}

		res.Room = &room
		pkt.Ok(res)

		common.ReqGateUpdate(pkt.Tid(), pkt.SenderId(), pkt.HeadId(), util.M{
			common.HdRoomId: req.RoomId,
		}, nil)

		idSlc := make([]string, 0, len(room.Players))
		for playerId := range room.Players {
			idSlc = append(idSlc, playerId)
		}

		p := playerRes.Player
		common.ReqGateToMultiId(pkt.Tid(), idSlc, &pb.RoomEntryPus{
			PlayerId: p.Id,
			Nick:     p.Nick,
			Hero:     p.Hero,
			Ready:    false,
			Seat:     seat,
		})
	})
}

func (s *Svc) OnRoomExit(pkt kiwi.IRcvRequest, req *pb.RoomExitReq, res *pb.RoomExitRes) {
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		pkt.Fail(EcRoomExit_NotEntryRoom)
		return
	}
	var (
		room  pb.Room
		empty bool
	)
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, &room, options.FindOne().SetProjection(
		bson.D{
			{Players, 1},
		}))
	if e != nil {
		if errors.Is(e, mongo.ErrNoDocuments) {
			pkt.Fail(EcRoomExit_NotExistRoomId)
			return
		}
		pkt.Fail(util.EcDbErr)
		return
	}
	//只有一个玩家，直接销毁房间
	playerCount := len(room.Players)
	if playerCount == 1 {
		empty = true
		_, e := mgo.DelOne(SchemaRoom, bson.D{
			{Id, roomId},
		})
		if e != nil {
			pkt.Fail(util.EcDbErr)
			return
		}
	}

	//移除玩家的房间数据
	update := bson.M{
		"$inc": bson.M{
			CurrPlayers: -1,
		},
		"$unset": bson.M{
			fmt.Sprintf("%s.%s", Players, pkt.HeadId()): "",
		},
	}
	if room.OwnerId == pkt.HeadId() {
		ownerId := getNewOwnerId(&room)
		update["$set"] = bson.M{
			OwnerId: ownerId,
		}
	}
	_, e = mgo.UpdateOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, update)
	if e != nil {
		pkt.Fail(util.EcDbErr)
		return
	}

	common.ReqGateRemove(pkt.Tid(), pkt.SenderId(), pkt.HeadId(),
		[]string{"rid"},
		[]string{strconv.Itoa(int(common.Room))},
	)

	if !empty {
		//广播给房间其他玩家
		ntcRoomToPlayers(pkt.Tid(), &room, &pb.RoomExitPus{
			PlayerId: pkt.HeadId(),
		})
	}
	pkt.Ok(res)
}

func (s *Svc) OnRoomReady(pkt kiwi.IRcvRequest, req *pb.RoomReadyReq, res *pb.RoomReadyRes) {
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		pkt.Fail(EcRoomReady_NotEntryRoom)
		return
	}
	var room pb.Room
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, &room, options.FindOne().SetProjection(bson.D{
		{Players, 1},
	}))
	if e != nil {
		if errors.Is(e, mongo.ErrNoDocuments) {
			pkt.Fail(EcRoomReady_NotExistRoomId)
			return
		}
		pkt.Fail(util.EcDbErr)
		return
	}
	pid := pkt.HeadId()
	ok = false
	for playerId, p := range room.Players {
		if playerId == pid {
			if p.Ready == req.IsReady {
				pkt.Ok(res)
				return
			}
			ok = true
			break
		}
	}
	if !ok {
		pkt.Fail(EcRoomReady_NotEntryRoom)
		return
	}
	_, e = mgo.UpdateOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, bson.M{
		"$set": bson.M{
			fmt.Sprintf("%s.%s.%s", Players, pid, RoomPlayerReady): req.IsReady,
		},
	})
	if e != nil {
		pkt.Fail(util.EcDbErr)
		return
	}
	pkt.Ok(res)
	ntcRoomToPlayers(pkt.Tid(), &room, &pb.RoomReadyPus{
		PlayerId: pkt.HeadId(),
		Ready:    req.IsReady,
	})
}

func (s *Svc) OnRoomStart(pkt kiwi.IRcvRequest, req *pb.RoomStartReq, res *pb.RoomStartRes) {
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		pkt.Fail(EcRoomStart_NotEntryRoom)
		return
	}
	var room pb.Room
	err := mgo.FindOneAndUpdate(SchemaRoom, bson.D{
		{Id, roomId},
		{OwnerId, pkt.HeadId()},
	}, bson.M{
		"$set": bson.D{
			{Status, pb.RoomStatus_InScene},
		},
	}, &room)
	if err != nil {
		pkt.Fail(EcRoomStart_NotExistRoomOrNotOwner)
		return
	}
	admitted := make([]*pb.SceneAdmittedPlayer, 0, len(room.Players))
	for playerId, p := range room.Players {
		if playerId != pkt.HeadId() && !p.Ready {
			pkt.Fail(EcRoomStart_SomebodyNotReady)
			return
		}
		admitted = append(admitted, &pb.SceneAdmittedPlayer{
			Id:      playerId,
			Entered: false,
		})
	}
	core.AsyncSubReq[*pb.NewSceneRes](pkt, nil, &pb.NewSceneReq{
		Scene: &pb.Scene{
			Id:       room.Id,
			TplId:    room.SceneTplId,
			Mode:     room.Mode,
			Admitted: nil,
		},
	}, func(tid int64, m util.M, code uint16) {
		pkt.Fail(EcRoomStart_CreateSceneFail)
	}, func(tid int64, m util.M, newSceneRes *pb.NewSceneRes) {
		pkt.Ok(res)
		ntcRoomToPlayers(tid, &room, &pb.RoomStartPus{})
	})
}

func (s *Svc) OnRoomGet(pkt kiwi.IRcvRequest, req *pb.RoomGetReq, res *pb.RoomGetRes) {
	var room pb.Room
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, req.RoomId},
	}, &room)
	if e != nil {
		if errors.Is(e, mongo.ErrNoDocuments) {
			pkt.Fail(EcRoomGet_NotExistRoomId)
			return
		}
		pkt.Err(util.WrapErr(util.EcDbErr, e))
		return
	}
	res.Room = &room
	pkt.Ok(res)
}

func (s *Svc) OnRoomModify(pkt kiwi.IRcvRequest, req *pb.RoomModifyReq, res *pb.RoomModifyRes) {
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		pkt.Fail(EcRoomModify_NotEntryRoom)
		return
	}
	var room pb.Room
	e := mgo.FindOneAndUpdate(SchemaRoom, bson.D{
		{Id, roomId},
		{OwnerId, pkt.HeadId()},
	}, bson.M{
		"$set": bson.D{
			{Name, req.Name},
			{SceneTplId, req.SceneTplId},
			{Password, req.Password},
			{MaxPlayers, req.MaxPlayers},
		},
	}, &room)
	if e != nil {
		pkt.Err(util.WrapErr(EcRoomModify_NotOwner, e))
		return
	}
	room.Name = req.Name
	room.SceneTplId = req.SceneTplId
	room.Password = req.Password
	room.MaxPlayers = req.MaxPlayers
	ntcRoomToPlayers(pkt.Tid(), &room, &pb.RoomModifyPus{
		Room: &room,
	})
	pkt.Ok(res)
}

func (s *Svc) OnRoomReconnect(pkt kiwi.IRcvRequest, req *pb.RoomReconnectReq, res *pb.RoomReconnectRes) {
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		pkt.Fail(util.EcServiceErr)
		return
	}

	var room pb.Room
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, &room, options.FindOne().SetProjection(bson.D{
		{CreateTime, 0},
	}))
	if e != nil {
		if errors.Is(e, mongo.ErrNoDocuments) {
			pkt.Fail(EcRoomReconnect_NotExistRoomId)
			return
		}
		pkt.Fail(util.EcDbErr)
		return
	}
	playerId := pkt.HeadId()
	ok = false
	for id := range room.Players {
		if id == playerId {
			ok = true
			break
		}
	}
	if !ok {
		pkt.Fail(EcRoomReconnect_KickOut)
		return
	}
	key := fmt.Sprintf("%s.%s.%s", Players, playerId, RoomPlayerDisconnect)
	_, e = mgo.UpdateOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, bson.M{
		"$set": bson.M{
			key: false,
		},
	})
	if e != nil && !errors.Is(e, mongo.ErrNoDocuments) {
		pkt.Fail(util.EcDbErr)
		return
	}
	res.Room = &room
	pkt.Ok(res)

	//广播给房间其他玩家
	if len(room.Players) > 0 {
		idSlc := make([]string, 0, len(room.Players))
		for playerId := range room.Players {
			if playerId == pkt.HeadId() {
				continue
			}
			idSlc = append(idSlc, playerId)
		}
		if len(idSlc) > 0 {
			common.ReqGateToMultiId(pkt.Tid(), idSlc, &pb.RoomReconnectPus{
				PlayerId: pkt.HeadId(),
			})
		}
	}
}
