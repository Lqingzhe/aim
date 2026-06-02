package service

import (
	"aim/app/messageservice/dao"
	"aim/app/messageservice/dao/message"
	"aim/app/messageservice/dao/offlinemessage"
	"aim/app/messageservice/model"
	"aim/commonmodel"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexfileservice"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
)

type MessageService struct {
	traceID          string
	messageConfig    commonmodel.MessageConfig
	snowFlake        *snowflake.Node
	dbContext        *model.DBContext
	serviceClient    model.ServiceClient
	messageTopic     sarama.SyncProducer
	groupNoticeTopic sarama.SyncProducer
	systemTopic      sarama.SyncProducer
}

func NewMessageService(traceID string, messageConfig commonmodel.MessageConfig, snowFlake *snowflake.Node, dbContext *model.DBContext, serviceClient model.ServiceClient, messageTopic sarama.SyncProducer, groupNoticeTopic sarama.SyncProducer, systemTopic sarama.SyncProducer) *MessageService {
	return &MessageService{
		traceID:          traceID,
		messageConfig:    messageConfig,
		snowFlake:        snowFlake,
		dbContext:        dbContext,
		serviceClient:    serviceClient,
		messageTopic:     messageTopic,
		groupNoticeTopic: groupNoticeTopic,
		systemTopic:      systemTopic,
	}
}
func judgeExistInGroupAndMute(ctx context.Context, m *MessageService, groupID int64, userID int64) (err error) {
	var finalErr error
	GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupOrSessionRoleAndExistResp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	if !GetGroupOrSessionRoleAndExistResp.Exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Send Message To Unjoin Group"), newerror.LevelInfo)
	}
	GetMuteStatusReq := kitexgroupservice.GetMuteStatusReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetMuteStatusResp, err := m.serviceClient.GroupService.GetMuteStatus(ctx, &GetMuteStatusReq)

	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	if GetMuteStatusResp.IsMute {
		muteReason := GetMuteStatusResp.MuteReason
		muteEndTime := GetMuteStatusResp.MuteEndTime
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, fmt.Sprintf("You Are In MutIng Because : %s ,Before : %s", muteReason, muteEndTime), fmt.Errorf("Try To Send Message In Muting"), newerror.LevelInfo)
	}
	return finalErr
}
func sendNewMessageNoticeToGroupMember(ctx context.Context, m *MessageService, groupID int64, userID int64) (err error) {
	var finalErr error
	GetGroupUserIDReq := &kitexgroupservice.GetGroupUserIDReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupUserIDResp, err := m.serviceClient.GroupService.GetGroupUserID(ctx, GetGroupUserIDReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	//广播
	groupMessageStruct := commonmodel.KafkaNewMessageNotice{
		TraceID:        m.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     GetGroupUserIDResp.UserIdList,
		SessionID:      groupID,
		MessageCode:    commonmodel.MessageCode_GroupMessage,
	}
	_, _, err = tool.SendKafkaNewMessageNotice(m.messageTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	return finalErr
}
func (m *MessageService) SendMessage(ctx context.Context, groupID int64, userID int64, messageContent string, IsAi bool) (messageID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:SendMessage")
	var finalErr error
	if int64(len(messageContent)) > m.messageConfig.MaxMessageByteLength {
		return 0, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Message Too Long", fmt.Errorf("Send Too Long Message"), newerror.LevelInfo)
	}
	if err = judgeExistInGroupAndMute(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, finalErr
	}
	if userIDList := tool.GetMessageEmphasizeUserID(messageContent); len(userIDList) != 0 {
		goalUserList := make([]int64, 0, len(userIDList))
		for _, id := range userIDList {
			GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
				CommonInfo: &kitexcommonmodel.CommonInfo{Trace: m.traceID},
				GroupId:    groupID,
				UserId:     id,
			}
			Resp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)

			if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
				return 0, finalErr
			}
			if Resp.Exist {
				goalUserList = append(goalUserList, id)
			}
		}
		systemMessageStruct := commonmodel.KafkaSystemMessage{
			TraceID:        m.traceID,
			SendTimeSecond: time.Now().Unix(),
			GoalUserID:     goalUserList,
			Data: map[string]interface{}{
				"group_id": strconv.FormatInt(groupID, 10),
				"user_id":  strconv.FormatInt(userID, 10),
			},
			MessageCode: commonmodel.MessageCode_GroupEmpphasizeMessage,
		}
		_, _, err = tool.SendKafkaSystemMessage(m.systemTopic, systemMessageStruct)
		if err != nil {
			return 0, newerror.MakeError(http.StatusInternalServerError, newerror.CodeMessageQueueError, "Internal Error, Use '@' Later", err, newerror.LevelError)
		}
	} else if aiMessage, isNeedAi := tool.GetMessageAiChatMessage(messageContent); isNeedAi {
		//向aiService调用rcp，由aiService决定是降级熔断还是向kafka发送消息
		//返回的error存在warn级别的，需IsInterrupt
		///err=
		_ = fmt.Sprintf("%s", aiMessage)

		if newerror.WhetherInterrupt(err, &finalErr) {
			return 0, err
		}
	}
	messageID = m.snowFlake.Generate().Int64()
	serviceStruct := message.NewStruct(groupID, messageID, userID, time.Now(), message.WithMessageInfo(messageContent), message.WithAI(IsAi))
	err = dao.Add(ctx, serviceStruct, m.dbContext)
	if err != nil {
		return 0, err
	}
	if err = sendNewMessageNoticeToGroupMember(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, err
	}
	return messageID, finalErr
}

