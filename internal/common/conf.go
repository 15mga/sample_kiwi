package common

import (
	"github.com/15mga/kiwi"
	"github.com/15mga/kiwi/loader"
	etcd "go.etcd.io/etcd/client/v3"
)

var (
	_Conf Config
)

func Conf() *Config {
	return &_Conf
}

type Mgo struct {
	Uri string
	Db  string
}

type Register struct {
	Expire int64
}

type NodeConf struct {
	Ip   string
	Port int
}

type GateConf struct {
	Ip           string
	Tcp          int
	Web          int
	Udp          int
	Http         int
	ConnCap      int32
	DeadlineSecs int
	PacketLimit  uint16
	ErrLimit     uint16
	TickSeconds  uint16
	Auth         map[string][]string
	Msg          map[string]map[string]uint16
	AdminMsg     map[string]map[string]uint16
}

type LogStd struct {
	Enable bool
	Color  bool
	Log    []string
	Trace  []string
}

type LogMgo struct {
	Enable bool
	Log    []string
	Trace  []string
	Uri    string
	Db     string
	Ttl    int32
}

type LogConfig struct {
	Std     LogStd
	Mgo     LogMgo
	Exclude map[string][]kiwi.TCode
}

type RedisConf struct {
	Addr     string
	User     string
	Password string
	Db       int
}

type TestConf struct {
	MaxRobot int32
}

type Config struct {
	Id       int64
	Mode     string
	Log      LogConfig
	Etcd     *etcd.Config
	Redis    RedisConf
	Mongo    Mgo
	Node     NodeConf
	Gate     GateConf
	Local    bool
	SvcToVer map[string]string
	Test     TestConf
}

func LoadConf(confFolder string) {
	loader.SetConfRoot(confFolder)

	// todo 添加其他加载方式,远程,etcd
	convPath := loader.ConvertConfLocalPath
	//if strings.HasPrefix(confFolder, "etcd") {
	//
	//} else if strings.HasPrefix(confFolder, "http") {
	//
	//} else if strings.HasPrefix(confFolder, "https") {
	//
	//}
	slc := []string{
		"game.yml",       //配置模本
		"game_local.yml", //本地配置,用于覆盖模板配置
	}
	loader.LoadConf(Conf(), convPath(slc...)...)
}
