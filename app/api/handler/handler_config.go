package handler

import (
	"aim/app/api/model"
	"aim/commonmodel"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
)

type HandlerConfig struct {
	snowNode          *snowflake.Node
	dbContext         *model.DBContext
	tokenConfig       commonmodel.TokenConfig
	serviceClient     model.ServiceClient
	websocketUpgrader websocket.Upgrader
	hub               *model.WebSockedHub
	consumer          sarama.Consumer
}

func NewHandlerConfig(snowNode *snowflake.Node, dbContext *model.DBContext, tokenConfig commonmodel.TokenConfig, serviceClient model.ServiceClient, websocketUpgrader websocket.Upgrader, consumer sarama.Consumer) *HandlerConfig {
	return &HandlerConfig{
		snowNode:          snowNode,
		dbContext:         dbContext,
		tokenConfig:       tokenConfig,
		serviceClient:     serviceClient,
		websocketUpgrader: websocketUpgrader,
		hub:               &model.WebSockedHub{Client: make(map[int64]map[string]*model.Client)},
		consumer:          consumer,
	}
}
