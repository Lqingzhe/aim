package config

import (
	"aim/app/api/model"
	"aim/pkg/config"
)

func InitConfig() *model.Config {
	data := commonconfig.OpenYaml()

	newConfig := &model.Config{}

	newConfig.CommonConfig = commonconfig.GetCommonConfig(data)
	newConfig.GatewayConfig = commonconfig.GetGatewayConfig(data)
	newConfig.DBConfig = initDBConfig(data)
	newConfig.LimiterConfig = commonconfig.GetLimitersConfig(data)
	newConfig.TokenConfig = commonconfig.GetTokenConfig(data)
	newConfig.KafkaConfig = commonconfig.GetKafkaConfig(data)
	return newConfig
}
