package commonconfig

import (
	"aim/commonmodel"
	"context"
	"fmt"
	"log"
	"time"

	"aim/pkg/log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MakeMongoDB(cfg *commonmodel.MongoDBConfig) (*commonmodel.MongoDBContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Url))
	if err != nil {
		return nil, fmt.Errorf("MongoDB Connect Error: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("MongoDB Ping Error: %v", err)
	}
	log.Println("MongoDB Connect Success")
	return &commonmodel.MongoDBContext{
		Client: client,
	}, nil
}
func MakeMysql(cfg *commonmodel.MysqlConfig) (*commonmodel.MysqlContext, error) {
	client, err := gorm.Open(mysql.Open(cfg.Url), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("Mysql Connect Error: %v", err)
	}
	config, err := client.DB()
	if err != nil {
		return nil, fmt.Errorf("Mysql DB Error: %v", err)
	}
	config.SetMaxIdleConns(cfg.SetMaxIdleConns)
	config.SetMaxOpenConns(cfg.SetMaxOpenConns)
	config.SetConnMaxLifetime(cfg.SetConnMaxLifetime)
	log.Println("Mysql Connect Success")
	return &commonmodel.MysqlContext{
		Client: client,
	}, nil
}
func AutoMysql(client *commonmodel.MysqlContext, DBModel ...commonmodel.DataModel) []error {
	var errs []error
	var err error
	for _, model := range DBModel {
		err = client.Client.AutoMigrate(model)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func MakeRedis(cfg *commonmodel.RedisConfig) (*commonmodel.RedisContext, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              0,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("Redis Connect Error: %v", err)
	}
	log.Println("Redis Connect Success")
	return &commonmodel.RedisContext{
		Client: client,
		Script: commonmodel.InitRedisScript(),
	}, nil
}
func RestartRedis(redisContext *commonmodel.RedisContext, cfg *commonmodel.RedisConfig, Logger *zap.Logger) {
	timer := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
			}
			if redisContext.Client != nil {
				err := redisContext.Client.Ping(context.Background()).Err()
				if err == nil {
					continue
				}
				newlog.LogInitWarn(Logger, err, "Redis Connect Interrupt")
				_ = redisContext.Client.Close()
				redisContext.Client = nil
			}
			newClient := redis.NewClient(&redis.Options{
				Addr:            cfg.Addr,
				Password:        cfg.Password,
				DB:              0,
				PoolSize:        cfg.PoolSize,
				MinIdleConns:    cfg.MinIdleConns,
				MaxIdleConns:    cfg.MaxIdleConns,
				DialTimeout:     cfg.DialTimeout,
				ReadTimeout:     cfg.ReadTimeout,
				WriteTimeout:    cfg.WriteTimeout,
				ConnMaxLifetime: cfg.ConnMaxLifetime,
				ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			})
			if err := newClient.Ping(context.Background()).Err(); err == nil {
				redisContext.Client = newClient
				newlog.LogInitInfo(Logger, "Redis Reconnect Success")
			} else {
				newlog.LogInitWarn(Logger, err, "Redis Reconnect Error")
			}
		}
	}()
}
func DBClose[T *gorm.DB | *redis.Client | *mongo.Client](db T) {
	switch v := any(db).(type) {
	case *gorm.DB:
		cfg, _ := v.DB()
		_ = cfg.Close()
	case *redis.Client:
		_ = v.Close()
	case *mongo.Client:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = v.Disconnect(ctx)
		cancel()
	}
}
