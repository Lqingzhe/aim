package service

import (
	myagent "aim/app/aiservice/agent"
	"aim/app/aiservice/dao"
	"aim/app/aiservice/dao/aichatcontext"
	"aim/app/aiservice/dao/useraiconfig"
	"aim/app/aiservice/model"
	"aim/commonmodel"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	agenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"golang.org/x/time/rate"
)

type AiChat struct {
	traceID              string
	aiTopic              sarama.SyncProducer
	dbContext            *model.DBContext
	aiConfig             commonmodel.AiConfig
	serviceClient        model.ServiceClient
	Limiter              *rate.Limiter
	TraceWithUserManager *myagent.TraceWithUserManager
	tools                []agenttool.BaseTool
}

func NewAiChat(traceID string, aiTopic sarama.SyncProducer, dbContext *model.DBContext, aiConfig commonmodel.AiConfig, serviceClient model.ServiceClient, Limiter *rate.Limiter, TraceWithUserManager *myagent.TraceWithUserManager, tools []agenttool.BaseTool) *AiChat {
	return &AiChat{
		traceID:              traceID,
		aiTopic:              aiTopic,
		dbContext:            dbContext,
		aiConfig:             aiConfig,
		serviceClient:        serviceClient,
		Limiter:              Limiter,
		TraceWithUserManager: TraceWithUserManager,
		tools:                tools,
	}
}
func (a *AiChat) SendMessageToAi(ctx context.Context, userID int64, groupID int64, message string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("aichat:MakeKafkaAiMessage")
	select {
	case <-ctx.Done():
		_, err = newerror.IsContextError(ctx.Err())
		return err
	default:
		{
			if a.Limiter.Allow() {
				userAiConfigStruct := useraiconfig.NewStruct(userID)
				exist, err := dao.Get(ctx, userAiConfigStruct, a.dbContext)
				if err != nil {
					err2 := newerror.TranslateError(err)
					err2.IsNeedInterrupt = false
					return err2
				}
				if !exist {
					return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "You Did Not Have Ai Chat Config", fmt.Errorf("Try To Use Ai Chat Without Set Config"), newerror.LevelInfo, newerror.WithContinueError)
				}
				messageStruct := commonmodel.KafkaAiMessage{
					UserID:  userID,
					GroupID: groupID,
					Message: message,
					TraceID: a.traceID,
				}
				_, _, err = tool.SendKafkaAiMessage(a.aiTopic, messageStruct)
				if err != nil {
					return err
				}
				return nil
			} else {
				return newerror.MakeError(http.StatusOK, newerror.CodeServiceBusy, "Ai Service Is Busy, Please Repeat Later", fmt.Errorf("Ai Service Too Busy And Touch Limiter"), newerror.LevelWarn, newerror.WithContinueError)
			}
		}
	}
}
func (a *AiChat) SendMessageToUser(ctx context.Context, msg *sarama.ConsumerMessage) (traceID string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("aichat:SendMessage")
	var finalErr error
	data := &commonmodel.KafkaAiMessage{}
	err = sonic.Unmarshal(msg.Value, &data)
	if err != nil {
		return "", newerror.MakeKafkaError(newerror.CodeInvalidJSON, err, newerror.LevelWarn)
	}
	traceID = data.TraceID
	var message string
	defer func() {
		sendMessageReq := kitexmessageservice.SendMessageReq{
			CommonInfo: &kitexcommonmodel.CommonInfo{
				Trace: traceID,
			},
			GroupId:        data.GroupID,
			UserId:         data.UserID,
			MessageContent: "@" + strconv.FormatInt(data.UserID, 10) + " " + message,
			IsAi:           true,
		}
		_, err = a.serviceClient.MessageService.SendMessage(ctx, &sendMessageReq)
		err = newerror.UnMarshalError(err)
	}()
	a.TraceWithUserManager.SetTraceID(data.UserID, traceID)
	defer a.TraceWithUserManager.ReleaseTrace(data.UserID)
	userAiConfigStruct := useraiconfig.NewStruct(data.UserID)
	exist, err := dao.Get(ctx, userAiConfigStruct, a.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	if !exist {
		message = "You Did Not Have Ai Chat Config"
		return data.TraceID, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "You Did Not Have Ai Chat Config", fmt.Errorf("Try To Use Ai Chat Without Set Config"), newerror.LevelInfo, newerror.WithContinueError)
	}
	Agent, err := myagent.CreateAiAgent(ctx, userAiConfigStruct.ModelName, userAiConfigStruct.BaseUrl, userAiConfigStruct.ApiKey, a.tools, int(a.aiConfig.MaxThinkStep))
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	messageFormate := myagent.CreateFormate(userAiConfigStruct.Role, userAiConfigStruct.Prompt)
	aiChatContextStruct := aichatcontext.NewStruct(data.UserID, 0, nil)
	exist, err = dao.Get(ctx, aiChatContextStruct, a.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	var history []*schema.Message
	if !exist {
		history = myagent.CreateHistory(int(a.aiConfig.MaxChatTurns))
	} else {
		history = myagent.TranslateHistoryToContext(aiChatContextStruct.Info.Messages)
	}
	chatMessage, err := myagent.CreateMessage(ctx, messageFormate, data.Message, history)
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	aiMessage, err := myagent.AgentGenerateChat(ctx, chatMessage, Agent, data.UserID)
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	message = aiMessage.Content
	history, addByteLength := myagent.AddHistory(history, data.Message, aiMessage)

	history, releaseByteLength := myagent.CleanHistory(int(a.aiConfig.MaxChatTurns), history, aiChatContextStruct.Info.SumByteLength+addByteLength, a.aiConfig.MaxMessageByteLength)
	aiChatContextStruct = aichatcontext.NewStruct(data.UserID, aiChatContextStruct.Info.SumByteLength+addByteLength-releaseByteLength, myagent.TranslateContextToHistory(history))
	exist, err = dao.Update(ctx, aiChatContextStruct, a.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		err2 := newerror.TranslateError(err)
		message = err2.HttpMessage
		return data.TraceID, finalErr
	}
	if !exist {
		err = dao.Add(ctx, aiChatContextStruct, a.dbContext)
		if newerror.WhetherInterrupt(err, &finalErr) {
			err2 := newerror.TranslateError(err)
			message = err2.HttpMessage
			return data.TraceID, finalErr
		}
	}
	return data.TraceID, finalErr
} //不需要rpc调用
