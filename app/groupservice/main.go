package main

import (
	"aim/app/groupservice/config"
	"aim/app/groupservice/dao"
	"aim/app/groupservice/handler"
	"aim/app/groupservice/model"
	"aim/kitex_gen/kitexgroupservice/kitexgroupservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"

	commonconfig "aim/pkg/config"
	"aim/pkg/id"
	newlog "aim/pkg/log"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {

	Config := config.InitConfig()

	logger := newlog.InitLog(Config.Service, Config.EquipID)
	defer logger.Sync()

	snowNode := id.InitSnowNode(Config.EquipID, logger)

	dbContext := dao.InitDB(&Config.DBConfig.Mysql, &Config.DBConfig.Redis, logger)
	defer dao.CloseDB(dbContext)
	commonconfig.AutoMysql(
		dbContext.Mysql,
		&model.GroupMuteInfo{},
		&model.GroupApplyInfo{},
		&model.GroupWithUserInfo{},
		&model.GroupInfo{},
		&model.SessionInfo{},
	)

	kafkaProducerConfig := commonconfig.GetKafkaProducerConfig()
	groupNoticeTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)
	systemTopic := commonconfig.MakeKafkaProducer(Config.KafkaConfig.Broker, kafkaProducerConfig, logger)

	MessageClient := kitexmessageservice.MustNewClient(
		"message_service",
		commonconfig.ResolverService("message_service", Config.CommonConfig.ServiceInfo, logger),
	)

	svr := kitexgroupservice.NewServer(
		handler.NewGroupServiceImpl(
			snowNode,
			dbContext,
			logger,
			Config.GroupConfig,
			model.ServiceClient{
				MessageClient: MessageClient,
			},
			int64(Config.EquipID),
			groupNoticeTopic,
			systemTopic,
		),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "group-service",
			},
		),
		commonconfig.RegisterService(
			Config.ServiceConfig,
			logger,
		),
	)
	err := svr.Run()
	if err != nil {
		newlog.LogInitFatal(logger, err, "Grcp Begin Error")
	}
}
