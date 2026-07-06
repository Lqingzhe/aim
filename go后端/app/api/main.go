package main

import (
	api "aim/app/api/api"
	"aim/app/api/config"
	"aim/app/api/dao"
	"aim/app/api/model"
	"aim/kitex_gen/kitexaiservice/kitexaiservice"
	"aim/kitex_gen/kitexfileservice/kitexfileservice"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	"aim/kitex_gen/kitexuserservice/kitexuserservice"
	"net/http"

	commonconfig "aim/pkg/config"
	"aim/pkg/id"
	"aim/pkg/log"

	"github.com/cloudwego/kitex/client"
	"github.com/gorilla/websocket"
)

func main() {
	Config := config.InitConfig()

	logger := newlog.InitLog(Config.Service, Config.EquipID)
	defer logger.Sync()

	snowNode := id.InitSnowNode(Config.EquipID, logger)

	dbContext := dao.InitDB(&Config.DBConfig.Redis, logger)
	defer dao.CloseDB(dbContext)

	UserClient := kitexuserservice.MustNewClient(
		"user_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8889"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["user_service"].KitexTimeOut),
	)
	GroupClient := kitexgroupservice.MustNewClient(
		"group_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8890"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["group_service"].KitexTimeOut),
	)
	FileClient := kitexfileservice.MustNewClient(
		"file_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8892"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["file_service"].KitexTimeOut),
	)
	MessageClient := kitexmessageservice.MustNewClient(
		"message_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8891"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["message_service"].KitexTimeOut),
	)
	AiClient := kitexaiservice.MustNewClient(
		"ai_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8891"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["ai_service"].KitexTimeOut),
	)

	kafkaConfig := commonconfig.GetKafkaConsumerConfig()
	consumer := commonconfig.MakeKafkaConsumer(Config.KafkaConfig.Broker, kafkaConfig, logger)

	httpStruct := api.NewConfig(
		logger,
		snowNode,
		dbContext,
		Config.LimiterConfig,
		Config.TokenConfig,
		int64(Config.EquipID),
		Config.RoutTimeOut,
		model.ServiceClient{
			UserClient:    UserClient,
			GroupClient:   GroupClient,
			FileClient:    FileClient,
			MessageClient: MessageClient,
			AiClient:      AiClient,
		},
		websocket.Upgrader{
			ReadBufferSize:  Config.ReadBufferSize,
			WriteBufferSize: Config.WriteBufferSize,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		consumer)

	httpStruct.Begin(Config.Port)
}
