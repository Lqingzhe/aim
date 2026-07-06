package config

import (
	"aim/app/messageservice/model"
	commonconfig "aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	return model.DBConfig{
		Mysql:   commonconfig.GetMysqlConfig(data),
		MongoDB: commonconfig.GetMongoDBConfig(data),
		Redis:   commonconfig.GetRedisConfig(data),
	}
}
