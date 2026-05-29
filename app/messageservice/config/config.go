package config

import (
	"aim/app/messageservice/model"
	commonconfig "aim/pkg/config"
)

func InitConfig() model.Config {
	data := commonconfig.OpenYaml()
	return model.Config{
		CommonConfig:  commonconfig.GetCommonConfig(data),
		MessageConfig: commonconfig.GetMessageConfig(data),
		ServiceConfig: commonconfig.GetServiceConfig(data),
		KafkaConfig:   commonconfig.GetKafkaConfig(data),
		DBConfig:      initDBConfig(data),
	}

}
