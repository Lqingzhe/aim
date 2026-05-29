package userlogin

import (
	"aim/app/userservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type Login struct {
	model.UserLoginInfo
	Info *model.UserLoginInfo
}

func NewStruct(UserID int64, Password string, Salt string) *Login {
	return &Login{
		UserLoginInfo: model.UserLoginInfo{
			UserID:   UserID,
			Password: Password,
			Salt:     Salt,
		},
		Info: &model.UserLoginInfo{},
	}
}

func (l *Login) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, l)
	if err != nil {
		return err
	}
	return nil
}
func (l *Login) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, l)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (l *Login) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, l)
	if err != nil {
		return err
	}
	return nil
}
func (l *Login) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, l)
	if err != nil {
		return false, err
	}
	return exist, nil
}
