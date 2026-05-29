package nickname

import (
	"aim/app/userservice/model"
	newerror "aim/pkg/error"
	"context"
)

func setMysql(ctx context.Context, dbContext *model.DBContext, info *NickNameInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.RemarkInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *NickNameInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.RemarkInfo{}).Where("user_id = ?", info.RemarkInfo.UserID).Find(info.Info)
	if result.Error == nil && len(info.Info) == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *NickNameInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.RemarkInfo{}).Where("user_id = ?", info.RemarkInfo.UserID).Where("goal_id = ?", info.RemarkInfo.GoalUserID).Updates(info.RemarkInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *NickNameInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("user_id = ?", info.RemarkInfo.UserID).Where("goal_id = ?", info.RemarkInfo.GoalUserID).Delete(info.RemarkInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
