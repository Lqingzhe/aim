package model

import "aim/commonmodel"

type DBConfig struct {
	MongoDB commonmodel.MongoDBConfig
	Mysql   commonmodel.MysqlConfig
	Redis   commonmodel.RedisConfig
}
type DBContext struct {
	MongoDB *commonmodel.MongoDBContext
	Mysql   *commonmodel.MysqlContext
	Redis   *commonmodel.RedisContext
}
