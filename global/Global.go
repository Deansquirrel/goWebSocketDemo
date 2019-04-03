package global

import (
	"context"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goWebSocketDemo/config"
	"time"
)

const (
	//PreVersion = "0.0.0 Build20190328"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

const (
	HttpConnectTimeout = 30
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()

func init() {
	goToolMSSql.SetMaxIdleConn(15)
	goToolMSSql.SetMaxOpenConn(15)
	goToolMSSql.SetMaxLifetime(time.Second * 60)
}
