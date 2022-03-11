package setup

import (
	"fmt"
	kratos "github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/config/setting"
	config2 "github.com/liangguifeng/kratos-app/internal/config"
	"github.com/liangguifeng/kratos-app/logging/applog"
	helper2 "github.com/liangguifeng/kratos-app/module/helper"
	"os"
	"runtime"
)

// 初始化全局日志配置（指定公司内部约定路径和服务名）
func NewLogGlobalConfig(application *kratos.Application) error {
	rootPath := "/var/log/service"
	if runtime.GOOS == "windows" {
		rootPath = os.TempDir()
	}
	if application.LoggerRootPath != "" {
		rootPath = application.LoggerRootPath
	} else if projectLoggerPath := config2.GetProjectLoggerPath(); projectLoggerPath != "" {
		rootPath = projectLoggerPath
	}

	err := applog.InitGlobalConfig(rootPath, "debug", application.Name)
	if err != nil {
		return fmt.Errorf("applog.InitGlobalConfig: %v", err)
	}

	kratos.Logger, err = helper2.NewLogger()
	if err != nil {
		return err
	}
	return nil
}

// 初始化全局包变量
func NewGlobalVars() error {
	var err error
	if setting.MysqlSetting != nil && setting.MysqlSetting.Host != "" {
		kratos.MySQLConn, err = helper2.NewMySQLConn(setting.MysqlSetting)
		if err != nil {
			return err
		}
	}

	if setting.RedisSetting != nil && setting.RedisSetting.Host != "" {
		kratos.RedisConn, err = helper2.NewRedisConn(setting.RedisSetting)
		if err != nil {
			return err
		}
	}

	kratos.HTTPClient = helper2.NewHTTPClient()
	return nil
}
