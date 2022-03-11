package kratos_app

import "flag"

var (
	port       = flag.Int64("p", 0, "Set server port.")
	loggerPath = flag.String("logger_path", "", "Set Logger Root Path.")
)

const (
	APP_TYPE_GRPC     = 1
	APP_TYPE_CRON     = 2
	APP_TYPE_QUEUE    = 3
	APP_TYPE_HTTP     = 4
	APP_TYPE_LISTENER = 5
	APP_TYPE_GIN      = 6
)

const (
	// 调用 LoadConfig 方法后
	POS_LOAD_CONFIG CallbackPos = iota + 1
	// 调用 SetupVars 方法后
	POS_SETUP_VARS
	// 调用 NewRunner 方法后
	POS_NEW_RUNNER
	// 调用 RegisterEventHandler 方法后
	POS_REGISTER_EVENT_HANDLER
)

type CallbackPos int

type Runner struct {
	App *Application
}

// Application ...
type Application struct {
	Name             string
	Port             int64
	Type             int32
	LoggerRootPath   string
	LoadConfig       func() error
	SetupVars        func() error
	RegisterCallback map[CallbackPos]func() error
}
