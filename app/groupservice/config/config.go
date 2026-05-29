package config

import (
	"aim/app/groupservice/model"
	commonconfig "aim/pkg/config"
)

func InitConfig() model.Config {
	data := commonconfig.OpenYaml()
	newConfig := model.Config{}

	newConfig.CommonConfig = commonconfig.GetCommonConfig(data)
	newConfig.ServiceConfig = commonconfig.GetServiceConfig(data)
	newConfig.GroupConfig = commonconfig.GetGroupConfig(data)
	newConfig.KafkaConfig = commonconfig.GetKafkaConfig(data)
	newConfig.DBConfig = initDBConfig(data)

	return newConfig
}
