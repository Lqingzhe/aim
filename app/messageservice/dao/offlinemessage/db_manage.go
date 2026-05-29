package offlinemessage

import (
	"aim/app/messageservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
)

type MessageInfo struct {
	model.OfflineMessageInfo
	Info []*model.OfflineMessageInfo
}
type operate func(*MessageInfo)

func NewStruct(UserAndDeviceID []string, Operates ...operate) *MessageInfo {
	newStruct := &MessageInfo{
		OfflineMessageInfo: model.OfflineMessageInfo{
			UserAndDeviceID: UserAndDeviceID,
		},
	}
	for _, Operate := range Operates {
		Operate(newStruct)
	}
	return newStruct
}
func WithMessageInfo(MessageID int64, SendTimeSecond int64, JsonData []byte) operate {
	return func(info *MessageInfo) {
		info.MessageID = MessageID
		info.JsonData = JsonData
		info.SendTimeSecond = SendTimeSecond
	}
}
func (m *MessageInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setRedis(ctx, DB, m)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Useless Module Unexpectedly Used", fmt.Errorf("%s", "Useless Module Unexpectly Used"), newerror.LevelFatal)
}
func (m *MessageInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getRedis(ctx, DB, m)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	exist, err = getMysql(ctx, DB, m)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	err = deleteRedis(ctx, DB, m)
	if err != nil {
		return false, err
	}
	return true, nil
}
