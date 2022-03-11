package client_conn

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/liangguifeng/kratos-app/internal/config"
	"github.com/liangguifeng/kratos-app/internal/setup"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	google_grpc "google.golang.org/grpc"
)

// NewConn
func NewConn(serviceName string) (*google_grpc.ClientConn, error) {
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.GetNacosNamespaceId(),
			TimeoutMs:           setup.NACOS_TIMEOU_MS,
			NotLoadCacheAtStart: true,
			LogDir:              setup.NACOS_LOG_DIR + "/" + serviceName,
			CacheDir:            setup.NACOS_CACHE_DIR + "/" + serviceName,
			LogLevel:            setup.NACOS_LOG_LEVEL,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(config.GetNacosAddress(), config.GetNacosEndpoint()),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	dis := nacos.New(client)

	return grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///"+serviceName+".grpc"),
		grpc.WithDiscovery(dis),
	)
}
