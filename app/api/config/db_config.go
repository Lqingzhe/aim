package config

import (
	"aim/app/api/model"
	"aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	newConfig := model.DBConfig{}

	newConfig.Redis = commonconfig.GetRedisConfig(data)

	return newConfig
}