func (m *MessageService) SendFile(ctx context.Context, groupID int64, userID int64, fileName string, contentType string, dataStream []byte) (messageID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:SendFile")
	var finalErr error
	messageID = m.snowFlake.Generate().Int64()
	if int64(len(dataStream)) > m.messageConfig.MaxFileByteLength {
		return 0, newerror.MakeError(http.StatusRequestEntityTooLarge, newerror.CodeRequestBodyTooBig, "The File Is Too Large", fmt.Errorf("User Send Too Large Message"), newerror.LevelInfo)
	}
	if err = judgeExistInGroupAndMute(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, finalErr
	}
	CreateFileReq := kitexfileservice.CreateFileReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		FileName:    fileName,
		ContentType: contentType,
		DataStream:  dataStream,
		FileType:    string(commonmodel.FileType_File),
	}
	CreateFileResp, err := m.serviceClient.FileService.CreateFile(ctx, &CreateFileReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return 0, finalErr
	}
	fileStorageID := CreateFileResp.FileId
	serviceStruct := message.NewStruct(groupID, messageID, userID, time.Now(), message.WithFileInfo(fileStorageID, contentType))
	err = dao.Add(ctx, serviceStruct, m.dbContext)
	if err != nil {
		return 0, err
	}
	if err = sendNewMessageNoticeToGroupMember(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, err
	}
	return messageID, finalErr
}
func (m *MessageService) SendVoice(ctx context.Context, groupID int64, userID int64, contentType string, voiceTimeSecond int64, dataStream []byte) (messageID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:SendVoice")
	var finalErr error
	messageID = m.snowFlake.Generate().Int64()
	if voiceTimeSecond > m.messageConfig.MaxVoiceTimeSecond {
		return 0, newerror.MakeError(http.StatusRequestEntityTooLarge, newerror.CodeRequestBodyTooBig, "The Voice Is Too Large", fmt.Errorf("User Send Too Long Voice"), newerror.LevelInfo)
	}
	CreateFileReq := kitexfileservice.CreateFileReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		DataStream:              dataStream,
		FileType:                string(commonmodel.FileType_Voice),
		ContentType:             contentType,
		VoiceDurationTimeSecond: voiceTimeSecond,
	}
	CreateFileResp, err := m.serviceClient.FileService.CreateFile(ctx, &CreateFileReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return 0, finalErr
	}
	fileStorageID := CreateFileResp.FileId
	messageStruct := message.NewStruct(groupID, messageID, userID, time.Now(), message.WithVoiceInfo(voiceTimeSecond, fileStorageID, contentType))
	err = dao.Add(ctx, messageStruct, m.dbContext)
	if err != nil {
		return 0, err
	}
	if err = sendNewMessageNoticeToGroupMember(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, err
	}
	return messageID, finalErr
}
func (m *MessageService) SendPicture(ctx context.Context, groupID int64, userID int64, contentType string, dataStream []byte) (messageID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:SendPicture")
	var finalErr error
	messageID = m.snowFlake.Generate().Int64()
	if err = judgeExistInGroupAndMute(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, err
	}
	CreateFileReq := kitexfileservice.CreateFileReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		ContentType: contentType,
		DataStream:  dataStream,
	}
	CreateFileResp, err := m.serviceClient.FileService.CreateFile(ctx, &CreateFileReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return 0, finalErr
	}
	fileStorageID := CreateFileResp.FileId
	messageStruct := message.NewStruct(groupID, messageID, userID, time.Now(), message.WithPictureInfo(fileStorageID, contentType))
	err = dao.Add(ctx, messageStruct, m.dbContext)
	if err != nil {
		return 0, err
	}
	if err = sendNewMessageNoticeToGroupMember(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return 0, err
	}
	return messageID, finalErr
}
func (m *MessageService) WithdrawMessage(ctx context.Context, groupID int64, userID int64, messageID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:WithdrawMessage")
	var finalErr error
	if err = judgeExistInGroupAndMute(ctx, m, groupID, userID); newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	messageStruct := message.NewStruct(groupID, messageID, 0, time.Now(), message.GetWithMessageID)
	exist, err := dao.Get(ctx, messageStruct, m.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "This Message Is Not Exist", fmt.Errorf("Try To Withdraw Unexist Message"), newerror.LevelInfo)
	}
	if messageStruct.Info[0].UserID != userID {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Only Can Withdraw The Message You Send", fmt.Errorf("Try To Withdraw Other's Message"), newerror.LevelInfo)
	}
	err = dao.Delete(ctx, messageStruct, m.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	DeleteFileReq := kitexfileservice.DeleteFileReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		FileId: messageStruct.Info[0].FileStorageID,
	}
	_, err = m.serviceClient.FileService.DeleteFile(ctx, &DeleteFileReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	GetGroupUserIDReq := &kitexgroupservice.GetGroupUserIDReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupUserIDResp, err := m.serviceClient.GroupService.GetGroupUserID(ctx, GetGroupUserIDReq)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	//广播
	groupMessageStruct := commonmodel.KafkaNewMessageNotice{
		TraceID:        m.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     GetGroupUserIDResp.UserIdList,
		SessionID:      groupID,
		MessageCode:    commonmodel.MessageCode_GroupWithdrawMessage,
	}
	_, _, err = tool.SendKafkaNewMessageNotice(m.messageTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	return finalErr
}
func (m *MessageService) DeleteMessageAllGroup(ctx context.Context, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:DeleteMessageAllGroup")
	messageStruct := message.NewStruct(groupID, 0, 0, time.Time{}, message.GetWithGroupID)
	err = dao.Delete(ctx, messageStruct, m.dbContext)
	if err != nil {
		return err
	}
	return nil
}
func (m *MessageService) GetMessageList(ctx context.Context, groupID int64, userID int64, startTimeSecond int64, endTimeSecond int64) (messageList []*model.MessageInfo, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:GetMessageList")
	if endTimeSecond < startTimeSecond || endTimeSecond < 0 || startTimeSecond < 0 {
		return nil, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "The Time Formate Is Error", fmt.Errorf("User Send Error Formate Time, startTime:%d ,endTime:%d", startTimeSecond, endTimeSecond), newerror.LevelInfo)
	}
	var finalErr error
	GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupOrSessionRoleAndExistResp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, finalErr
	}
	if !GetGroupOrSessionRoleAndExistResp.Exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Send Message To Unjoin Group"), newerror.LevelInfo)
	}
	startTime := time.Unix(startTimeSecond, 0)
	endTime := time.Unix(endTimeSecond, 0)
	messageStruct := message.NewStruct(groupID, 0, 0, time.Time{}, message.GetWithGroupID, message.GetWithStartAndEndTime(&startTime, &endTime))
	exist, err := dao.Get(ctx, messageStruct, m.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return nil, err
	}
	setLastVisitTimeReq := kitexgroupservice.SetLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: m.traceID},
		UserId:     userID,
		GroupId:    groupID,
	}
	_, err = m.serviceClient.GroupService.SetLastVisitTime(ctx, &setLastVisitTimeReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, finalErr
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusOK, newerror.CodeSuccess, "Do Not Have Message", fmt.Errorf("Do Not Find Message"), newerror.LevelInfo)
	}
	return messageStruct.Info, finalErr
}
func (m *MessageService) GetNewMessage(ctx context.Context, groupID int64, userID int64) (messageID []int64, sendTimeSecond []int64, messageType []string, messageContent []string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:GetNewMessage")
	var finalErr error
	GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupOrSessionRoleAndExistResp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, nil, nil, nil, finalErr
	}
	if !GetGroupOrSessionRoleAndExistResp.Exist {
		return nil, nil, nil, nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Send Message To Unjoin Group"), newerror.LevelInfo)
	}
	getLastVisitTimeReq := &kitexgroupservice.GetLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: m.traceID},
		GroupId:    groupID,
		UserId:     userID,
	}
	getLastVisitTimeResp, err := m.serviceClient.GroupService.GetLastVisitTime(ctx, getLastVisitTimeReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, nil, nil, nil, finalErr
	}
	var LastVisitTimeSecond int64
	for i := range getLastVisitTimeResp.UserIdList {
		if getLastVisitTimeResp.UserIdList[i] == userID {
			LastVisitTimeSecond = getLastVisitTimeResp.LastVisitTimeList[i]
			break
		}
	}
	LastVisitTime := time.Unix(LastVisitTimeSecond, 0)
	messageStruct := message.NewStruct(groupID, 0, 0, time.Time{}, message.GetWithStartAndEndTime(&LastVisitTime, nil), message.GetWithGroupID)
	exist, err := dao.Get(ctx, messageStruct, m.dbContext)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, nil, nil, nil, finalErr
	}
	if !exist {
		return nil, nil, nil, nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "You No Not Have New Message", fmt.Errorf("Try To Get New Message Without Any One New"), newerror.LevelInfo)
	}
	setLastVisitTimeReq := kitexgroupservice.SetLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: m.traceID},
		UserId:     userID,
		GroupId:    groupID,
	}
	_, err = m.serviceClient.GroupService.SetLastVisitTime(ctx, &setLastVisitTimeReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, nil, nil, nil, finalErr
	}
	messageID = make([]int64, 0, len(messageStruct.Info))
	messageType = make([]string, 0, len(messageStruct.Info))
	messageContent = make([]string, 0, len(messageStruct.Info))
	sendTimeSecond = make([]int64, 0, len(messageStruct.Info))
	for _, info := range messageStruct.Info {
		messageID = append(messageID, info.MessageID)
		messageType = append(messageType, info.MessageType)
		messageContent = append(messageContent, info.MessageContent)
		sendTimeSecond = append(sendTimeSecond, info.SendTimeSecond)
	}
	return messageID, sendTimeSecond, messageType, messageContent, nil
}
func (m *MessageService) GetFileContent(ctx context.Context, groupID int64, userID int64, messageID int64) (dataStream []byte, contentType string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:GetFileContent")
	var finalErr error
	GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupOrSessionRoleAndExistResp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return nil, "", finalErr
	}
	if !GetGroupOrSessionRoleAndExistResp.Exist {
		return nil, "", newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Send Message To Unjoin Group"), newerror.LevelInfo)
	}
	messageStruct := message.NewStruct(groupID, messageID, 0, time.Time{}, message.GetWithGroupID, message.GetWithMessageID)
	exist, err := dao.Get(ctx, messageStruct, m.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return nil, "", err
	}
	if !exist || messageStruct.Info[0].FileStorageID == 0 {
		return nil, "", newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The File Is Not Exist", fmt.Errorf("Try To Get Unexist File"), newerror.LevelInfo)
	}
	getFileReq := kitexfileservice.GetFileReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: m.traceID},
		FileId:     messageStruct.Info[0].FileStorageID,
	}
	getFileResp, err := m.serviceClient.FileService.GetFile(ctx, &getFileReq)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return nil, "", err
	}
	return getFileResp.DataStream, getFileResp.ContentType, nil
}
func (m *MessageService) SendGroupNotice(ctx context.Context, groupID int64, userID int64, messageContent string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:SendGroupNotice")
	var finalErr error
	if tool.CalculateLength(messageContent) > m.messageConfig.MaxGroupNoticeLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "The Group Notice Is Too Long", fmt.Errorf("Try To Send Too Long Group Notice"), newerror.LevelInfo)
	}
	GetGroupOrSessionRoleAndExistReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupOrSessionRoleAndExistResp, err := m.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &GetGroupOrSessionRoleAndExistReq)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	if !GetGroupOrSessionRoleAndExistResp.Exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Send Message To Unjoin Group"), newerror.LevelInfo)
	}
	if GetGroupOrSessionRoleAndExistResp.Role == string(commonmodel.Member) {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Do Not Have Enough Permission", fmt.Errorf("Try To Send Group Notice Without Enough Permission"), newerror.LevelInfo)
	}
	GetMuteStatusReq := kitexgroupservice.GetMuteStatusReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetMuteStatusResp, err := m.serviceClient.GroupService.GetMuteStatus(ctx, &GetMuteStatusReq)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	if GetMuteStatusResp.IsMute {
		muteReason := GetMuteStatusResp.MuteReason
		muteEndTime := GetMuteStatusResp.MuteEndTime
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, fmt.Sprintf("You Are In MutIng Because : %s ,Before : %s", muteReason, muteEndTime), fmt.Errorf("Try To Send Message In Muting"), newerror.LevelInfo)
	}
	GetGroupUserIDReq := &kitexgroupservice.GetGroupUserIDReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: m.traceID,
		},
		GroupId: groupID,
		UserId:  userID,
	}
	GetGroupUserIDResp, err := m.serviceClient.GroupService.GetGroupUserID(ctx, GetGroupUserIDReq)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	userIDList := GetGroupUserIDResp.UserIdList
	groupNoticeStruct := commonmodel.KafkaGroupNotice{
		TraceID:        m.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     userIDList,
		SessionID:      groupID,
		MessageCode:    commonmodel.MessageCode_GroupNotice,
	}
	_, _, err = tool.SendKafkaGroupNotice(m.groupNoticeTopic, groupNoticeStruct)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return err
	}
	return finalErr
}

func (m *MessageService) SetOffLineMessage(ctx context.Context, UserAndDeviceID []string, JsonData []byte) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:KeepOffLineMessage")
	MessageID := m.snowFlake.Generate().Int64()
	offlineMessageStruct := offlinemessage.NewStruct(UserAndDeviceID, offlinemessage.WithMessageInfo(MessageID, time.Now().Unix(), JsonData))
	err = dao.Add(ctx, offlineMessageStruct, m.dbContext)
	if err != nil {
		return err
	}
	return nil
}
func (m *MessageService) GetOffLineMessageList(ctx context.Context, UserAndDeviceID string) (JsonDataList [][]byte, exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("message:GetOffLineMessageList")
	offlineMessageStruct := offlinemessage.NewStruct([]string{UserAndDeviceID})
	exist, err = dao.Get(ctx, offlineMessageStruct, m.dbContext)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}
	JsonDataList = make([][]byte, 0, len(offlineMessageStruct.Info))
	for _, info := range offlineMessageStruct.Info {
		JsonDataList = append(JsonDataList, info.JsonData)
	}
	return JsonDataList, true, nil
}
