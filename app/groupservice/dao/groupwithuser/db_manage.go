package groupwithuser

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type GroupWithUser struct {
	model.GroupWithUserInfo
	Info             []*model.GroupWithUserInfo
	whereWithGroupID bool
	whereWithUserID  bool
}
type operation func(*GroupWithUser)

func NewStruct(GroupID int64, UserID int64, GroupRemarkName string, role commonmodel.GroupRole, Operations ...operation) *GroupWithUser {
	newStruct := &GroupWithUser{
		GroupWithUserInfo: model.GroupWithUserInfo{
			GroupID:         GroupID,
			UserID:          UserID,
			GroupRemarkName: GroupRemarkName,
			Role:            role,
		},
		Info: []*model.GroupWithUserInfo{},
	}

	if len(Operations) > 0 {
		for _, Operate := range Operations {
			Operate(newStruct)
		}
	}
	return newStruct
}
func WithGroupID(groupWithUser *GroupWithUser) {
	groupWithUser.whereWithGroupID = true
}
func WithUserID(groupWithUser *GroupWithUser) {
	groupWithUser.whereWithUserID = true
}

func (g *GroupWithUser) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, g)
	if err != nil {
		return err
	}
	return nil
}
func (g *GroupWithUser) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, g)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (g *GroupWithUser) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, g)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (g *GroupWithUser) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, g)
	return err
}
