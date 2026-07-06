package commonconfig

import (
	"aim/commonmodel"
	newlog "aim/pkg/log"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
)

func initNacosClient(logger *zap.Logger, nacocAddr ...commonmodel.NacosConfig) naming_client.INamingClient {
	sc := make([]constant.ServerConfig, len(nacocAddr))
	for i, j := range nacocAddr {
		sc[i] = *constant.NewServerConfig(j.Host, uint64(j.Port))
	}
	cc := constant.ClientConfig{
		NamespaceId: "public",
	}
	nacosClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	})
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Nacos Service Register Error")
	}
	return nacosClient
}
func ResolverService(nacosConfig commonmodel.NacosConfig, logger *zap.Logger) client.Option {
	return client.WithResolver(
		resolver.NewNacosResolver(
			initNacosClient(
				logger,
				nacosConfig,
			),
		),
	)
}
func RegisterService(nacosConfig commonmodel.NacosConfig, logger *zap.Logger) server.Option {
	return server.WithRegistry(
		registry.NewNacosRegistry(
			initNacosClient(
				logger,
				nacosConfig,
			),
		),
	)
}
