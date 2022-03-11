package helper

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liangguifeng/kratos-app/internal/config"
	"github.com/liangguifeng/kratos-app/internal/setup"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type registerClient struct {
	client naming_client.INamingClient
}

// NewRegisterConn
func NewRegisterConn() (*registerClient, error) {
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.GetNacosNamespaceId(),
			TimeoutMs:           setup.NACOS_TIMEOU_MS,
			NotLoadCacheAtStart: true,
			LogDir:              setup.NACOS_LOG_DIR,
			CacheDir:            setup.NACOS_CACHE_DIR,
			LogLevel:            setup.NACOS_LOG_LEVEL,
		},
		ServerConfigs: []constant.ServerConfig{
			*constant.NewServerConfig(config.GetNacosAddress(), config.GetNacosEndpoint()),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	return &registerClient{client: client}, err
}

// GetClient
func (r registerClient) GetClient() naming_client.INamingClient {
	return r.client
}
