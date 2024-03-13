package room

import (
	"fmt"
	"game/internal/common"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/util/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

func (s *Svc) OnPlayerDisconnectNtc(pkt kiwi.IRcvNotice, ntc *pb.PlayerDisconnectNtc) {
	defer pkt.Complete()
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		return
	}

	kiwi.TD(pkt.Tid(), "room receive player disconnect", util.M{
		"ntc": ntc,
	})
	var room pb.Room
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, &room, options.FindOne().SetProjection(bson.D{
		{Players, 1},
	}))
	if e != nil {
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
		return
	}
	key := fmt.Sprintf("%s.%s.%s", Players, playerId, RoomPlayerDisconnect)
	_, e = mgo.UpdateOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, bson.M{
		"$set": bson.D{
			{key, true},
		},
	})
	if e != nil {
		return
	}

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
			common.ReqGateToMultiId(pkt.Tid(), idSlc, &pb.RoomDisconnectPus{
				PlayerId: pkt.HeadId(),
			})
		}
	}
}

func (s *Svc) OnPlayerOfflineNtc(pkt kiwi.IRcvNotice, ntc *pb.PlayerOfflineNtc) {
	defer pkt.Complete()
	roomId, ok := util.MGet[string](pkt.Head(), common.HdRoomId)
	if !ok {
		return
	}

	kiwi.TD(pkt.Tid(), "room receive player offline", util.M{
		"ntc": ntc,
	})
	var (
		room        pb.Room
		playerCount int
	)
	e := mgo.FindOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, &room, options.FindOne().SetProjection(
		bson.D{
			{Players, 1},
			{OwnerId, 1},
		}))
	if e != nil {
		return
	}
	isMember := false
	for playerId := range room.Players {
		if playerId == pkt.HeadId() {
			isMember = true
			break
		}
	}
	if !isMember {
		return
	}
	//只有一个玩家，直接销毁房间
	playerCount = len(room.Players)
	if playerCount == 1 {
		_, _ = mgo.DelOne(SchemaRoom, bson.D{
			{Id, roomId},
		})
		return
	}

	key := util.StringsJoin(".", Players, pkt.HeadId())
	update := bson.M{
		"$inc": bson.M{
			CurrPlayers: -1,
		},
		"$unset": bson.M{
			key: "",
		},
	}
	if room.OwnerId == pkt.HeadId() {
		ownerId := getNewOwnerId(&room)
		update["$set"] = bson.M{
			OwnerId: ownerId,
		}
	}
	//移除玩家的房间数据
	_, e = mgo.UpdateOne(SchemaRoom, bson.D{
		{Id, roomId},
	}, update)
	if e != nil {
		return
	}

	common.ReqGateRemove(pkt.Tid(), pkt.SenderId(), pkt.HeadId(),
		[]string{"rid"},
		[]string{strconv.Itoa(int(common.Room))},
	)

	if playerCount > 1 {
		//广播给房间其他玩家
		idSlc := make([]string, 0, playerCount)
		for id := range room.Players {
			idSlc = append(idSlc, id)
		}

		//广播给房间其他玩家
		ntcRoomToPlayers(pkt.Tid(), &room, &pb.RoomExitPus{
			PlayerId: pkt.HeadId(),
		})
	}
}
