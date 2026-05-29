package groupapply

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"context"

	"gorm.io/gorm"
)

func addWhereInfo(mysqlClient *gorm.DB, info *GroupApplyInfo) *gorm.DB {
	if info.withWhereGoalID {
		mysqlClient.Where("goal_id = ?", info.GoalID)
	}
	if info.withWhereApplyUserID {
		mysqlClient.Where("apply_user_id = ?", info.ApplyUserID)
	}
	return mysqlClient
}
func setMysql(ctx context.Context, dbContext *model.DBContext, info *GroupApplyInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.GroupApplyInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *GroupApplyInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx), info).Model(&model.GroupApplyInfo{}).Find(info.GroupApplyInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *GroupApplyInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx), info).Delete(info.GroupApplyInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
