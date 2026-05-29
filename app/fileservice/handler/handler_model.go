package handler

import (
	"aim/app/fileservice/model"
	"aim/commonmodel"

	"github.com/bwmarrin/snowflake"
)

type KitexFileServiceImpl struct {
	fileConfig commonmodel.FileConfig
	snowFlake  *snowflake.Node
	dbContext  *model.DBContext
}

func NewFileServiceImpl(fileConfig commonmodel.FileConfig, snowFlake *snowflake.Node, dbContext *model.DBContext) *KitexFileServiceImpl {
	return &KitexFileServiceImpl{
		fileConfig: fileConfig,
		snowFlake:  snowFlake,
		dbContext:  dbContext,
	}
}
