package model

import "aim/commonmodel"

type DBConfig struct {
	Mysql commonmodel.MysqlConfig
}

type DBContext struct {
	Mysql *commonmodel.MysqlContext
}
