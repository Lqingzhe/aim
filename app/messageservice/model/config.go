package model

import "aim/commonmodel"

type Config struct {
	commonmodel.CommonConfig
	commonmodel.MessageConfig
	commonmodel.ServiceConfig
	commonmodel.KafkaConfig
	DBConfig
}
