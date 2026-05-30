package groupmuteinfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"time"
)

type MuteInfo struct {
	model.GroupMuteInfo
	Info            []*model.GroupMuteInfo
	whereWithUserID bool
}
type operation func(*MuteInfo)

func NewStruct(groupID int64, userID int64, muteEndTime time.Time, muteReason string, Operations ...operation) *MuteInfo {
	newStruct := &MuteInfo{
		GroupMuteInfo: model.GroupMuteInfo{
			GroupID:     groupID,
			UserID:      userID,
			MuteEndTime: muteEndTime,
			MuteReason:  muteReason,
		},
		Info: []*model.GroupMuteInfo{},
	}
	if len(Operations) > 0 {
		for _, Operate := range Operations {
			Operate(newStruct)
		}
	}
	return newStruct
}
func WithWhereUserID(muteInfo *MuteInfo) {
	muteInfo.whereWithUserID = true
}
func (m *MuteInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, m)
	if err != nil {
		return err
	}
	return nil
}
func (m *MuteInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	var finalErr error
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	if m.whereWithUserID {
		exist, err = getRedis(ctx, DB, m)
		if newerror.WhetherInterrupt(err, &finalErr) {
			return false, err
		}
		if exist {
			return true, nil
		} else {
			exist, err = getMysql(ctx, DB, m)
			if newerror.WhetherInterrupt(err, &finalErr) {
				return false, err
			}
			newStruct := &MuteInfo{
				GroupMuteInfo: *m.Info[0],
				Info:          []*model.GroupMuteInfo{},
			} //存在则Info【0】有值，不存在则默认为空值
			newStruct.GroupMuteInfo.UserID = m.UserID
			newStruct.GroupMuteInfo.GroupID = m.GroupID
			err = setRedis(ctx, DB, newStruct)
			if newerror.WhetherInterrupt(err, &finalErr) {
				return false, err
			}
			return exist, finalErr
		}

	} else {
		exist, err = getMysql(ctx, DB, m)
		if err != nil {
			return false, err
		}
		return exist, nil
	}

}
func (m *MuteInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	var finalErr error
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	err = deleteRedis(ctx, DB, m)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, m)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return false, err
	}
	return exist, nil
}
func (m *MuteInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	var finalErr error
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteRedis(ctx, DB, m)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	err = deleteMysql(ctx, DB, m)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	return nil
}
