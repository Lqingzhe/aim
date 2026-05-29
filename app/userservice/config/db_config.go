package config

import (
	"aim/app/userservice/model"
	"aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	newConfig := model.DBConfig{}

	newConfig.Mysql = commonconfig.GetMysqlConfig(data)

	return newConfig
}
