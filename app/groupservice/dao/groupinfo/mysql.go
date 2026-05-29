package groupinfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"context"

	"gorm.io/gorm"
)

func addWhereInfo(mysqlClient *gorm.DB, info *GroupInfo) *gorm.DB {
	if info.whereWithGroupName {
		mysqlClient.Where("group_name = ?", info.GroupName)
	} else {
		mysqlClient.Where("group_id = ?", info.GroupID)
	}
	return mysqlClient
}
func setMysql(ctx context.Context, dbContext *model.DBContext, info *GroupInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.GroupInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *GroupInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx).Model(&model.GroupInfo{}), info).Find(info.Info)
	if result.Error == nil && len(info.Info) == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *GroupInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.GroupInfo{}).Where("group_id = ?", info.GroupInfo.GroupID).Updates(info.GroupInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *GroupInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("group_id = ?", info.GroupInfo.GroupID).Delete(info.GroupInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
