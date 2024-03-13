package scene

import (
	"game/proto/pb"
	"github.com/15mga/kiwi/ds"
	"github.com/15mga/kiwi/ecs"
	"github.com/15mga/kiwi/sid"
	"github.com/15mga/kiwi/util"
	"github.com/15mga/kiwi/worker"
	"math/rand"
	"time"
)

var (
	MonsterTpl = []string{
		"Catcher",
		"Fishguard",
		"Imp",
		"Knight",
		"Monkeydong",
		"Nosedman",
		"Pitboy",
		"Spike",
		"Treestor",
		"Wedger",
	}
	MonsterLevel = []string{
		"Big",
		"Medium",
		"Small",
	}
)

func randMonsterTplId() string {
	tpl := MonsterTpl[rand.Intn(len(MonsterTpl))]
	lvl := MonsterLevel[rand.Intn(len(MonsterLevel))]
	return tpl + "_" + lvl
}

func NewSRobot(maxRobot int32) *SRobot {
	return &SRobot{
		System:   ecs.NewSystem(S_Robot),
		maxRobot: maxRobot,
	}
}

type SRobot struct {
	ecs.System
	sysTnf    *STransform
	maxRobot  int32
	currRobot int32
}

func (s *SRobot) OnBeforeStart() {
	s.System.OnBeforeStart()
	s.BindJob(JobRobotAdd, s.onRobotAdd)
	s.BindJob(JobRobotClear, s.onRobotClear)
}

func (s *SRobot) OnAfterStart() {
	s.System.OnAfterStart()
	sysTnf, _ := s.Frame().GetSystem(S_Transform)
	s.sysTnf = sysTnf.(*STransform)
}

func (s *SRobot) OnUpdate() {
	s.DoJob(JobRobotAdd)
	s.DoJob(JobRobotClear)
	s.updateRobot()
}

func (s *SRobot) onRobotAdd(data []any) {
	_, count, handler := util.SplitSlc3[int64, int32, func(int32)](data)
	m := s.maxRobot - s.currRobot
	if m == 0 {
		handler(s.currRobot)
		return
	}
	if count > m {
		count = m
	}
	s.currRobot += count
	for i := int32(0); i < count; i++ {
		pawnId := sid.GetStrId()
		e := ecs.NewEntity(pawnId)
		var pos pb.Vector2
		s.sysTnf.GenRandPos(&pos)
		//pos.X = rand.Float32() * 200
		//pos.Y = rand.Float32() * 200
		cpawn := NewCMonster(&pb.SceneMonster{
			TplId: randMonsterTplId(),
		})
		tnf := NewCTransform(&pos)
		e.AddComponents(
			NewCBehaviour(&pb.SceneBehaviour{
				Timestamp: time.Now().UnixMilli(),
				BehaviourType: &pb.SceneBehaviour_Idle{
					Idle: &pb.BehaviourIdle{},
				},
			}),
			tnf,
			NewCRobot(),
			cpawn,
			NewCEvent(cpawn),
		)
		_ = s.Scene().AddEntity(e)
		s.Scene().TagComponent(tnf, TagCompSceneEntry)
	}
	handler(s.currRobot)
}

func (s *SRobot) onRobotClear(data []any) {
	tid := data[0].(int64)
	components, ok := s.Scene().GetTagComponents(string(C_Robot))
	if !ok {
		return
	}
	ids := make([]string, len(components))
	for i, component := range components {
		ids[i] = component.Entity().Id()
	}
	s.currRobot = 0
	s.Frame().PutJob(JobEntityDel, tid, ids)
}

func (s *SRobot) updateRobot() {
	components, ok := s.Scene().GetTagComponents(string(C_Robot))
	if !ok {
		return
	}
	worker.PToLink[ecs.IComponent, ecs.IComponent](components, func(component ecs.IComponent, d *ds.Link[ecs.IComponent]) {
		if !component.(*CRobot).Update(s.Frame().DeltaMillSec()) {
			return
		}
		d.Push(component)
	}, func(d *ds.Link[ecs.IComponent]) {
		d.Iter(s.putMovementJob)
	})
}

func (s *SRobot) putMovementJob(component ecs.IComponent) {
	var dir pb.Vector2
	Vec2ToPbVec(util.RandDir(), &dir)
	s.Frame().PutJob(JobMovement, int64(0), component.Entity().Id(), &pb.SceneMovement{
		Direction: &dir,
		MoveSpeed: 4,
	})
}