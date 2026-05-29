package config

import (
	"aim/app/groupservice/model"
	commonconfig "aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	newConfig := model.DBConfig{}

	newConfig.Mysql = commonconfig.GetMysqlConfig(data)
	newConfig.Redis = commonconfig.GetRedisConfig(data)

	return newConfig
}
