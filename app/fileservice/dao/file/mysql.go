package file

import (
	"aim/app/fileservice/model"
	newerror "aim/pkg/error"
	"context"
)

func setMysql(ctx context.Context, dbContext *model.DBContext, info *FileInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Create(&info.FileModel)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *FileInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	getInfo := &model.FileModel{}
	result := dbContext.Mysql.Client.WithContext(ctx).Where("file_id = ?", info.FileID).First(getInfo)
	if result.Error == nil && result.RowsAffected == 0 {
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	info.Info = append(info.Info, getInfo)
	return true, nil
}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *FileInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Where("file_id = ?", info.FileID).Delete(&info.FileModel)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
