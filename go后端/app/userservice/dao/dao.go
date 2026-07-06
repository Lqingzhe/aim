package dao

import (
	"aim/app/userservice/model"
	"aim/commonmodel"
	"aim/pkg/config"
	newlog "aim/pkg/log"

	"context"

	"go.uber.org/zap"
)

func InitDB(mysqlConfig *commonmodel.MysqlConfig, logger *zap.Logger) *model.DBContext {
	Mysql, err := commonconfig.MakeMysql(mysqlConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Mysql Failed")
	}

	return &model.DBContext{
		Mysql: Mysql,
	}
}
func CloseDB(dbContext *model.DBContext) {
	commonconfig.DBClose(dbContext.Mysql.Client)
}

func Add(ctx context.Context, info commonmodel.DBOperater, dbContext *model.DBContext) error {
	return info.AddInfo(ctx, dbContext)
}
func Update(ctx context.Context, info commonmodel.DBOperater, dbContext *model.DBContext) (bool, error) {
	return info.UpdateInfo(ctx, dbContext)
}
func Delete(ctx context.Context, info commonmodel.DBOperater, dbContext *model.DBContext) error {
	return info.DeleteInfo(ctx, dbContext)
}
func Get(ctx context.Context, info commonmodel.DBOperater, dbContext *model.DBContext) (bool, error) {
	return info.GetInfo(ctx, dbContext)
}
