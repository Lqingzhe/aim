package config

import (
	"aim/app/userservice/model"
	"aim/pkg/config"
)

func InitConfig() model.Config {
	data := commonconfig.OpenYaml()
	newConfig := model.Config{}

	newConfig.CommonConfig = commonconfig.GetCommonConfig(data)
	newConfig.ServiceConfig = commonconfig.GetServiceConfig(data)
	newConfig.UserConfig = commonconfig.GetUserConfig(data)
	newConfig.DBConfig = initDBConfig(data)

	return newConfig
}
