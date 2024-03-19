package scene

import "github.com/15mga/kiwi/ecs"

const (
	TagCompSceneEntry = "scene_entry"
	TagCompSceneExit  = "scene_exit"
	TagCompMove       = "move"
	TagCompBehaviour  = "behaviour"
)

const (
	S_Tile       = "tile"
	S_Entity     = "entity"
	S_Event      = "event"
	S_Monster    = "monster"
	S_Ntc_Sender = "notice_sender"
	S_Transform  = "transform"
	S_Behaviour  = "behaviour"
	S_Robot      = "robot"
)

const (
	C_Behaviour ecs.TComponent = "behaviour"
	C_Event     ecs.TComponent = "event"
	C_Monster   ecs.TComponent = "monster"
	C_Player    ecs.TComponent = "player"
	C_Tile      ecs.TComponent = "tile"
	C_Transform ecs.TComponent = "transform"
	C_Robot     ecs.TComponent = "robot"
)

const (
	JobEntityAdd  = "entity_add"
	JobEntityDel  = "entity_del"
	JobMovement   = "movement"
	JobBehaviour  = "behaviour"
	JobSendNtc    = "send_ntc"
	JobRobotAdd   = "robot_add"
	JobRobotClear = "robot_clear"
)
