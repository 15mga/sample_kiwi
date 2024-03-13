// Code generated by protoc-gen-go-kiwi. DO NOT EDIT.

package common

import (
	"github.com/15mga/kiwi"
)

const (
	Gate   kiwi.TSvc = 2
	User   kiwi.TSvc = 3
	Player kiwi.TSvc = 4
	Room   kiwi.TSvc = 5
	Scene  kiwi.TSvc = 7
	Team   kiwi.TSvc = 8
	Chat   kiwi.TSvc = 9
)

const (
	SGate   = "gate"
	SUser   = "user"
	SPlayer = "player"
	SRoom   = "room"
	SScene  = "scene"
	STeam   = "team"
	SChat   = "chat"
)

var SvcToName = map[kiwi.TSvc]string{
	Gate:   SGate,
	User:   SUser,
	Player: SPlayer,
	Room:   SRoom,
	Scene:  SScene,
	Team:   STeam,
	Chat:   SChat,
}

var NameToSvc = map[string]kiwi.TSvc{
	SGate:   Gate,
	SUser:   User,
	SPlayer: Player,
	SRoom:   Room,
	SScene:  Scene,
	STeam:   Team,
	SChat:   Chat,
}
