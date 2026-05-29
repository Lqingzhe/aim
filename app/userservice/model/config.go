package model

import "aim/commonmodel"

type Config struct {
	commonmodel.CommonConfig
	commonmodel.ServiceConfig
	commonmodel.UserConfig
	DBConfig
}
