package userinfo

import (
	"aim/app/userservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type DaoUserInfo struct {
	model.UserInfo
	Info *model.UserInfo
}
type operate func(*DaoUserInfo)

func NewStruct(UserID int64, Operats ...operate) *DaoUserInfo {
	newstruct := &DaoUserInfo{
		UserInfo: model.UserInfo{
			UserID: UserID,
		},
		Info: &model.UserInfo{},
	}
	for _, Operate := range Operats {
		Operate(newstruct)
	}
	return newstruct
}
func WithBirthday(year int64, mouth int64, day int64) operate {
	return func(info *DaoUserInfo) {
		info.BirthdayDay = day
		info.BirthdayMonth = mouth
		info.BirthdayYear = year
	}
}
func WithIntroduction(introduction string) operate {
	return func(info *DaoUserInfo) {
		info.Introduction = introduction
	}
}
func WithUserName(username string) operate {
	return func(info *DaoUserInfo) {
		info.UserName = username
	}
}

func (d *DaoUserInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, d)
	return err
}
func (d *DaoUserInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, d)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (d *DaoUserInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, d)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (d *DaoUserInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, d)
	if err != nil {
		return err
	}
	return nil
}
