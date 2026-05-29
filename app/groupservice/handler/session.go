package handler

import (
	"aim/app/groupservice/service"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *GroupServiceImpl) CreatSession(ctx context.Context, req *kitexgroupservice.CreatSessionReq) (resp *kitexgroupservice.CreatSessionResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	sessionID, err := serviceStruct.CreatSession(ctx, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "CreateSession")
		return nil, err
	}
	resp = &kitexgroupservice.CreatSessionResp{
		SessionId: sessionID,
	}
	newlog.Log(logger, newerror.LevelInfo, "CreateSession")
	return resp, nil
}

// DeleteSession implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) DeleteSession(ctx context.Context, req *kitexgroupservice.DeleteSessionReq) (resp *kitexgroupservice.DeleteSessionResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	err = serviceStruct.DeleteSession(ctx, req.SessionId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "DeleteSession")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "DeleteSession")
	return &kitexgroupservice.DeleteSessionResp{}, nil
}

// GetFriendLastVisitTime implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetFriendLastVisitTime(ctx context.Context, req *kitexgroupservice.GetFriendLastVisitTimeReq) (resp *kitexgroupservice.GetFriendLastVisitTimeResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	lastVisitTimeString, err := serviceStruct.GetFriendLastVisitTime(ctx, req.SessionId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetFriendLastVisitTime")
		return nil, err
	}
	resp = &kitexgroupservice.GetFriendLastVisitTimeResp{
		LastVisitTime: lastVisitTimeString,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetFriendLastVisitTime")
	return resp, nil
}

// ApplyForFriend implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) ApplyForFriend(ctx context.Context, req *kitexgroupservice.ApplyForFriendReq) (resp *kitexgroupservice.ApplyForFriendResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	err = serviceStruct.ApplyForFriend(ctx, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "ApplyForFriend")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "ApplyForFriend")
	return &kitexgroupservice.ApplyForFriendResp{}, nil
}

// GetFriendApplyList implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetFriendApplyList(ctx context.Context, req *kitexgroupservice.GetFriendApplyListReq) (resp *kitexgroupservice.GetFriendApplyListResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	applyUserIDList, err := serviceStruct.GetFriendApplyList(ctx, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetFriendApplyList")
		return nil, err
	}
	resp = &kitexgroupservice.GetFriendApplyListResp{
		ApplyUserIdList: applyUserIDList,
	}
	return
}

// RefuseFriendApply implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) RefuseFriendApply(ctx context.Context, req *kitexgroupservice.RefuseFriendApplyReq) (resp *kitexgroupservice.RefuseFriendApplyResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewSession(req.CommonInfo.Trace, s.SystemTopic, s.DBContext, s.SnowNode, s.ServiceClient)
	err = serviceStruct.RefuseFriendApply(ctx, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "RefuseFriendApply")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "RefuseFriendApply")
	return &kitexgroupservice.RefuseFriendApplyResp{}, nil
}
