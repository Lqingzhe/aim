package config

import (
	"aim/app/aiservice/model"
	commonconfig "aim/pkg/config"
)

func InitConfig() *model.Config {
	newStruct := &model.Config{}
	data := commonconfig.OpenYaml()
	newStruct.CommonConfig = commonconfig.GetCommonConfig(data)
	newStruct.ServiceConfig = commonconfig.GetServiceConfig(data)
	newStruct.KafkaConfig = commonconfig.GetKafkaConfig(data)
	newStruct.DBConfig = initDBConfig(data)
	newStruct.AiConfig = commonconfig.GetAiConfig(data)
	newStruct.LimiterConfig = commonconfig.GetLimitersConfig(data)

	return newStruct
}
