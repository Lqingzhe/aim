package main

import (
	"aim/app/messageservice/config"
	"aim/app/messageservice/dao"
	"aim/app/messageservice/handler"
	"aim/app/messageservice/model"
	"aim/kitex_gen/kitexfileservice/kitexfileservice"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	commonconfig "aim/pkg/config"
	"aim/pkg/id"
	newlog "aim/pkg/log"
	"net"
	"strconv"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
)

func main() {

	Config := config.InitConfig()

	logger := newlog.InitLog(Config.Service, Config.EquipID)
	defer logger.Sync()

	snowNode := id.InitSnowNode(Config.EquipID, logger)

	dbContext := dao.InitDB(&Config.Mysql, &Config.MongoDB, &Config.Redis, logger)
	defer dao.CloseDB(dbContext)
	commonconfig.AutoMysql(dbContext.Mysql, &model.OfflineMessageInfo{})

	groupClient := kitexgroupservice.MustNewClient(
		"group_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["group_service"].KitexTimeOut),
	)
	fileClient := kitexfileservice.MustNewClient(
		"file_service",
		commonconfig.ResolverService(Config.NacosConfig, logger),
		client.WithRPCTimeout(Config.CommonConfig.ServiceInfo["file_service"].KitexTimeOut),
	)
	//aiClient := kitexaiservice.MustNewClient(
	//	"ai_service",
	//	commonconfig.ResolverService("ai_service", Config.CommonConfig.ServiceInfo, logger),
	//)

	kafkaProducerConfig := commonconfig.GetKafkaProducerConfig()
	groupNoticeTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)
	systemTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)
	newMessageTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)

	addr, err := net.ResolveTCPAddr("tcp", Config.ServiceConfig.ServiceAddr.Host+":"+strconv.FormatInt(Config.ServiceAddr.Port, 10))
	if err != nil {
		newlog.LogInitFatal(logger, err, "Make Addr Failed")
	}

	svr := kitexmessageservice.NewServer(
		handler.NewMessageServiceImpl(
			logger,
			Config.MessageConfig,
			snowNode,
			dbContext,
			model.ServiceClient{
				GroupService: groupClient,
				FileService:  fileClient,
				//AiService:    aiClient,
			},
			newMessageTopic,
			groupNoticeTopic,
			systemTopic,
		),
		server.WithServiceAddr(addr),
		//server.WithServerBasicInfo(
		//	&rpcinfo.EndpointBasicInfo{
		//		ServiceName: "user_service",
		//	},
		//),
		//server.WithServiceAddr(&net.TCPAddr{
		//	IP:   net.ParseIP(Config.ServiceConfig.Host),
		//	Port: int(Config.ServiceConfig.Port),
		//}),
		//commonconfig.RegisterService(
		//	Config.NacosConfig,
		//	logger,
		//),
	)
	err = svr.Run()
	if err != nil {
		newlog.LogInitFatal(logger, err, "Grcp Begin Error")
	}
}
