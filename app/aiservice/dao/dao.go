package dao

import (
	"aim/app/aiservice/model"
	"aim/commonmodel"
	commonconfig "aim/pkg/config"
	newlog "aim/pkg/log"
	"context"

	"go.uber.org/zap"
)

func InitDB(mysqlConfig *commonmodel.MysqlConfig, mongoConfig *commonmodel.MongoDBConfig, logger *zap.Logger) *model.DBContext {
	Mysql, err := commonconfig.MakeMysql(mysqlConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Mysql Failed")
	}
	MongoDB, err := commonconfig.MakeMongoDB(mongoConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init MongoDB Failed")
	}
	return &model.DBContext{
		Mysql:   Mysql,
		MongoDB: MongoDB,
	}
}
func CloseDB(dbContext *model.DBContext) {
	commonconfig.DBClose(dbContext.Mysql.Client)
	commonconfig.DBClose(dbContext.MongoDB.Client)
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
