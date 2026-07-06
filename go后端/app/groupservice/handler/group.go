package handler

import (
	"aim/app/groupservice/service"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *GroupServiceImpl) GetGroupInfo(ctx context.Context, req *kitexgroupservice.GetGroupInfoReq) (resp *kitexgroupservice.GetGroupInfoResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	groupInfo, err := serviceStruct.GetGroupInfo(ctx, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetGroupInfo")
		return nil, err
	}
	resp = &kitexgroupservice.GetGroupInfoResp{
		GroupId:   groupInfo.GroupID,
		GroupName: groupInfo.GroupName,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupInfo")
	return resp, nil
}

// ChangeGroupInfo implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) ChangeGroupInfo(ctx context.Context, req *kitexgroupservice.ChangeGroupInfoReq) (resp *kitexgroupservice.ChangeGroupInfoResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.ChangeGroupInfo(ctx, req.GroupId, req.UserId, req.GroupName)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "ChangeGroupInfo")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "ChangeGroupInfo")
	return &kitexgroupservice.ChangeGroupInfoResp{}, nil
}

// SearchGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) SearchGroup(ctx context.Context, req *kitexgroupservice.SearchGroupReq) (resp *kitexgroupservice.SearchGroupResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	groupID, err := serviceStruct.SearchGroup(ctx, req.GroupName)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SearchGroup")
		return nil, err
	}
	resp = &kitexgroupservice.SearchGroupResp{
		GroupIdList: groupID,
	}
	newlog.Log(logger, newerror.LevelInfo, "SearchGroup")
	return resp, nil
}

// CreateGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) CreateGroup(ctx context.Context, req *kitexgroupservice.CreateGroupReq) (resp *kitexgroupservice.CreateGroupResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	groupID, err := serviceStruct.CreateGroup(ctx, req.UserId, req.GroupName)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "CreateGroup")
		return nil, err
	}
	resp = &kitexgroupservice.CreateGroupResp{
		GroupId: groupID,
	}
	newlog.Log(logger, newerror.LevelInfo, "CreateGroup")
	return resp, nil
}

// DeleteGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) DeleteGroup(ctx context.Context, req *kitexgroupservice.DeleteGroupReq) (resp *kitexgroupservice.DeleteGroupResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.DeleteGroup(ctx, req.UserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "DeleteGroup")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "DeleteGroup")
	return &kitexgroupservice.DeleteGroupResp{}, nil
}

// LeaveGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) LeaveGroup(ctx context.Context, req *kitexgroupservice.LeaveGroupReq) (resp *kitexgroupservice.LeaveGroupResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.LeaveGroup(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "LeaveGroup")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "LeaveGroup")
	return &kitexgroupservice.LeaveGroupResp{}, nil
}

// SetGroupApply implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) SetGroupApply(ctx context.Context, req *kitexgroupservice.SetGroupApplyReq) (resp *kitexgroupservice.SetGroupApplyResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.SetGroupApply(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SetGroupApply")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "SetGroupApply")
	return &kitexgroupservice.SetGroupApplyResp{}, nil
}

// GetGroupApplyList implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetGroupApplyList(ctx context.Context, req *kitexgroupservice.GetGroupApplyListReq) (resp *kitexgroupservice.GetGroupApplyListResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	applyUserIDList, err := serviceStruct.GetGroupApplyList(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetGroupApplyList")
		return nil, err
	}
	resp = &kitexgroupservice.GetGroupApplyListResp{
		ApplyUserIdList: applyUserIDList,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupApplyList")
	return resp, nil
}

