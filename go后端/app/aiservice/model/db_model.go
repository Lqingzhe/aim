package model

import "aim/commonmodel"

type DBConfig struct {
	Mysql   commonmodel.MysqlConfig
	MongoDB commonmodel.MongoDBConfig
}
type DBContext struct {
	Mysql   *commonmodel.MysqlContext
	MongoDB *commonmodel.MongoDBContext
}
