package main

import (
	"aim/app/aiservice/agent"
	"aim/app/aiservice/config"
	"aim/app/aiservice/dao"
	"aim/app/aiservice/handler"
	"aim/app/aiservice/model"
	"aim/kitex_gen/kitexaiservice/kitexaiservice"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	commonconfig "aim/pkg/config"
	newlog "aim/pkg/log"
	"net"
	"strconv"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"golang.org/x/time/rate"
)

func main() {
	Config := config.InitConfig()

	logger := newlog.InitLog(Config.Service, Config.EquipID)
	defer logger.Sync()

	dbContext := dao.InitDB(&Config.DBConfig.Mysql, &Config.MongoDB, logger)
	defer dao.CloseDB(dbContext)
	commonconfig.AutoMysql(
		dbContext.Mysql,
		&model.BotConfig{},
	)

	GroupClient := kitexgroupservice.MustNewClient(
		"group_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8890"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["group_service"].KitexTimeOut),
	)
	MessageClient := kitexmessageservice.MustNewClient(
		"message_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		//client.WithHostPorts("127.0.0.1:8891"),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["message_service"].KitexTimeOut),
	)
	kafkaProducerConfig := commonconfig.GetKafkaProducerConfig()
	aiTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)

	kafkaConfig := commonconfig.GetKafkaConsumerConfig()
	consumer := commonconfig.MakeKafkaConsumer(Config.KafkaConfig.Broker, kafkaConfig, logger)

	TraceWithUserManager := agent.NewTraceWithUserManager()
	tools := agent.InitTools(logger, TraceWithUserManager)

	addr, err := net.ResolveTCPAddr("tcp", Config.ServiceConfig.ServiceAddr.Host+":"+strconv.FormatInt(Config.ServiceConfig.ServiceAddr.Port, 10))
	if err != nil {
		newlog.LogInitFatal(logger, err, "Make Addr Failed")
	}
	handlerStruct := handler.NewAiServiceImpl(
		logger,
		aiTopic,
		dbContext,
		Config.AiConfig,
		model.ServiceClient{
			GroupService:   GroupClient,
			MessageService: MessageClient,
		},
		consumer,
		TraceWithUserManager,
		rate.NewLimiter(rate.Limit(Config.LimiterConfig.GenerateToken), int(Config.LimiterConfig.MaxToken)),
		tools,
	)

	handlerStruct.BeginConsumer()

	svr := kitexaiservice.NewServer(
		handlerStruct,
		server.WithServiceAddr(addr),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "ai_service",
			},
		),
		commonconfig.RegisterService(
			Config.NacosConfig,
			logger,
		),
	)
	err = svr.Run()
	if err != nil {
		newlog.LogInitFatal(logger, err, "Grcp Begin Error")
	}
}
