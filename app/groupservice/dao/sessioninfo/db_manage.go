package sessioninfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type SessionInfo struct {
	*model.SessionInfo
	Info                []*model.SessionInfo
	whereWithSessionID  bool
	whereWithUserID     bool
	whereWithGoalUserID bool
}
type option func(*SessionInfo)

func NewStruct(SessionID int64, UserID int64, GoalUserID int64, Options ...option) *SessionInfo {
	newStruct := &SessionInfo{
		SessionInfo: &model.SessionInfo{
			SessionID:  SessionID,
			UserID:     UserID,
			GoalUserID: GoalUserID,
		},
		Info: []*model.SessionInfo{},
	}
	for _, Option := range Options {
		Option(newStruct)
	}
	return newStruct
}
func WithSessionID(sessionInfo *SessionInfo) {
	sessionInfo.whereWithSessionID = true
}
func WithUserID(sessionInfo *SessionInfo) {
	sessionInfo.whereWithUserID = true
}
func WithGoalUserID(sessionInfo *SessionInfo) {
	sessionInfo.whereWithGoalUserID = true
}

func (s *SessionInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, s)
	return err
}
func (s *SessionInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, s)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (s *SessionInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, s)
	if err != nil {
		return err
	}
	return nil
}
func (s *SessionInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, s)
	if err != nil {
		return false, err
	}
	return exist, nil
}
