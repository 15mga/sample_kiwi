package scene

import (
	"fmt"
	"game/proto/pb"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/log"
	"github.com/15mga/kiwi/sid"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/worker"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func initEgFrame() {
	kiwi.AddLogger(log.NewStd())
	core.InitNodeTest()
	core.InitCodec()
	core.InitRouter()
	worker.InitParallel()
	sid.SetNodeId(0)
}

func TestNewScene(t *testing.T) {
	initEgFrame()
	sceneId := sid.GetStrId()
	pbScene := &pb.Scene{
		Id:       sceneId,
		TplId:    0,
		Mode:     pb.ESceneMode_PVE,
		Admitted: nil,
	}
	NewScene(0, pbScene)
	scene, ok := GetSceneDataById(sceneId)
	assert.True(t, ok)
	assert.Equal(t, pbScene, scene)
	var playerIds []string
	tc := int32(8)
	pc := int32(1)
	kiwi.Debug("test", util.M{
		"count": tc * tc * pc,
	})
	for y := int32(0); y < tc; y++ {
		for x := int32(0); x < tc; x++ {
			for i := int32(0); i < pc; i++ {
				var player pb.ScenePawn
				var pos pb.Vector2
				getTilePos(x, y, &pos)
				playerId := fmt.Sprintf("player_%d_%d", x, y)
				newPlayer(&player, playerId, &pos)
				PushJob(sceneId, JobEntityAdd, int64(0), &player, func(err *util.Err) {
					//t.Log(code)
				})
				playerIds = append(playerIds, playerId)
			}
		}
	}
	time.Sleep(time.Second)

	var player pb.ScenePawn
	var pos pb.Vector2
	getTilePos(0, 0, &pos)
	playerId := "move_player"
	newPlayer(&player, playerId, &pos)
	kiwi.Debug("add move player", nil)
	PushJob(sceneId, JobEntityAdd, int64(0), &player, func(err *util.Err) {
		//t.Log(code)
	})

	kiwi.Debug("start move player", nil)
	time.Sleep(time.Second)
	PushJob(sceneId, JobMovement, int64(0), playerId, &pb.SceneTransform{
		Position: &pos,
		Direction: &pb.Vector2{
			X: 1,
			Y: 1,
		},
		MoveSpeed: 10,
		Timestamp: time.Now().UnixMilli(),
	})
	time.Sleep(time.Second * 500)
}

func getTilePos(tx, ty int32, pos *pb.Vector2) {
	pos.X = (float32(tx) + 0.5) * _TileSize
	pos.Y = (float32(ty) + 0.5) * _TileSize
}

func newPlayer(player *pb.ScenePawn, id string, pos *pb.Vector2) {
	player.Transform = &pb.SceneTransform{
		Position: pos,
		Direction: &pb.Vector2{
			X: 1,
			Y: 0,
		},
		MoveSpeed: 0,
		Timestamp: time.Now().UnixMilli(),
	}
	player.PawnType = &pb.ScenePawn_Player{
		Player: &pb.ScenePlayer{
			Id:     id,
			Nick:   "player_" + id,
			Avatar: "avatar_" + id,
			TeamId: "team_" + id,
			Gender: rand.Int31n(3),
		},
	}
	player.PlayerGateAddr = "127.0.0.1:7737"
	player.PlayerGateNodeId = 1
}

func TestTileChange(t *testing.T) {
	kiwi.AddLogger(log.NewStd())
	width := float32(100)
	height := float32(100)
	sys := NewSTransform(10, 100, 100, 1)
	pos := &pb.Vector2{
		X: 35,
		Y: 35,
	}
	tnf := NewCTransform(&pb.SceneTransform{
		Position: pos,
		Direction: &pb.Vector2{
			X: 1,
			Y: 1,
		},
		MoveSpeed: 30,
		Timestamp: 0,
	})
	tnf.InitTile(sys.posToTile(tnf.TransformData.Position))
	tnf.ProcessMovement(1000, width, height)
	tnf.UpdateTile(sys.posToTile(tnf.TransformData.Position))
	var tcd TileChangeData
	sys.getInterestTileChanged(tnf, &tcd)
	kiwi.Debug("tile changed", util.M{
		"interest": tcd,
	})
}
