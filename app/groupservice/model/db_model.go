package model

import "aim/commonmodel"

type DBConfig struct {
	Mysql commonmodel.MysqlConfig
	Redis commonmodel.RedisConfig
}

type DBContext struct {
	Mysql *commonmodel.MysqlContext
	Redis *commonmodel.RedisContext
}
