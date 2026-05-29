package config

import (
	"aim/app/fileservice/model"
	commonconfig "aim/pkg/config"
)

func InitConfig() model.Config {
	data := commonconfig.OpenYaml()
	return model.Config{
		CommonConfig:  commonconfig.GetCommonConfig(data),
		ServiceConfig: commonconfig.GetServiceConfig(data),
		DBConfig:      initDBConfig(data),
		FileConfig:    commonconfig.GetFileConfig(data),
	}
}
