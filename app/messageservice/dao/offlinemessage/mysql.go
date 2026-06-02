package offlinemessage

import (
	"aim/app/messageservice/model"
	newerror "aim/pkg/error"
	"context"
	"time"

	"gorm.io/gorm"
)

func setMysql(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:SetMysql")
	result := dbContext.Mysql.Client.WithContext(ctx).Model(&model.OfflineMessageInfo{}).Create(info.OfflineMessageInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil

}
func deleteMysql(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:DeleteMysql")
	messageIDList := make([]int64, 0, len(info.Info))
	result := dbContext.Mysql.Client.WithContext(ctx).Where("message_id IN ?", messageIDList).Delete(&info.OfflineMessageInfo)
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	return nil
}
func getMysql(ctx context.Context, dbContext *model.DBContext, info *MessageInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("mysql:GetMysql")
	messageList := make([]int64, 0, len(info.Info))
	for _, i := range info.Info {
		messageList = append(messageList, i.MessageID)
	}
	GetInfo := make([]*model.OfflineMessageInfo, 0, len(messageList))
	result := dbContext.Mysql.Client.WithContext(ctx).Where("message_id IN ?", messageList).Find(&GetInfo)
	if result.Error == nil && len(GetInfo) == 0 {
		info.Info = nil
		return false, nil
	}
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return false, err2
	}
	info.Info = GetInfo
	return true, nil
}
func ClearMysql(dbContext *model.DBContext, ErrChan chan error) {
	go func() {
		var err error
		defer func(trace string) {
			err = newerror.TranslateError(err).AddErrorTrace(trace)
			ErrChan <- err
		}("mysql:clearOfflineMessageMysql")
		now := time.Now()
		needClearTimeOfKey := now.Unix() - 7*24*60*60
		clearTime := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
		if now.After(clearTime) {
			clearTime = clearTime.AddDate(0, 0, 1)
		}
		select {
		case <-time.After(clearTime.Sub(now)):
		}
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		var result *gorm.DB
		for {
			result = dbContext.Mysql.Client.WithContext(ctx).Where("send_time_second < ?", needClearTimeOfKey).Limit(100).Delete(&model.OfflineMessageInfo{})
			if result.Error != nil {
				err = newerror.MakeError(-1, newerror.CodeDatabaseError, "", result.Error, newerror.LevelError)
				return
			}
			if result.RowsAffected == 0 {
				break
			}
		}
	}()
}
