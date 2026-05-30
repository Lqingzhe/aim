package handler

import (
	"aim/app/groupservice/service"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *GroupServiceImpl) GetGroupAndSessionID(ctx context.Context, req *kitexgroupservice.GetGroupAndSessionIDReq) (resp *kitexgroupservice.GetGroupAndSessionIDResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewUserInfoOfGroup(s.DBContext)
	groupIDList, sessionIDList, err := serviceStruct.GetGroupAndSessionID(ctx, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetGroupAndSessionID")
		return nil, err
	}
	resp = &kitexgroupservice.GetGroupAndSessionIDResp{
		GroupIdList:   groupIDList,
		SessionIdList: sessionIDList,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupAndSessionID")
	return resp, nil
}
func (s *GroupServiceImpl) GetGroupOrSessionRoleAndExist(ctx context.Context, req *kitexgroupservice.GetGroupOrSessionRoleAndExistReq) (resp *kitexgroupservice.GetGroupOrSessionRoleAndExistResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	var finalErr error
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewUserInfoOfGroup(s.DBContext)
	role, exist, err := serviceStruct.GetGroupOrSessionExistAndRole(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "UpdateGroupInfoWithUser")
		if newerror.WhetherInterrupt(err, &finalErr) {
			return &kitexgroupservice.GetGroupOrSessionRoleAndExistResp{}, err
		}
		return &kitexgroupservice.GetGroupOrSessionRoleAndExistResp{
			Role:  string(role),
			Exist: exist,
		}, finalErr
	}
	newlog.Log(logger, newerror.LevelInfo, "GetGroupOrSessionRoleAndExist")
	return &kitexgroupservice.GetGroupOrSessionRoleAndExistResp{
		Role:  string(role),
		Exist: exist,
	}, nil
}
