package handler

import (
	"aim/app/groupservice/service"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *GroupServiceImpl) SetMute(ctx context.Context, req *kitexgroupservice.SetMuteReq) (resp *kitexgroupservice.SetMuteResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewMute(req.CommonInfo.Trace, s.GroupNoticeTopic, s.DBContext, s.GroupConfig)
	err = serviceStruct.SetMute(ctx, req.UserId, req.GroupId, req.GoalUserId, req.MuteTimeSecond, req.MuteReason)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SetMute")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "SetMute")
	return &kitexgroupservice.SetMuteResp{}, nil
}

// ReleaseMute implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) ReleaseMute(ctx context.Context, req *kitexgroupservice.ReleaseMuteReq) (resp *kitexgroupservice.ReleaseMuteResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewMute(req.CommonInfo.Trace, s.GroupNoticeTopic, s.DBContext, s.GroupConfig)
	err = serviceStruct.ReleaseMute(ctx, req.UserId, req.GroupId, req.GoalUserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "ReleaseMute")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "ReleaseMute")
	return &kitexgroupservice.ReleaseMuteResp{}, nil
}

// GetMuteStatus implements the GroupServiceImpl interface.
func (s *GroupServiceImpl) GetMuteStatus(ctx context.Context, req *kitexgroupservice.GetMuteStatusReq) (resp *kitexgroupservice.GetMuteStatusResp, err error) {
	var finalErr error
	logger := newlog.AddTraceAndEquipID(s.Logger, req.CommonInfo.Trace, s.EquipID)
	serviceStruct := service.NewMute(req.CommonInfo.Trace, s.GroupNoticeTopic, s.DBContext, s.GroupConfig)
	muteReason, muteEndTime, isMute, err := serviceStruct.GetMuteStatus(ctx, req.UserId, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetMuteStatus")
		if newerror.WhetherInterrupt(err, &finalErr) {
			return nil, err
		}
	}
	resp = &kitexgroupservice.GetMuteStatusResp{
		MuteReason:  muteReason,
		MuteEndTime: muteEndTime,
		IsMute:      isMute,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetMuteStatus")
	return resp, finalErr
}
