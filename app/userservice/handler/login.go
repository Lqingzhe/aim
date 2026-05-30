package handler

import (
	"aim/app/userservice/model"
	"aim/app/userservice/service"
	"aim/kitex_gen/kitexuserservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *UserServiceImpl) Register(ctx context.Context, req *kitexuserservice.RegisterReq) (resp *kitexuserservice.RegisterResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewLoginInfo(s.DBContext, &model.UserLoginInfo{Password: req.Password}, nil)
	userID, err := serviceStruct.Register(ctx, s.UserConfig, s.SnowNode)
	if err != nil {
		err2 := newerror.TranslateError(err)
		newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "Register")
		return nil, err2
	}
	resp = &kitexuserservice.RegisterResp{UserId: userID}
	newlog.Log(logger, newerror.LevelInfo, "Register")
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *kitexuserservice.LoginReq) (resp *kitexuserservice.LoginResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.Logger, req.CommonInfo.Trace)
	serviceStruct := service.NewLoginInfo(s.DBContext, &model.UserLoginInfo{
		UserID:   req.UserId,
		Password: req.Password,
	}, nil)
	err = serviceStruct.Login(ctx)
	if err != nil {
		err2 := newerror.TranslateError(err)
		newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "Login")
		return nil, err2
	}
	newlog.Log(logger, newerror.LevelInfo, "Login")
	return &kitexuserservice.LoginResp{}, nil
}
