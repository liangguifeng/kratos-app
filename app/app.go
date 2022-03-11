package app

import (
	"flag"
	"fmt"
	"github.com/liangguifeng/kratos-app"
	config2 "github.com/liangguifeng/kratos-app/internal/config"
	"github.com/liangguifeng/kratos-app/internal/setup"
	"github.com/pkg/errors"
)

var (
	port       = flag.Int64("p", 0, "Set server port.")
	loggerPath = flag.String("logger_path", "", "Set Logger Root Path.")
)

type Runner struct {
	App *kratos.Application
}

func NewRunner(app *kratos.Application) (*Runner, error) {
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
	goEnv := config2.GetBuildEnv()
	if goEnv == "" {
		return nil, fmt.Errorf("Can't not found env '%s' or '%s'", config2.PROJECT_ENV, config2.GO_ENV)
	}

	// 日志组件初始化
	var err error
	// todo: liangguifeng 初始化日志
	err = setup.NewLogGlobalConfig(app)
	if err != nil {
		return nil, err
	}

	// 配置中心
	kratos.Configer, err = setup.NewConfiger(app)
	if err != nil {
		return nil, err
	}

	// 加载初始化配置到全局变量中
	err = kratos.Configer.LoadAppConfig(app)
	if err != nil {
		return nil, err
	}

	// 如果有需要手动设定全局变量的配置在此添加
	if app.LoadConfig != nil {
		// 加载手动声明的配置
		err = app.LoadConfig()
		if err != nil {
			return nil, err
		}
		// 执行回调
		err = app.RunLoadConfigCallback()
		if err != nil {
			return nil, err
		}
	}

	if app.SetupVars != nil {
		err = app.SetupVars()
		if err != nil {
			return nil, err
		}
		err = app.RunSetupVarsCallback()
		if err != nil {
			return nil, err
		}
	}

	// todo: liangguifeng 监听更新
	//kratos.Configer.WatchUpdateConfig()

	// 设置全局mysql连接池、redis连接池、http客户端
	err = setup.NewGlobalVars()
	if err != nil {
		return nil, err
	}

	// 运行回调
	err = app.RunNewRunnerCallback()
	if err != nil {
		return nil, err
	}

	return &Runner{App: app}, nil
}
