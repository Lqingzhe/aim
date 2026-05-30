package handler

import (
	"aim/app/messageservice/model"
	"aim/commonmodel"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type KitexMessageServiceImpl struct {
	logger           *zap.Logger
	messageConfig    commonmodel.MessageConfig
	snowFlake        *snowflake.Node
	dbContext        *model.DBContext
	serviceClient    model.ServiceClient
	messageTopic     sarama.SyncProducer
	groupNoticeTopic sarama.SyncProducer
	systemTopic      sarama.SyncProducer
}

func NewMessageServiceImpl(logger *zap.Logger, messageConfig commonmodel.MessageConfig, snowFlake *snowflake.Node, dbContext *model.DBContext, serviceClient model.ServiceClient, messageTopic sarama.SyncProducer, groupNoticeTopic sarama.SyncProducer, systemTopic sarama.SyncProducer) *KitexMessageServiceImpl {
	return &KitexMessageServiceImpl{
		logger:           logger,
		messageConfig:    messageConfig,
		snowFlake:        snowFlake,
		dbContext:        dbContext,
		serviceClient:    serviceClient,
		messageTopic:     messageTopic,
		groupNoticeTopic: groupNoticeTopic,
		systemTopic:      systemTopic,
	}
}
