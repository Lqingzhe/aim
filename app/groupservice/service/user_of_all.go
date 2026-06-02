package service

import (
	"aim/app/groupservice/dao"
	"aim/app/groupservice/dao/groupmember"
	"aim/app/groupservice/dao/groupwithuser"
	"aim/app/groupservice/dao/sessioninfo"
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"context"
	"fmt"
	"net/http"
	"time"
)

type UserOfAll struct {
	dbContext *model.DBContext
}

func NewUserInfoOfGroup(dbContext *model.DBContext) *UserOfAll {
	return &UserOfAll{
		dbContext: dbContext,
	}
}
func (u *UserOfAll) GetGroupAndSessionID(ctx context.Context, userID int64) (groupID []int64, sessionID []int64, userOfSessionID []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("user_of_all:GetGroupAndSessionInfo")
	groupStruct := groupwithuser.NewStruct(0, userID, "", "", groupwithuser.WithUserID)
	sessionStruct := sessioninfo.NewStruct(0, userID, 0, sessioninfo.WithUserID)
	exist, err := dao.Get(ctx, groupStruct, u.dbContext)
	if err != nil {
		return nil, nil, nil, err
	}
	if !exist {
		groupID = nil
	} else {
		groupID = make([]int64, len(groupStruct.Info))
		for i, v := range groupStruct.Info {
			groupID[i] = v.GroupID
		}
	}
	exist, err = dao.Get(ctx, sessionStruct, u.dbContext)
	if err != nil {
		return nil, nil, nil, err
	}
	if !exist {
		sessionID = nil
	} else {
		sessionID = make([]int64, len(sessionStruct.Info))
		userOfSessionID = make([]int64, len(sessionStruct.Info))
		for i, v := range sessionStruct.Info {
			sessionID[i] = v.SessionID
			userOfSessionID[i] = v.GoalUserID
		}
	}
	return groupID, sessionID, userOfSessionID, nil
} //未获取到的Info空值为nil
func (u *UserOfAll) GetGroupOrSessionExistAndRole(ctx context.Context, groupID int64, userID int64) (role commonmodel.GroupRole, exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("user_of_all:GetGroupOrSessionExistAndRole")
	var finalErr error
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err = dao.Get(ctx, groupMemberStruct, u.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return "", false, err
	}
	if !exist {
		return "", false, nil
	}
	return groupMemberStruct.Info[0].Role, true, finalErr
}
func (u *UserOfAll) SetLastVisitTime(ctx context.Context, userID int64, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("user_of_all:SetLastVisitTime")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithVisitTime([]time.Time{time.Now()}), groupmember.WithWhereMemberID)
	exist, err := dao.Update(ctx, groupMemberStruct, u.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Set Last Visit Time To Unexist User In Group"), newerror.LevelInfo)
	}
	return nil
}
