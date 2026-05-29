package groupapply

import (
	"aim/app/groupservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
)

type GroupApplyInfo struct {
	model.GroupApplyInfo
	Info                 []*model.GroupApplyInfo
	withWhereGoalID      bool
	withWhereApplyUserID bool
}
type operation func(*GroupApplyInfo)

func NewStruct(goalID int64, applyUserID int64, Operations ...operation) *GroupApplyInfo {
	newStruct := &GroupApplyInfo{
		GroupApplyInfo: model.GroupApplyInfo{
			GoalID:      goalID,
			ApplyUserID: applyUserID,
		},
		Info: []*model.GroupApplyInfo{},
	}
	for _, Operation := range Operations {
		Operation(newStruct)
	}
	return newStruct
}
func WithGoalID(info *GroupApplyInfo) {
	info.withWhereGoalID = true
}
func WithApplyUserID(info *GroupApplyInfo) {
	info.withWhereApplyUserID = true
}
func (g *GroupApplyInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
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
func (g *GroupApplyInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
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
func (g *GroupApplyInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Internal Error", fmt.Errorf("This Module Is Not Using"), newerror.LevelFatal)
}
func (g *GroupApplyInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
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
