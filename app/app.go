package app

import (
	"flag"
	"fmt"
	kratos_app "github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/config"
	"github.com/liangguifeng/kratos-app/internal/setup"
	"github.com/pkg/errors"
)

var (
	port       = flag.Int64("p", 0, "Set server port.")
	loggerPath = flag.String("logger_path", "", "Set Logger Root Path.")
)

func NewRunner(app *kratos_app.Application) (*kratos_app.Runner, error) {
	if app.Name == "" {
		return nil, errors.New("Application name can't not be empty")
	}
	if app.Type <= 0 {
		return nil, errors.New("Application type can't not be empty")
	}

	// 加载指定运行端口和日志路径
	flag.Parse()
	app.Port = *port
	app.LoggerRootPath = *loggerPath

	// 获取当前环境变量
	goEnv := config.GetBuildEnv()
	if goEnv == "" {
		return nil, fmt.Errorf("Can't not found env '%s' or '%s'", config.PROJECT_ENV, config.GO_ENV)
	}

	// 日志组件初始化
	var err error
	err = setup.NewLogGlobalConfig(app)
	if err != nil {
		return nil, err
	}

	return &kratos_app.Runner{App: app}, nil
}
