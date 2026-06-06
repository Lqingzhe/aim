package handler

import (
	"aim/app/fileservice/model"
	"aim/commonmodel"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

type KitexFileServiceImpl struct {
	logger     *zap.Logger
	fileConfig commonmodel.FileConfig
	snowFlake  *snowflake.Node
	dbContext  *model.DBContext
}

func NewFileServiceImpl(logger *zap.Logger, fileConfig commonmodel.FileConfig, snowFlake *snowflake.Node, dbContext *model.DBContext) *KitexFileServiceImpl {
	return &KitexFileServiceImpl{
		logger:     logger,
		fileConfig: fileConfig,
		snowFlake:  snowFlake,
		dbContext:  dbContext,
	}
}
