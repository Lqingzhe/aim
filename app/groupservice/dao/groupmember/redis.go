package groupmember

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func translateRoleToFloat(role commonmodel.GroupRole) float64 {
	switch role {
	case commonmodel.Member:
		return 1
	case commonmodel.Manager:
		return 2
	case commonmodel.GroupOwner:
		return 3
	}
	return 0
}
func translateFloatToRole(roleNumber float64) commonmodel.GroupRole {
	switch roleNumber {
	case 1:
		return commonmodel.Member
	case 2:
		return commonmodel.Manager
	case 3:
		return commonmodel.GroupOwner
	}
	return ""
}
func setRedis(ctx context.Context, dbContext *model.DBContext, info *GroupMember) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:AddRedis")
	pipe := dbContext.Redis.Client.TxPipeline()

	roleValue := make([]redis.Z, 0, len(info.members))
	visitTimeValue := make([]redis.Z, 0, len(info.members))
	for i := range info.members {
		roleValue = append(roleValue, redis.Z{
			Member: info.members[i],
			Score:  translateRoleToFloat(info.Role[i]),
		})
		visitTimeValue = append(visitTimeValue, redis.Z{
			Member: info.members[i],
			Score:  float64(info.LastReadTime[i].Unix()),
		})
	}
	pipe.ZAdd(ctx, "group_member:role:"+strconv.FormatInt(info.GroupID, 10), roleValue...)
	pipe.ZAdd(ctx, "group_member:visit_time:"+strconv.FormatInt(info.GroupID, 10), visitTimeValue...)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func deleteRedisWithGroup(ctx context.Context, dbContext *model.DBContext, info *GroupMember) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:DeleteRedisWithGroup")
	pipe := dbContext.Redis.Client.TxPipeline()
	pipe.Del(ctx, "group_member:role:"+strconv.FormatInt(info.GroupID, 10))
	pipe.Del(ctx, "group_member:visit_time:"+strconv.FormatInt(info.GroupID, 10))
	_, err = pipe.Exec(ctx)
	if err != nil {
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func deleteRedisWithMember(ctx context.Context, dbContext *model.DBContext, info *GroupMember) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:DeleteRedisWithMember")
	pipe := dbContext.Redis.Client.TxPipeline()
	members := make([]any, 0, len(info.members))
	for i := range info.members {
		members = append(members, info.members[i])
	}
	pipe.ZRem(ctx, "group_member:role:"+strconv.FormatInt(info.GroupID, 10), members)
	pipe.ZRem(ctx, "group_member:visit_time:"+strconv.FormatInt(info.GroupID, 10), members)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	return nil
}
func getRedisWithUser(ctx context.Context, dbContext *model.DBContext, info *GroupMember) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:getRedisWithUser")
	members := make([]string, len(info.members))
	for i := range info.members {
		members[i] = strconv.FormatInt(info.members[i], 10)
	}
	roleScores, err := dbContext.Redis.Client.ZMScore(ctx, "group_member:role:"+strconv.FormatInt(info.GroupID, 10), members...).Result()
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	visitTime, err := dbContext.Redis.Client.ZMScore(ctx, "group_member:visit_time:"+strconv.FormatInt(info.GroupID, 10), members...).Result()
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	isEmpty := false
	for i, roleScore := range roleScores {
		if !math.IsNaN(roleScore) {
			isEmpty = true
		}
		info.Info = append(info.Info, &GroupMemberInfo{
			GroupID:      info.GroupID,
			UserID:       info.members[i],
			Role:         translateFloatToRole(roleScore),
			LastReadTime: time.Unix(int64(visitTime[i]), 0),
		})
	}
	if !isEmpty {
		return false, nil
	}
	return true, nil
}
func getRedisWithGroup(ctx context.Context, dbContext *model.DBContext, info *GroupMember) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("redis:getRedisWithGroup")
	roleScores, err := dbContext.Redis.Client.ZRangeWithScores(ctx, "group_member:role:"+strconv.FormatInt(info.GroupID, 10), 0, -1).Result()
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	visitTime, err := dbContext.Redis.Client.ZRangeWithScores(ctx, "group_member:visit_time:"+strconv.FormatInt(info.GroupID, 10), 0, -1).Result()
	if err != nil {
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", err, newerror.LevelError)
	}
	if len(roleScores) == 0 || len(visitTime) == 0 {
		return false, nil
	}
	for i := range roleScores {
		userID, _ := strconv.ParseInt(roleScores[i].Member.(string), 10, 64)
		info.Info = append(info.Info, &GroupMemberInfo{
			GroupID:      info.GroupID,
			UserID:       userID,
			Role:         translateFloatToRole(roleScores[i].Score),
			LastReadTime: time.Unix(int64(visitTime[i].Score), 0),
		})
	}
	return true, nil
}
