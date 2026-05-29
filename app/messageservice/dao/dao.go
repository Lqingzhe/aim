package dao

import (
	"aim/app/messageservice/dao/offlinemessage"
	"aim/app/messageservice/model"
	"aim/commonmodel"
	commonconfig "aim/pkg/config"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"

	"context"

	"go.uber.org/zap"
)

func InitDB(mysqlConfig *commonmodel.MysqlConfig, mongoDBConfig *commonmodel.MongoDBConfig, redifConfig *commonmodel.RedisConfig, logger *zap.Logger) *model.DBContext {
	dbContext := &model.DBContext{}
	Mysql, err := commonconfig.MakeMysql(mysqlConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Mysql Failed")
	}

	MongoDB, err := commonconfig.MakeMongoDB(mongoDBConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init MongoDB Failed")
	}
	Redis, err := commonconfig.MakeRedis(redifConfig)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Redis Failed")
	}
	dbContext.Mysql = Mysql
	dbContext.MongoDB = MongoDB
	dbContext.Redis = Redis

	ErrChan := make(chan error)
	offlinemessage.ClearMysql(dbContext, ErrChan)
	go func() {
		for {
			select {
			case err := <-ErrChan:
				err2 := newerror.TranslateError(err)
				newlog.Log(newlog.AddError(logger, err, err2.StatusCode), err2.LogLevel, "Database Error")
			}
		}
	}()
	return &model.DBContext{
		Mysql:   Mysql,
		MongoDB: MongoDB,
		Redis:   Redis,
	}
}

func CloseDB(dbContext *model.DBContext) {
	commonconfig.DBClose(dbContext.Mysql.Client)
	commonconfig.DBClose(dbContext.MongoDB.Client)
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
