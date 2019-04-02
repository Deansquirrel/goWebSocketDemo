package common

import (
	"github.com/BurntSushi/toml"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goWebSocketDemo/config"
)

//获取配置
func GetSysConfig(fileName string) (*config.SysConfig, error) {
	path, err := goToolCommon.GetCurrPath()
	if err != nil {
		return nil, err
	}
	var c config.SysConfig
	fileFullPath := path + "\\" + fileName
	b, err := goToolCommon.PathExists(fileFullPath)
	if err != nil {
		log.Warn("检查路径是否存在时遇到错误:" + err.Error() + ",使用默认配置;filePath:" + fileFullPath)
		c = config.SysConfig{}
	} else if !b {
		log.Info("未发现配置文件,使用默认配置" + ";filePath:" + fileFullPath)
		c = config.SysConfig{}
	} else {
		_, err = toml.DecodeFile(fileFullPath, &c)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

//刷新服务端配置
func RefreshSysConfig(c config.SysConfig) error {
	refreshTotalConfig(c.Total)
	return nil
}

//刷新Total配置
func refreshTotalConfig(t config.Total) {
	setLogLevel(t.LogLevel)
	setStdOut(t.StdOut)
}

//设置标准输出
func setStdOut(isStdOut bool) {
	log.StdOut = isStdOut
}

//设置日志级别
func setLogLevel(logLevel string) {
	switch logLevel {
	case "debug":
		log.Level = log.LevelDebug
		return
	case "info":
		log.Level = log.LevelInfo
		return
	case "warn":
		log.Level = log.LevelWarn
		return
	case "error":
		log.Level = log.LevelError
		return
	default:
		log.Level = log.LevelWarn
	}
}
