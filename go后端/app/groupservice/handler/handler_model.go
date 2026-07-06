package handler

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type GroupServiceImpl struct {
	SnowNode      *snowflake.Node
	DBContext     *model.DBContext
	Logger        *zap.Logger
	GroupConfig   commonmodel.GroupConfig
	ServiceClient model.ServiceClient

	GroupNoticeTopic sarama.SyncProducer
	SystemTopic      sarama.SyncProducer
}

func NewGroupServiceImpl(SnowNode *snowflake.Node, DBContext *model.DBContext, Logger *zap.Logger, GroupConfig commonmodel.GroupConfig, ServiceClient model.ServiceClient, GroupNoticeTopic sarama.SyncProducer, SystemTopic sarama.SyncProducer) *GroupServiceImpl {
	return &GroupServiceImpl{
		SnowNode:         SnowNode,
		DBContext:        DBContext,
		Logger:           Logger,
		GroupConfig:      GroupConfig,
		ServiceClient:    ServiceClient,
		GroupNoticeTopic: GroupNoticeTopic,
		SystemTopic:      SystemTopic,
	}
}
