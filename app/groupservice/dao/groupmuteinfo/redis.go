package groupmuteinfo

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func setRedis(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:SetRedis")
	MuteEndTime := info.MuteEndTime
	if MuteEndTime.Second() == 0 {
		MuteEndTime = time.Now().Add(10 * time.Minute)
	}
	KEY := []string{fmt.Sprintf("group_mute_info:%d%d", info.GroupID, info.UserID)}
	VALUE := []any{int64(info.MuteEndTime.Sub(time.Now()).Seconds()), "mute_endtime", MuteEndTime.Unix(), "mute_reason", info.MuteReason}
	result := dbContext.Redis.Script[commonmodel.HSETEX].Run(ctx, dbContext.Redis.Client, KEY, VALUE)
	if result.Err() != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeCacheError, "Success", result.Err(), newerror.LevelWarn, newerror.WithContinueError)
	}
	return nil
}
func getRedis(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:GetRedis")
	result, err := dbContext.Redis.Client.HGetAll(ctx, fmt.Sprintf("group_mute_info:%d%d", info.GroupID, info.UserID)).Result()
	if err != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return false, err2
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeCacheError, "Success", err, newerror.LevelWarn, newerror.WithContinueError)
	}
	if result["mute_reason"] == "" {
		return false, nil
	}
	t1, err := strconv.ParseInt(result["mute_endtime"], 10, 64)
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeCacheError, "Success", err, newerror.LevelWarn, newerror.WithContinueError)
	}
	t2 := time.Unix(t1, 0)
	info.Info = append(info.Info, &model.GroupMuteInfo{
		MuteEndTime: t2,
		MuteReason:  result["mute_reason"],
	})
	return true, nil
}
func deleteRedis(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:DeleteRedis")
	result := dbContext.Redis.Client.Del(ctx, fmt.Sprintf("group_mute_info:%d%d", info.GroupID, info.UserID))
	if result.Err() != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeCacheError, "Success", result.Err(), newerror.LevelWarn, newerror.WithContinueError)
	}
	return nil
}
