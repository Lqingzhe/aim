package refreshtoken

import (
	"aim/app/api/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func setRedis(ctx context.Context, dbContext *model.DBContext, Info *RefreshToken) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:SetRedis")
	KEY := []string{"refresh_token:" + Info.RefreshTokenID}
	ARGV := []any{Info.ExpireTime.Seconds(), "user_id", Info.UserID, "device_id", Info.DeviceID}
	result := dbContext.Redis.Script[commonmodel.HSETEX].Run(ctx, dbContext.Redis.Client, KEY, ARGV...)
	if result.Err() != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", result.Err(), newerror.LevelError)
	}
	return nil
} //向Redis中添加以RefreshToken为key的HASH

func deleteBlurryRedis(ctx context.Context, dbContext *model.DBContext, Info *RefreshToken) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:DeleteBlurryRedis")
	KEY := []string{"refresh_token:" + Info.RefreshTokenID}
	result := dbContext.Redis.Script[commonmodel.DELBLURRY].Run(ctx, dbContext.Redis.Client, KEY)
	if result.Err() != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", result.Err(), newerror.LevelError)
	}
	return nil
}
func getRedis(ctx context.Context, dbContext *model.DBContext, Info *RefreshToken) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:GetRedis")
	result, err := dbContext.Redis.Client.HGetAll(ctx, "refresh_token:"+Info.RefreshTokenID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	userID, err := strconv.ParseInt(result["user_id"], 10, 64)
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Database Error", err, newerror.LevelFatal)
	}
	Info.Info = &model.TokenInfo{
		UserID:   userID,
		DeviceID: result["device_id"],
	}
	return true, nil
}
