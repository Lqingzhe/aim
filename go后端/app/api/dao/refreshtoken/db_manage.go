package refreshtoken

import (
	"aim/app/api/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"time"
)

type RefreshToken struct {
	model.TokenInfo
	Info *model.TokenInfo
}

func NewStruct(refreshToken string, userID int64, deviceID string, expireTime time.Duration) *RefreshToken {
	newStruct := &RefreshToken{
		TokenInfo: model.TokenInfo{
			RefreshTokenID: refreshToken,
			UserID:         userID,
			DeviceID:       deviceID,
			ExpireTime:     expireTime,
		},
	}
	return newStruct
}

func (t *RefreshToken) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setRedis(ctx, DB, t)
	return err
}
func (t *RefreshToken) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getRedis(ctx, DB, t)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (t *RefreshToken) UpdateInfo(_ context.Context, _ any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Useless Module Unexpectedly Used", fmt.Errorf("%s", "Useless Module Unexpectly Used"), newerror.LevelFatal)
}
func (t *RefreshToken) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return newerror.TranslateError(err).AddErrorTrace("db_manage:DeleteInfo")
	}
	err = deleteBlurryRedis(ctx, DB, t)
	return newerror.TranslateError(err).AddErrorTrace("db_manage:DeleteInfo")
}
