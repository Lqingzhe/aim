package handler

import (
	"aim/app/userservice/model"
	"aim/commonmodel"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type UserServiceImpl struct {
	SnowNode   *snowflake.Node
	DBContext  *model.DBContext
	Logger     *zap.Logger
	UserConfig commonmodel.UserConfig
	EquipID    int64
}

func NewUserServiceImpl(SnowNode *snowflake.Node, DBContext *model.DBContext, Logger *zap.Logger, UserConfig commonmodel.UserConfig, EquipID int64) *UserServiceImpl {
	return &UserServiceImpl{
		SnowNode:   SnowNode,
		DBContext:  DBContext,
		Logger:     Logger,
		UserConfig: UserConfig,
		EquipID:    EquipID,
	}
}
