package config

import (
	"aim/app/aiservice/model"
	commonconfig "aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	newConfig := model.DBConfig{}

	newConfig.MongoDB = commonconfig.GetMongoDBConfig(data)
	newConfig.Mysql = commonconfig.GetMysqlConfig(data)

	return newConfig
}
