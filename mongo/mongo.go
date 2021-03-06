package mongo

import (
	"fmt"
	"github.com/dishine/libary/log"
	"github.com/dishine/libary/node"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"time"
)

var session *mgo.Session

var (
	Nil = mgo.ErrNotFound
)

type Config struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	Open     int64  `json:"open" yaml:"open"`
	Idle     int64  `json:"idle" yaml:"idle"`
}

func NewLoad(config *Config) *Load {
	return &Load{
		config: config,
	}
}

type Load struct {
	config *Config
}

func (m *Load) GetOrder() node.Order {
	return node.Before
}
func (m *Load) GetOptionFunc() node.Func {
	return m.Connect
}

func (m *Load) Connect() error {
	config := m.config
	dialInfo := mgo.DialInfo{
		Addrs:    []string{config.Host + ":" + config.Port},
		Timeout:  time.Second * 3,
		Database: config.Database,
		Username: config.User,
		Password: config.Password,
	}
	ses, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Info("mongo 连接错误")
		panic(err.Error())
	}
	ses.SetSocketTimeout(time.Hour)
	//if c.IsDevelop() {
	//mgo.SetDebug(true)
	//mgo.SetLogger(new(MongoLog))
	//}
	session = ses
	log.Info("load mongo success", zap.String("conn", config.Host+":"+config.Port))
	return nil
}

func Conn() *mgo.Session {
	if session == nil {
		panic("获取mongo连接为空")
	}
	if err := session.Clone().Ping(); err != nil {
		panic(fmt.Sprintf("获取mongo连接 Ping出错 - [%s]", err.Error()))
	}
	return session
}

func InitConnect(config *Config) {
	dialInfo := mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  time.Second * 3,
		Database: config.Database,
		Username: config.User,
		Password: config.Password,
	}
	ses, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Info("mongo 连接错误")
		panic(err.Error())
	}
	ses.SetSocketTimeout(time.Hour)

	//if c.IsDevelop() {
	//mgo.SetDebug(true)
	//mgo.SetLogger(new(MongoLog))
	//}
	session = ses
	log.Info("load mongo success", zap.String("conn", config.Host))
}

type MongoLog struct {
}

func (MongoLog) Output(calldepth int, s string) error {
	log.Info(fmt.Sprintf(" %v , %v", calldepth, s))
	return nil
}
