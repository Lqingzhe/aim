package model

import "aim/commonmodel"

type Config struct {
	commonmodel.CommonConfig
	commonmodel.ServiceConfig
	commonmodel.KafkaConfig
	DBConfig
	commonmodel.AiConfig
	commonmodel.LimiterConfig
}
