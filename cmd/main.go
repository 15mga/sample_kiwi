package main

import (
	"fmt"
	"game/internal/chat"
	"game/internal/common"
	"game/internal/gate"
	"game/internal/player"
	"game/internal/room"
	"game/internal/scene"
	"game/internal/team"
	"game/internal/user"
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/core"
	"github.com/15mga/kiwi/log"
	"github.com/15mga/kiwi/util/etd"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go http.ListenAndServe("0.0.0.0:6060", nil)
	wd, _ := os.Getwd()
	//加载配置文件
	common.LoadConf(fmt.Sprintf("%s/config", wd))
	conf := common.Conf()
	nodeInfo := kiwi.GetNodeMeta()
	nodeInfo.Mode = conf.Mode

	// 分布式服务使用，暂时没用
	if conf.Etcd != nil {
		err := etd.Conn(*conf.Etcd)
		if err != nil {
			panic(err)
		}
	}

	//日志过滤
	for svc, codes := range conf.Log.Exclude {
		core.ExcludeLog(common.NameToSvc[svc], codes...)
	}

	//文件日志
	var loggers []kiwi.ILogger
	logStd := conf.Log.Std
	if logStd.Enable {
		loggers = append(loggers, log.NewStd(
			log.StdColor(logStd.Color),
			log.StdLogStrLvl(logStd.Log...),
			log.StdTraceStrLvl(logStd.Trace...),
			log.StdFile(fmt.Sprintf("%s/log/game.log", wd)),
		))
	}

	//mongo日志
	logMgo := conf.Log.Mgo
	if logMgo.Enable {
		loggers = append(loggers, log.NewMgo(
			log.MgoLogLvl(logMgo.Log...),
			log.MgoTraceLvl(logMgo.Trace...),
			log.MgoClientOptions(options.Client().ApplyURI(logMgo.Uri)),
			log.MgoDb(logMgo.Db),
			log.MgoTtl(logMgo.Ttl),
		))
	}

	// 设置服务相关功能
	core.StartDefault(
		core.SetMeta(&core.Meta{
			Id:       conf.Id, //服务节点唯一 id，不能重复
			SvcToVer: conf.SvcToVer,
			SvcNameConv: func(s string) kiwi.TSvc {
				return common.NameToSvc[s]
			},
		}),
		core.SetLoggers(loggers...),
		core.SetMongoDB(conf.Mongo.Uri, conf.Mongo.Db, nil),
		core.SetRedis(conf.Redis.Addr, conf.Redis.User, conf.Redis.Password, conf.Redis.Db),
		//设置服务
		core.SetServices([]kiwi.IService{
			chat.Service(),
			gate.Service(),
			player.Service(),
			room.Service(),
			scene.Service(),
			team.Service(),
			user.Service(),
		}),
		//网关设置
		core.SetGate(gate.SocketReceiver,
			core.GateRoles(common.MsgRole),
			core.GateIp(conf.Gate.Ip),
			core.GateWebsocketPort(conf.Gate.Web),
			core.GateConnCap(conf.Gate.ConnCap),
			core.GateDeadlineSecs(conf.Gate.DeadlineSecs),
			core.GateCheckIp(gate.CheckIp),
			core.GateDisconnected(gate.Disconnected),
			core.GateHeadLen(4),
			core.GateConnected(gate.InitAgent),
		))
	core.StartAllService()
	kiwi.WaitExit()
}
