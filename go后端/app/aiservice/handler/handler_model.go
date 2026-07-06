package handler

import (
	myagent "aim/app/aiservice/agent"
	"aim/app/aiservice/model"
	"aim/commonmodel"

	"github.com/IBM/sarama"
	agenttool "github.com/cloudwego/eino/components/tool"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct {
	logger               *zap.Logger
	aiTopic              sarama.SyncProducer
	dbContext            *model.DBContext
	aiConfig             commonmodel.AiConfig
	serviceClient        model.ServiceClient
	consumer             sarama.Consumer
	Limiter              *rate.Limiter
	TraceWithUserManager *myagent.TraceWithUserManager
	tools                []agenttool.BaseTool
}

func NewAiServiceImpl(logger *zap.Logger, aiTopic sarama.SyncProducer, dbContext *model.DBContext, aiConfig commonmodel.AiConfig, serviceClient model.ServiceClient, consumer sarama.Consumer, TraceWithUserManager *myagent.TraceWithUserManager, Limiter *rate.Limiter, tools []agenttool.BaseTool) *AiServiceImpl {
	return &AiServiceImpl{
		dbContext:            dbContext,
		aiConfig:             aiConfig,
		logger:               logger,
		aiTopic:              aiTopic,
		serviceClient:        serviceClient,
		consumer:             consumer,
		Limiter:              Limiter,
		TraceWithUserManager: TraceWithUserManager,
		tools:                tools,
	}
}

func (s *AiServiceImpl) BeginConsumer() {
	s.SendMessageToUser(75)

}
