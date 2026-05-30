package groupmuteinfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"context"

	"gorm.io/gorm"
)

func addWhereInfo(mysqlClient *gorm.DB, info *MuteInfo) *gorm.DB {
	if info.whereWithUserID {
		mysqlClient = mysqlClient.Where("user_id = ?", info.UserID)
	}
	return mysqlClient
}
func setMysql(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(&info.GroupMuteInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx).Where("group_id = ?", info.GroupID), info).Model(&model.GroupMuteInfo{}).Find(&info.Info)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("group_id = ?", info.GroupID).Where("user_id = ?", info.UserID).Updates(info.GroupMuteInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *MuteInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("group_id = ?", info.GroupID).Where("user_id = ?", info.UserID).Delete(&info)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
