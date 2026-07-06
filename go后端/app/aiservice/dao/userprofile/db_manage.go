package userprofile

import (
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type UserProfile struct {
	model.UserProfile
	Info *model.UserProfile
}

func NewStruct(userID int64, profile string) *UserProfile {
	return &UserProfile{
		UserProfile: model.UserProfile{
			UserID:  userID,
			Profile: profile,
		},
		Info: &model.UserProfile{},
	}
}
func (u *UserProfile) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMongo(ctx, DB, u)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserProfile) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMongo(ctx, DB, u)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (u *UserProfile) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMongo(ctx, DB, u)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	err = setMongo(ctx, DB, u)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (u *UserProfile) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMongo(ctx, DB, u)
	if err != nil {
		return err
	}
	return nil
}