// GetLastVisitTime implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetLastVisitTime(ctx context.Context, req *kitexgroupservice.GetLastVisitTimeReq) (resp *kitexgroupservice.GetLastVisitTimeResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	userIDList, lastVisitTimeList, err := serviceStruct.GetLastVisitTime(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetLastVisitTime")
		return nil, err
	}
	resp = &kitexgroupservice.GetLastVisitTimeResp{
		UserIdList:        userIDList,
		LastVisitTimeList: lastVisitTimeList,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetLastVisitTime")
	return resp, nil
}
func (s *GroupServiceImpl) SetLastVisitTime(ctx context.Context, req *kitexgroupservice.SetLastVisitTimeReq) (resp *kitexgroupservice.SetLastVisitTimeResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewUserInfoOfGroup(s.DBContext)
	err = serviceStruct.SetLastVisitTime(ctx, req.UserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SetLastVisitTime")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "SetLastVisitTime")
	return &kitexgroupservice.SetLastVisitTimeResp{}, nil
}

// AgreeGroupApply implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) AgreeGroupApply(ctx context.Context, req *kitexgroupservice.AgreeGroupApplyReq) (resp *kitexgroupservice.AgreeGroupApplyResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.AgreeGroupApply(ctx, req.GroupId, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "AgreeGroupApply")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "AgreeGroupApply")
	return &kitexgroupservice.AgreeGroupApplyResp{}, nil
}
func (s *GroupServiceImpl) RefuseGroupApply(ctx context.Context, req *kitexgroupservice.RefuseGroupApplyReq) (resp *kitexgroupservice.RefuseGroupApplyResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.RefuseGroupApply(ctx, req.GroupId, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "RefuseGroupApply")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "RefuseGroupApply")
	return &kitexgroupservice.RefuseGroupApplyResp{}, nil
}

// TransformGroupOwner implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) TransformGroupOwner(ctx context.Context, req *kitexgroupservice.TransformGroupOwnerReq) (resp *kitexgroupservice.TransformGroupOwnerResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.TransformGroupOwner(ctx, req.GroupId, req.UserId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "TransformGroupOwner")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "TransformGroupOwner")
	return &kitexgroupservice.TransformGroupOwnerResp{}, nil
}

// KickOutGroup implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) KickOutGroup(ctx context.Context, req *kitexgroupservice.KickOutGroupReq) (resp *kitexgroupservice.KickOutGroupResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.KickOutGroup(ctx, req.UserId, req.GoalUserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "KickOutGroup")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "KickOutGroup")
	return &kitexgroupservice.KickOutGroupResp{}, nil
}

// SetManager implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) SetManager(ctx context.Context, req *kitexgroupservice.SetManagerReq) (resp *kitexgroupservice.SetManagerResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.SetManager(ctx, req.UserId, req.GoalUserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SetManager")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "SetManager")
	return &kitexgroupservice.SetManagerResp{}, nil
}

// RevokeManager implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) RevokeManager(ctx context.Context, req *kitexgroupservice.RevokeManagerReq) (resp *kitexgroupservice.RevokeManagerResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.RevokeManager(ctx, req.UserId, req.GoalUserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "RevokeManager")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "RevokeManager")
	return &kitexgroupservice.RevokeManagerResp{}, nil
}

// GetGroupInfoWithUser implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetGroupInfoWithUser(ctx context.Context, req *kitexgroupservice.GetGroupInfoWithUserReq) (resp *kitexgroupservice.GetGroupInfoWithUserResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	groupWithUserInfo, err := serviceStruct.GetGroupInfoWithUser(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetGroupInfoWithUser")
		return nil, err
	}
	resp = &kitexgroupservice.GetGroupInfoWithUserResp{
		GroupId:         groupWithUserInfo.GroupID,
		GroupRemarkName: groupWithUserInfo.GroupRemarkName,
		Role:            string(groupWithUserInfo.Role),
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupInfoWithUser")
	return resp, nil
}

// UpdateGroupInfoWithUser implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) UpdateGroupInfoWithUser(ctx context.Context, req *kitexgroupservice.UpdateGroupInfoWithUserReq) (resp *kitexgroupservice.UpdateGroupInfoWithUserResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	err = serviceStruct.UpdateGroupInfoWithUser(ctx, req.UserId, req.GroupId, req.GroupRemarkName)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "UpdateGroupInfoWithUser")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "UpdateGroupInfoWithUser")
	return &kitexgroupservice.UpdateGroupInfoWithUserResp{}, nil
}
func (s *GroupServiceImpl) GetGroupUserID(ctx context.Context, req *kitexgroupservice.GetGroupUserIDReq) (resp *kitexgroupservice.GetGroupUserIDResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewGroup(req.CommonInfo.Trace, s.GroupNoticeTopic, s.SystemTopic, s.DBContext, s.GroupConfig, s.SnowNode, s.ServiceClient)
	userListID, err := serviceStruct.GetGroupUserID(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetGroupUserID")
		return nil, err
	}
	resp = &kitexgroupservice.GetGroupUserIDResp{
		UserIdList: userListID,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupUserID")
	return resp, nil
}
