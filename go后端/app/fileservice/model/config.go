package model

import "aim/commonmodel"

type Config struct {
	commonmodel.CommonConfig
	commonmodel.ServiceConfig
	DBConfig
	commonmodel.FileConfig
}
