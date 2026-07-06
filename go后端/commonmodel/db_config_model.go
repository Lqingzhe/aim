package commonmodel

import (
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

//Mysql

type MysqlConfig struct {
	Url      string `yaml:"-"`
	Username string
	Password string
	Host     string
	Port     int
	DBName   string

	SetMaxIdleConns    int
	SetMaxOpenConns    int
	SetConnMaxLifetime time.Duration
	Timeout            time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
}
type MysqlContext struct {
	Client *gorm.DB
}

//MongoDB

type MongoDBConfig struct {
	Url string
}
type MongoDBContext struct {
	Client *mongo.Client
}

//Redis

type LuaOperate string
type RedisConfig struct {
	Addr            string
	Password        string
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}
type RedisContext struct {
	Client *redis.Client
	Script map[LuaOperate]*redis.Script
}
