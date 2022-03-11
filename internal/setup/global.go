package setup

import (
	"fmt"
	kratos_app "github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/config"
	"os"
	"runtime"
)

// 初始化全局日志配置（指定公司内部约定路径和服务名）
func NewLogGlobalConfig(application *kratos_app.Application) error {
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

//// 初始化全局包变量
//func NewGlobalVars() error {
//	var err error
//	if setting.MysqlSetting != nil && setting.MysqlSetting.Host != "" {
//		stark.MySQLConn, err = helper.NewMySQLConn(setting.MysqlSetting)
//		if err != nil {
//			return err
//		}
//	}
//
//	if setting.RedisSetting != nil && setting.RedisSetting.Host != "" {
//		stark.RedisConn, err = helper.NewRedisConn(setting.RedisSetting)
//		if err != nil {
//			return err
//		}
//	}
//
//	stark.HTTPClient = helper.NewHTTPClient()
//	return nil
//}
