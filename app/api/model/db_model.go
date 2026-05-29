package model

import "aim/commonmodel"

type DBConfig struct {
	Redis commonmodel.RedisConfig
}
type DBContext struct {
	Redis *commonmodel.RedisContext
}
