package setup

import (
	"fmt"
	kratos "github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/config"
	"github.com/liangguifeng/kratos-app/config/setting"
	"github.com/liangguifeng/kratos-app/internal/module/helper"
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
	} else if projectLoggerPath := config.GetProjectLoggerPath(); projectLoggerPath != "" {
		rootPath = projectLoggerPath
	}
	//err := yklog.InitGlobalConfig(rootPath, "debug", application.Name)
	//if err != nil {
	//	return fmt.Errorf("Yklog.InitGlobalConfig: %v", err)
	//}

	fmt.Println(rootPath)
	return nil
}

// 初始化全局包变量
func NewGlobalVars() error {
	var err error
	if setting.MysqlSetting != nil && setting.MysqlSetting.Host != "" {
		kratos.MySQLConn, err = helper.NewMySQLConn(setting.MysqlSetting)
		if err != nil {
			return err
		}
	}

	if setting.RedisSetting != nil && setting.RedisSetting.Host != "" {
		kratos.RedisConn, err = helper.NewRedisConn(setting.RedisSetting)
		if err != nil {
			return err
		}
	}

	kratos.HTTPClient = helper.NewHTTPClient()
	return nil
}
