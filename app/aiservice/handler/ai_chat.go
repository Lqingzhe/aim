package handler

import (
	"aim/app/aiservice/service"
	"aim/kitex_gen/kitexaiservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// SendMessageToAi implements the AiServiceImpl interface.
func (s *AiServiceImpl) SendMessageToAi(ctx context.Context, req *kitexaiservice.SendMessageToAiReq) (resp *kitexaiservice.SendMessageToAiResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	aiChatService := service.NewAiChat(req.CommonInfo.Trace, s.aiTopic, s.dbContext, s.aiConfig, s.serviceClient, s.Limiter, s.TraceWithUserManager, s.tools)
	err = aiChatService.SendMessageToAi(ctx, req.UserId, req.GroupId, req.Message)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "SendMessageToAi")
		return nil, err
	}
	return
}
func (s *AiServiceImpl) SendMessageToUser(poolLimit int64) {
	logger2 := newlog.AddTraceID(s.logger, "-1")
	serviceStruct := service.NewAiChat("", s.aiTopic, s.dbContext, s.aiConfig, s.serviceClient, s.Limiter, s.TraceWithUserManager, s.tools)
	errPool := make(chan *newerror.Error, poolLimit)
	taskPool := make(chan func(), poolLimit*2)

	aiChatTopic, err := s.consumer.ConsumePartition("ai-topic", 0, sarama.OffsetNewest)
	if err != nil {
		newlog.LogInitFatal(logger2, err, "Init Consumer Error")
	}
	go func() {
		for {
			select {
			case err2 := <-errPool:
				newlog.Log(newlog.AddError(s.logger, err2, err2.StatusCode), err2.LogLevel, "SendMessageToUser")
			}
		}
	}()
	for range poolLimit {
		go func() {
			for {
				select {
				case task := <-taskPool:
					task()
				}
			}
		}()
	}
	var msg *sarama.ConsumerMessage
	go func() {
		for {
			select {
			case msg = <-aiChatTopic.Messages():
			}
			if msg == nil {
				taskPool <- func() {
					ctx, cancel := context.WithTimeout(context.Background(), s.aiConfig.AiChatTimeout)
					traceID, err2 := serviceStruct.SendMessageToUser(ctx, msg)
					logger := newlog.AddTraceID(s.logger, traceID)
					if err2 != nil {
						select {
						case errPool <- newerror.TranslateError(err2):
						default:
							newlog.Log(newlog.AddError(logger, fmt.Errorf("ErrorPool Full, Drop Error"), -1), newerror.LevelError, "SendMessageToUser")
						}
					}
					cancel()
					newlog.Log(logger, newerror.LevelInfo, "SendMessageToUser")
				}
			}
		}
	}()
}
