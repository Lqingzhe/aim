package useraiconfig

import (
	"aim/app/aiservice/model"
	"aim/pkg/error"
	"context"
	"errors"

	"gorm.io/gorm"
)

func setMysql(ctx context.Context, dbContext *model.DBContext, info *UserAiConfig) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(&info.BotConfig)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *UserAiConfig) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.BotConfig{}).Where("user_id = ?", info.BotConfig.UserID).First(info.Info)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *UserAiConfig) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.BotConfig{}).Where("user_id = ?", info.BotConfig.UserID).Updates(&info.BotConfig)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *UserAiConfig) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.BotConfig{}).Where("user_id = ?", info.BotConfig.UserID).Delete(&info.BotConfig)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
