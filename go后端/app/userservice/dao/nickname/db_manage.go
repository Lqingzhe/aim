package nickname

import (
	"aim/app/userservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type NickNameInfo struct {
	model.RemarkInfo
	Info []*model.RemarkInfo
}

func NewStruct(UserID int64, GoalID int64, NickName string) *NickNameInfo {
	return &NickNameInfo{
		RemarkInfo: model.RemarkInfo{
			UserID:     UserID,
			GoalUserID: GoalID,
			NickName:   NickName,
		},
		Info: []*model.RemarkInfo{},
	}
}
func (n *NickNameInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, n)
	if err != nil {
		return err
	}
	return nil
}
func (n *NickNameInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, n)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (n *NickNameInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, n)
	if err != nil {
		return err
	}
	return nil
}
func (n *NickNameInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, n)
	if err != nil {
		return false, err
	}
	return exist, nil
}
