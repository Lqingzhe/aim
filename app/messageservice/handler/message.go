package handler

import (
	"aim/app/messageservice/service"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

// SendMessage implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SendMessage(ctx context.Context, req *kitexmessageservice.SendMessageReq) (resp *kitexmessageservice.SendMessageResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageID, err := serviceStruct.SendMessage(ctx, req.GroupId, req.UserId, req.MessageContent, req.IsAi)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:SendMessage")
		newlog.Log(logger, err2.LogLevel, "SendMessage")
		return nil, err2
	}
	resp = &kitexmessageservice.SendMessageResp{MessageId: messageID}
	newlog.Log(logger, newerror.LevelInfo, "SendMessage")
	return resp, nil
}

// SendFile implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SendFile(ctx context.Context, req *kitexmessageservice.SendFileReq) (resp *kitexmessageservice.SendFileResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageID, err := serviceStruct.SendFile(ctx, req.GroupId, req.UserId, req.FileName, req.ContentType, req.DataStream)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:SendFile")
		newlog.Log(logger, err2.LogLevel, "SendFile")
		return nil, err2
	}
	resp = &kitexmessageservice.SendFileResp{MessageId: messageID}
	newlog.Log(logger, newerror.LevelInfo, "SendFile")
	return resp, nil
}

// SendVoice implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SendVoice(ctx context.Context, req *kitexmessageservice.SendVoiceReq) (resp *kitexmessageservice.SendVoiceResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageID, err := serviceStruct.SendVoice(ctx, req.GroupId, req.UserId, req.ContentType, req.VoiceTimeSecond, req.DataStream)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:SendVoice")
		newlog.Log(logger, err2.LogLevel, "SendVoice")
		return nil, err2
	}
	resp = &kitexmessageservice.SendVoiceResp{MessageId: messageID}
	newlog.Log(logger, newerror.LevelInfo, "SendVoice")
	return resp, nil
}

// SendPicture implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SendPicture(ctx context.Context, req *kitexmessageservice.SendPictureReq) (resp *kitexmessageservice.SendPictureResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageID, err := serviceStruct.SendPicture(ctx, req.GroupId, req.UserId, req.ContentType, req.DataStream)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:SendPicture")
		newlog.Log(logger, err2.LogLevel, "SendPicture")
		return nil, err2
	}
	resp = &kitexmessageservice.SendPictureResp{MessageId: messageID}
	newlog.Log(logger, newerror.LevelInfo, "SendPicture")
	return resp, nil
}

// WithdrawMessage implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) WithdrawMessage(ctx context.Context, req *kitexmessageservice.WithdrawMessageReq) (resp *kitexmessageservice.WithdrawMessageResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	err = serviceStruct.WithdrawMessage(ctx, req.GroupId, req.UserId, req.MessageId)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:WithdrawMessage")
		newlog.Log(logger, err2.LogLevel, "WithdrawMessage")
		return nil, err2
	}
	newlog.Log(logger, newerror.LevelInfo, "WithdrawMessage")
	return &kitexmessageservice.WithdrawMessageResp{}, nil
}
func (s *KitexMessageServiceImpl) DeleteMessageAllGroup(ctx context.Context, req *kitexmessageservice.DeleteMessageAllGroupReq) (resp *kitexmessageservice.DeleteMessageAllGroupResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	err = serviceStruct.DeleteMessageAllGroup(ctx, req.GroupId)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:DeleteMessageAllGroup")
		newlog.Log(logger, err2.LogLevel, "DeleteMessageAllGroup")
		return nil, err2
	}
	newlog.Log(logger, newerror.LevelInfo, "DeleteMessageAllGroup")
	return &kitexmessageservice.DeleteMessageAllGroupResp{}, nil
}

// GetMessageList implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) GetMessageList(ctx context.Context, req *kitexmessageservice.GetMessageListReq) (resp *kitexmessageservice.GetMessageListResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageList, err := serviceStruct.GetMessageList(ctx, req.GroupId, req.UserId, req.StartTimeSecond, req.EndTimeSecond)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:GetMessageList")
		newlog.Log(logger, err2.LogLevel, "GetMessageList")
		return nil, err2
	}
	MessageInfo := make([]*kitexcommonmodel.KitexMessageInfo, 0, len(messageList))
	for _, v := range messageList {
		MessageInfo = append(MessageInfo, &kitexcommonmodel.KitexMessageInfo{
			GroupId:             v.GroupID,
			UserId:              v.UserID,
			MessageId:           v.MessageID,
			MessageContent:      v.MessageContent,
			ContentType:         v.ContentType,
			VoiceDurationSecond: v.VoiceDurationSecond,
			IsAi:                v.IsAI,
			MessageType:         v.MessageType,
			SendTimeSecond:      v.SendTimeSecond,
		})
	}
	resp = &kitexmessageservice.GetMessageListResp{
		MessageInfo: MessageInfo,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetMessageList")
	return resp, nil
}

// GetMessageIDList implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) GetNewMessage(ctx context.Context, req *kitexmessageservice.GetNewMessageReq) (resp *kitexmessageservice.GetNewMessageResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	messageID, sendTimeSecond, messageType, messageContent, err := serviceStruct.GetNewMessage(ctx, req.GroupId, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:GetNewMessage")
		newlog.Log(logger, err2.LogLevel, "GetNewMessage")
		return nil, err2
	}
	resp = &kitexmessageservice.GetNewMessageResp{
		MessageId:      messageID,
		SendTimeSecond: sendTimeSecond,
		MessageType:    messageType,
		MessageContent: messageContent,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetMessageIDList")
	return resp, nil
}

func (s *KitexMessageServiceImpl) GetFileContent(ctx context.Context, req *kitexmessageservice.GetFileContentReq) (resp *kitexmessageservice.GetFileContentResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	dataStream, contentType, err := serviceStruct.GetFileContent(ctx, req.GroupId, req.UserId, req.MessageId)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:GetFileContent")
		newlog.Log(logger, err2.LogLevel, "GetFileContent")
		return nil, err2
	}
	resp = &kitexmessageservice.GetFileContentResp{
		DataStream:  dataStream,
		ContentType: contentType,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetFileContent")
	return resp, nil
}

// SendGroupNotice implements the KitexMessageServiceImpl interface.
func (s *KitexMessageServiceImpl) SendGroupNotice(ctx context.Context, req *kitexmessageservice.SendGroupNoticeReq) (resp *kitexmessageservice.SendGroupNoticeResp, err error) {
	logger := newlog.AddTraceAndEquipID(s.logger, req.CommonInfo.Trace, s.equipID)
	serviceStruct := service.NewMessageService(req.CommonInfo.Trace, s.messageConfig, s.snowFlake, s.dbContext, s.serviceClient, s.messageTopic, s.groupNoticeTopic, s.systemTopic)
	err = serviceStruct.SendGroupNotice(ctx, req.GroupId, req.UserId, req.MessageContent)
	if err != nil {
		err2 := newerror.TranslateError(err).AddErrorTrace("message:SendGroupNotice")
		newlog.Log(logger, err2.LogLevel, "SendGroupNotice")
		return nil, err2
	}
	newlog.Log(logger, newerror.LevelInfo, "SendGroupNotice")
	return &kitexmessageservice.SendGroupNoticeResp{}, nil
}
