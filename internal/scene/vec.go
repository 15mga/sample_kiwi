package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/util"
)

func PbVecToVec2(v1 *pb.Vector2) util.Vec2 {
	return util.Vec2{
		X: v1.X,
		Y: v1.Y,
	}
}

func Vec2ToPbVec(v1 util.Vec2, v2 *pb.Vector2) {
	v2.X, v2.Y = v1.X, v1.Y
}
