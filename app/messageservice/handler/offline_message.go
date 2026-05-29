package handler

import (
	"aim/app/messageservice/service"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

// SetOfflineMessage implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SetOfflineMessage(ctx context.Context, req *kitexmessageservice.SetOfflineMessageReq) (resp *kitexmessageservice.SetOfflineMessageResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	err = serviceStruct.SetOffLineMessage(ctx, req.GoalUserAndDeviceId, req.JsonData)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("offline_message:SetOfflineMessage")
		newlog.Log(logger, err2.LogLevel, "SetOfflineMessage")
		return nil, err2
	}
	newlog.Log(s.logger, newerror.LevelInfo, "SetOfflineMessage")
	return &kitexmessageservice.SetOfflineMessageResp{}, nil
}

// GetOfflineMessageList implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) GetOfflineMessageList(ctx context.Context, req *kitexmessageservice.GetOfflineMessageListReq) (resp *kitexmessageservice.GetOfflineMessageListResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	jsondataList, exist, err := serviceStruct.GetOffLineMessageList(ctx, req.UserAndDeviceId)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("offline_message:GetOfflineMessageList")
		newlog.Log(logger, err2.LogLevel, "GetOfflineMessageList")
		return nil, err2
	}

	resp = &kitexmessageservice.GetOfflineMessageListResp{
		JsonData: jsondataList,
		Exist:    exist,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetOfflineMessageList")
	return resp, nil
}
