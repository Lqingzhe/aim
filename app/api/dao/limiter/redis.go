package limiter

import (
	"aim/app/api/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"context"

	"github.com/redis/go-redis/v9"
)

func setRedis(ctx context.Context, dbContext *model.DBContext, Info *Limiter) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:SetRedis")
	KEY := []string{"limiter:" + fmt.Sprintf("%d%s", Info.UserID, Info.DeviceID)}
	ARGV := []any{Info.ExpireTime.Seconds(), "last_time", Info.LastTime, "last_tokens", Info.LastTokens}
	result := dbContext.Redis.Script[commonmodel.HSETEX].Run(ctx, dbContext.Redis.Client, KEY, ARGV)
	if result.Err() != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", result.Err(), newerror.LevelError)
	}
	return nil
}
func getRedis(ctx context.Context, dbContext *model.DBContext, Info *Limiter) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:GetRedis")
	result, err := dbContext.Redis.Client.HGetAll(ctx, "limiter:"+fmt.Sprintf("%d%s", Info.UserID, Info.DeviceID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	Info.Info = &model.LimiterInfo{
		UserID:   Info.UserID,
		DeviceID: Info.DeviceID,
	}
	Info.Info.LastTime, err = strconv.ParseInt(result["last_time"], 10, 64)
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Database Error", err, newerror.LevelFatal)
	}
	Info.Info.LastTokens, err = strconv.ParseInt(result["last_tokens"], 10, 64)
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Database Error", err, newerror.LevelFatal)
	}
	return true, nil
}
