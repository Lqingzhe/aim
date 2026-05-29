package handler

import (
	"aim/app/userservice/model"
	"aim/app/userservice/service"
	"aim/kitex_gen/kitexuserservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *kitexuserservice.GetUserInfoReq) (resp *kitexuserservice.GetUserInfoResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewUserInfo(s.DBContext, &model.UserInfo{
		UserID: req.UserId,
	}, nil)
	result, err := serviceStruct.GetOtherUserInfo(ctx)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("info:GetUserInfo")
		newlog.Log(logger, err2.LogLevel, "GetUserInfo")
		return nil, err
	}
	resp = &kitexuserservice.GetUserInfoResp{UserInfo: service.TranslateUserInfoModel(result)}
	newlog.Log(logger, newerror.LevelInfo, "GetUserInfo")
	return resp, nil
}

// GetOtherUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetOtherUserInfo(ctx context.Context, req *kitexuserservice.GetOtherUserInfoReq) (resp *kitexuserservice.GetOtherUserInfoResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewUserInfo(s.DBContext, &model.UserInfo{
		UserID: req.GoalUserId,
	}, nil)
	result, err := serviceStruct.GetOtherUserInfo(ctx)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("info:GetOtherUserInfo")
		newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetOtherUserInfo")
		return nil, err
	}
	resp = &kitexuserservice.GetOtherUserInfoResp{UserInfo: service.TranslateUserInfoModel(result)}
	newlog.Log(logger, newerror.LevelInfo, "GetOtherUserInfo")
	return resp, nil
}

// UpdateUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserInfo(ctx context.Context, req *kitexuserservice.UpdateUserInfoReq) (resp *kitexuserservice.UpdateUserInfoResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)

	serviceStruct := service.NewUserInfo(s.DBContext, &model.UserInfo{
		UserID:        req.UserInfo.UserID,
		UserName:      req.UserInfo.UserName,
		Introduction:  req.UserInfo.Introduction,
		BirthdayYear:  req.UserInfo.BirthdayYear,
		BirthdayMonth: req.UserInfo.BirthdayMonth,
		BirthdayDay:   req.UserInfo.BirthdayDay,
	}, nil)
	err = serviceStruct.UpdateUserInfo(ctx, s.UserConfig)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("info:UpdateUserInfo")
		newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "UpdateUserInfo")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "UpdateUserInfo")
	return &kitexuserservice.UpdateUserInfoResp{}, nil
}

// Remark implements the UserServiceImpl interface.
func (s *UserServiceImpl) Remark(ctx context.Context, req *kitexuserservice.RemarkReq) (resp *kitexuserservice.RemarkResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)

	serviceStruct := service.NewUserInfo(s.DBContext, nil, &model.RemarkInfo{
		UserID:     req.RemarkInfo.UserID,
		GoalUserID: req.RemarkInfo.GoalUserID,
		NickName:   req.RemarkInfo.NickName,
	})
	err = serviceStruct.Remark(ctx, s.UserConfig)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("info:Remark")
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "Remark")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "Remark")
	return &kitexuserservice.RemarkResp{}, nil
}
