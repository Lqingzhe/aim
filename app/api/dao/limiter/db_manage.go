package limiter

import (
	"aim/app/api/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"time"

	"context"
	"fmt"
	"net/http"
)

type Limiter struct {
	model.LimiterInfo
	Info *model.LimiterInfo
}

func NewStruct(UserID int64, Device string, LastTime int64, ExpireTime time.Duration) *Limiter {
	return &Limiter{
		LimiterInfo: model.LimiterInfo{
			UserID:     UserID,
			DeviceID:   Device,
			LastTime:   LastTime,
			ExpireTime: ExpireTime},
	}
}

func (l *Limiter) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setRedis(ctx, DB, l)
	return err
}
func (l *Limiter) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getRedis(ctx, DB, l)
	return exist, err
}
func (l *Limiter) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	err = setRedis(ctx, DB, l)
	return exist, err
}
func (l *Limiter) DeleteInfo(_ context.Context, _ any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	return newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Useless Module Unexpectedly Used", fmt.Errorf("%s", "Useless Module Unexpectly Used"), newerror.LevelFatal)
}
