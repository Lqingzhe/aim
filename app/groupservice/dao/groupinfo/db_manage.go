package groupinfo

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type GroupInfo struct {
	model.GroupInfo
	whereWithGroupName bool
	Info               []*model.GroupInfo
}
type option func(*GroupInfo)

func NewStruct(GroupID int64, GroupName string, Options ...option) *GroupInfo {
	newStruct := &GroupInfo{
		GroupInfo: model.GroupInfo{
			GroupID:   GroupID,
			GroupName: GroupName,
		},
		Info: []*model.GroupInfo{},
	}
	for _, Option := range Options {
		Option(newStruct)
	}
	return newStruct
}
func WithGroupName(groupInfo *GroupInfo) {
	groupInfo.whereWithGroupName = true
}

func (g *GroupInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, g)
	return err
}
func (g *GroupInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
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
func (g *GroupInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
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
func (g *GroupInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, g)
	if err != nil {
		return err
	}
	return nil
}
