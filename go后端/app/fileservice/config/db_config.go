package config

import (
	"aim/app/fileservice/model"
	commonconfig "aim/pkg/config"
)

func initDBConfig(data []byte) model.DBConfig {
	return model.DBConfig{
		Mysql: commonconfig.GetMysqlConfig(data),
	}
}
