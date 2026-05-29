package offlinemessage

import (
	"aim/app/messageservice/model"
	newerror "aim/pkg/error"
	"context"
	"net/http"
	"strconv"
	"time"
)

func setRedis(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:SetRedis")
	pipe := dbContext.Redis.Client.TxPipeline()
	for i := range len(info.UserAndDeviceID)/75 + 1 {
		startSize := i * 75
		var endSize int
		if i*75 > len(info.UserAndDeviceID) {
			endSize = len(info.UserAndDeviceID) - 1
		} else {
			endSize = (i+1)*75 - 1
		}
		for j := startSize; j <= endSize; j++ {
			pipe.SAdd(ctx, "offset_message:"+info.UserAndDeviceID[j], info.MessageID)
			pipe.Expire(ctx, "offset_message:"+info.UserAndDeviceID[j], 7*24*time.Hour)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		if isContextErr, err2 := newerror.IsContextError(err); isContextErr {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func getRedis(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:GetRedis")
	var cursor uint64
	info.Info = make([]*model.OfflineMessageInfo, 0)
	for {
		var batch []string
		batch, cursor, err = dbContext.Redis.Client.SScan(ctx, "offset_message:"+info.UserAndDeviceID[0], cursor, "", 100).Result()
		if err != nil {
			info.Info = nil
			if isContextErr, err2 := newerror.IsContextError(err); isContextErr {
				return false, err2
			}
			return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
		}
		if cursor == 0 {
			break
		}
		for i := range batch {
			messageID, _ := strconv.ParseInt(batch[i], 10, 64)
			info.Info = append(info.Info, &model.OfflineMessageInfo{MessageID: messageID})
		}
	}
	return len(info.Info) > 0, nil
}
func deleteRedis(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:DeleteRedis")
	pipe := dbContext.Redis.Client.TxPipeline()
	for i := range info.UserAndDeviceID {
		pipe.Del(ctx, "offset_message:"+info.UserAndDeviceID[i])
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		if isContextErr, err2 := newerror.IsContextError(err); isContextErr {
			return err2
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
