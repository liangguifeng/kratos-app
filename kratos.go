package kratos

import (
	"context"
	grpc_server "github.com/go-kratos/kratos/v2/transport/grpc"
	http_server "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
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
)

type CallbackPos int

// Application ...
type Application struct {
	Name             string
	HTTPPort         string
	GRPCPort         string
	Version          string
	Type             int32
	LoggerRootPath   string
	LoadConfig       func() error
	SetupVars        func() error
	RegisterCallback map[CallbackPos]func() error
}

// GRPCApplication ...
type GRPCApplication struct {
	App                         *Application
	GRPCServer                  *grpc_server.Server
	HttpServer                  *http_server.Server
	UnaryServerInterceptors     []grpc.UnaryServerInterceptor
	ServerOptions               []grpc.ServerOption
	RegisterGRPCServer          func(*grpc_server.Server) error
	RegisterHttpRoute           func(*http_server.Server) error
	RegisterGracefulStopHandler func()
}

// Logger is a global vars for writing to the log.
var Logger LoggerIface

type LoggerIface interface {
	// 业务日志
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugw(ctx context.Context, err error)
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infow(ctx context.Context, err error)

	// 警告日志
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnw(ctx context.Context, err error)

	// 错误日志
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Errorw(ctx context.Context, err error)
}

// Configer is a global vars for read app config.
var Configer ConfigerIface

type ConfigerIface interface {
	LoadAppConfig(app *Application) error
	GetStringValue(key, defaultValue string) string
	GetAllKeys() map[string]interface{}
	GetIntValue(key string, defaultValue int) int
	GetBoolValue(key string, defaultValue bool) bool
	WatchUpdateConfig()
}

// DBConn is a global vars for mysql tracing connect.
var MySQLConn MySQLConnIface

type MySQLConnIface interface {
	GetClient() *gorm.DB
}

// RedisConn is a global vars for redis tracing connect.
var RedisConn RedisConnIface

type RedisConnIface interface {
	GetClient() *redis.Pool
	ActiveCount() int
	Close() error
	IdleCount() int
	Stats() redis.PoolStats
}

// HTTPClient is a global vars for http tracing client.
var HTTPClient HTTPClientIface

type HTTPClientIface interface {
	HttpGet(ctx context.Context, url string) (*http.Response, error)
	HttpPost(ctx context.Context, url string, bodyType string, body io.Reader) (*http.Response, error)
	PostForm(ctx context.Context, url string, data url.Values) (*http.Response, error)
	HttpHead(ctx context.Context, url string) (*http.Response, error)
	HttpDo(ctx context.Context, r *http.Request) (*http.Response, error)
}

// RunNewRunnerCallback
func (app *Application) RunNewRunnerCallback() error {
	if f, ok := app.RegisterCallback[POS_NEW_RUNNER]; ok {
		return f()
	}

	return nil
}

// RunLoadConfigCallback
func (app *Application) RunLoadConfigCallback() error {
	if f, ok := app.RegisterCallback[POS_LOAD_CONFIG]; ok {
		return f()
	}

	return nil
}

// RunSetupVarsCallback
func (app *Application) RunSetupVarsCallback() error {
	if f, ok := app.RegisterCallback[POS_SETUP_VARS]; ok {
		return f()
	}

	return nil
}
