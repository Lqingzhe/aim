package model

import (
	"aim/commonmodel"
)

type Config struct {
	commonmodel.CommonConfig
	commonmodel.GatewayConfig
	DBConfig
	commonmodel.LimiterConfig
	commonmodel.TokenConfig
	commonmodel.KafkaConfig
}
