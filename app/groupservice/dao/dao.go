package dao

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"
	commonconfig "aim/pkg/config"
	newlog "aim/pkg/log"
	"context"

	"go.uber.org/zap"
)

func InitDB(MysqlConfig *commonmodel.MysqlConfig, RedisConfig *commonmodel.RedisConfig, logger *zap.Logger) *model.DBContext {
	MysqlCtx, err := commonconfig.MakeMysql(MysqlConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Mysql Failed")
	}

	RedisCtx, err := commonconfig.MakeRedis(RedisConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Redis Failed")
	}
	return &model.DBContext{
		Mysql: MysqlCtx,
		Redis: RedisCtx,
	}
}
func CloseDB(dbContext *model.DBContext) {
	commonconfig.DBClose(dbContext.Mysql.Client)
	commonconfig.DBClose(dbContext.Redis.Client)
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
