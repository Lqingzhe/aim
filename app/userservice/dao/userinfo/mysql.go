package userinfo

import (
	"aim/app/userservice/model"
	newerror "aim/pkg/error"
	"context"
	"errors"

	"gorm.io/gorm"
)

func setMysql(ctx context.Context, dbContext *model.DBContext, info *DaoUserInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.UserInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *DaoUserInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.UserInfo{}).Where("user_id = ?", info.UserInfo.UserID).First(info.Info)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *DaoUserInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.UserInfo{}).Where("user_id = ?", info.UserInfo.UserID).Updates(info.UserInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *DaoUserInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("user_id = ?", info.UserID).Delete(info.UserInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
