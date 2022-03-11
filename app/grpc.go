package app

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	server "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liangguifeng/kratos-app"
	"github.com/liangguifeng/kratos-app/internal/config"
	"github.com/liangguifeng/kratos-app/internal/setup"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
)

// ListenGRPCServer run grpc application.
func (r *Runner) ListenGRPCServer(grpcApp *kratos.GRPCApplication) error {
	grpcApp.App = r.App
	err := r.handleGRPC(grpcApp)
	if err != nil {
		return err
	}

	return nil
}

// Run gRPC application handle.
func (r *Runner) handleGRPC(grpcApp *kratos.GRPCApplication) error {
	httpSrv := http.NewServer(
		http.Address(":"+grpcApp.App.HTTPPort),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":"+grpcApp.App.GRPCPort),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	// 注册HTTP路由
	err := grpcApp.RegisterHttpRoute(httpSrv)
	if err != nil {
		return err
	}
	// 注册GRPC路由
	err = grpcApp.RegisterGRPCServer(grpcSrv)
	if err != nil {
		return err
	}

	id, _ := os.Hostname()
	serverId := id + grpcApp.App.Name + " _service"
	logger := log.With(
		log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", serverId,
		"service.name", grpcApp.App.Name,
		"service.version", grpcApp.App.Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.GetNacosNamespaceId(),
			TimeoutMs:           setup.NACOS_TIMEOU_MS,
			NotLoadCacheAtStart: true,
			LogDir:              setup.NACOS_LOG_DIR + "/" + grpcApp.App.Name,
			CacheDir:            setup.NACOS_CACHE_DIR + "/" + grpcApp.App.Name,
			LogLevel:            setup.NACOS_LOG_LEVEL,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(config.GetNacosAddress(), config.GetNacosEndpoint()),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	kratosServer := server.New(
		server.ID(serverId),
		server.Name(grpcApp.App.Name),
		server.Version(grpcApp.App.Version),
		server.Metadata(map[string]string{}),
		server.Logger(logger),
		server.Server(httpSrv, grpcSrv),
		server.Registrar(nacos.New(client)),
	)

	if err = kratosServer.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
