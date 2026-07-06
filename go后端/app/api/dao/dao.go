package dao

import (
	"aim/app/api/model"
	"aim/commonmodel"
	"aim/pkg/config"
	newlog "aim/pkg/log"
	"context"

	"go.uber.org/zap"
)

func InitDB(redisConfig *commonmodel.RedisConfig, logger *zap.Logger) *model.DBContext {
	Redis, err := commonconfig.MakeRedis(redisConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Redis Failed.")
	}
	return &model.DBContext{
		Redis: Redis,
	}
}
func CloseDB(dbContext *model.DBContext) {
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
