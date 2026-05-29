package sessioninfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"context"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func addWhereInfo(mysqlClient *gorm.DB, info *SessionInfo) *gorm.DB {
	if info.whereWithSessionID {
		mysqlClient = mysqlClient.Where("session_id = ?", info.SessionID)
	}
	if info.whereWithUserID {
		mysqlClient = mysqlClient.Where("user_id = ?", info.UserID)
	}
	if info.whereWithGoalUserID {
		mysqlClient = mysqlClient.Where("goal_user_id = ?", info.GoalUserID)
	}
	return mysqlClient
}
func setMysql(ctx context.Context, dbContext *model.DBContext, info *SessionInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:setMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(info.SessionInfo)

	if result.Error != nil {
		if isContext, err2 := newerror.IsContextError(err); isContext {
			return err2
		}
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return newerror.MakeError(http.StatusConflict, newerror.CodeResourceDuplicate, "SessionID Already Exist", fmt.Errorf(`SessionID Already Set, Should Use "Update"`), newerror.LevelFatal)
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeDatabaseError, "Database Error", result.Error, newerror.LevelError)
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *SessionInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx).Model(&model.SessionInfo{}), info).Find(info.Info)
	if result.Error == nil && len(info.Info) == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func updateMysql(ctx context.Context, dbContext *model.DBContext, info *SessionInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:UpdateMysql")
	result := addWhereInfo(dbContext.Mysql.Client.WithContext(ctx).Model(&model.SessionInfo{}), info).Updates(info.SessionInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *SessionInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("session_id = ?", info.SessionID).Delete(info.SessionInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
