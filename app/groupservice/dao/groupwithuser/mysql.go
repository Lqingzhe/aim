package groupwithuser

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"context"

	"gorm.io/gorm"
)

func addWhereInfo(mysqlClient *gorm.DB, info *GroupWithUser) *gorm.DB {
	if info.whereWithGroupID {
		mysqlClient = mysqlClient.Where("group_id = ?", info.GroupID)
	}
	if info.whereWithUserID {
		mysqlClient = mysqlClient.Where("user_id = ?", info.UserID)
	}
	return mysqlClient
}

func setMysql(ctx context.Context, dbContext *model.DBContext, info *GroupWithUser) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.GroupWithUserInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *GroupWithUser) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx).Model(&model.GroupWithUserInfo{}), info).Find(info.Info)
	if result.Error == nil && len(info.Info) == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *GroupWithUser) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.GroupWithUserInfo{}).Where("group_id = ?", info.GroupID).Where("user_id = ?", info.UserID).Updates(info.GroupWithUserInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *GroupWithUser) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx), info).Delete(info.GroupWithUserInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
