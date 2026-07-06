package service

import (
	"aim/app/aiservice/dao"
	"aim/app/aiservice/dao/userprofile"
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"context"
	"fmt"
)

type UserProfile struct {
	traceID   string
	dbContext *model.DBContext
}

func NewUserProfile(traceID string, dbContext *model.DBContext) *UserProfile {
	return &UserProfile{
		traceID:   traceID,
		dbContext: dbContext,
	}
}
func (u *UserProfile) GetUserProfile(ctx context.Context, userID int64) (profile string, exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("userprofile:GetUserProfile")
	var finalErr error
	userProfileStruct := userprofile.NewStruct(userID, "")
	exist, err = dao.Get(ctx, userProfileStruct, u.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return "", false, finalErr
	}
	if exist {
		return "", false, finalErr
	}
	return userProfileStruct.Info.Profile, true, finalErr
}
func (u *UserProfile) UpdateUserProfile(ctx context.Context, userID int64, profile string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("userprofile:UpdateUserProfile")
	var finalErr error
	if profile == "" {
		return newerror.MakeError(-1, newerror.CodeMissingParam, "", fmt.Errorf("Lack User Profile"), newerror.LevelInfo)
	}
	userProfileStruct := userprofile.NewStruct(userID, profile)
	exist, err := dao.Update(ctx, userProfileStruct, u.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return finalErr
	}
	if !exist {
		err = dao.Add(ctx, userProfileStruct, u.dbContext)
		if newerror.WhetherInterrupt(err, &finalErr) {
			return finalErr
		}
	}
	return finalErr
} //tool
