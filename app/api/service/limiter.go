package service

import (
	"aim/app/api/dao"
	"aim/app/api/dao/limiter"
	"aim/app/api/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"context"
	"time"
)

type Limiter struct {
	dbContext     *model.DBContext
	limiterConfig commonmodel.LimiterConfig
}

func NewLimiter(dbContext *model.DBContext, config commonmodel.LimiterConfig) *Limiter {
	return &Limiter{
		dbContext:     dbContext,
		limiterConfig: config,
	}
}
func (l *Limiter) Exceed(ctx context.Context, userID int64, deviceID string) (needLimit bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("limiter:Exceed")
	Limit := limiter.NewStruct(userID, deviceID, 0, 0)
	exist, err := dao.Get(ctx, Limit, l.dbContext)
	if err != nil {
		return true, err
	}
	if !exist {
		Limit = limiter.NewStruct(userID, deviceID, time.Now().Unix(), l.limiterConfig.RedisExpireTime)
		Limit.LastTokens = l.limiterConfig.MaxToken
		err = dao.Add(ctx, Limit, l.dbContext)
		if err != nil {
			return true, err
		}
		return false, nil
	}
	Limit.LimiterInfo = *Limit.Info
	Tokens := Limit.LastTokens + l.limiterConfig.GenerateToken*(time.Now().Unix()-Limit.LastTime) - 1
	if Tokens < 0 {
		return true, nil
	}
	if Tokens > l.limiterConfig.MaxToken {
		Tokens = l.limiterConfig.MaxToken
	}
	Limit.LastTokens = Tokens
	Limit.ExpireTime = l.limiterConfig.RedisExpireTime
	Limit.LastTime = time.Now().Unix()
	_, err = dao.Update(ctx, Limit, l.dbContext)
	if err != nil {
		return true, err
	}
	return false, nil
	//先GET,不存在就SET，返回false
	//存在就检验，直接取GET的Info字段，加上ExpireTime后直接Update刷新时间，需要限流就返回true，不需要就返回false
}
