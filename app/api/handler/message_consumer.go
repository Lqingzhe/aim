package handler

import (
	"aim/app/api/service"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func (h *HandlerConfig) Consumer(logger *zap.Logger, poolLimit int64) {
	logger2 := newlog.AddTraceID(logger, "-1")
	ConsumerService := service.NewConsumerService(service.NewWebSocket(h.hub), h.serviceClient.MessageClient)
	errPool := make(chan *newerror.Error, poolLimit)
	taskPool := make(chan func(), poolLimit*2)

	GroupNoticeTopic, err := h.consumer.ConsumePartition("group-notice-topic", 0, sarama.OffsetNewest)
	if err != nil {
		newlog.LogInitFatal(logger2, err, "Init Consumer Error")
	}
	MessageTopic, err := h.consumer.ConsumePartition("message-topic", 0, sarama.OffsetOldest)
	if err != nil {
		newlog.LogInitFatal(logger2, err, "Init Consumer Error")
	}
	SystemTopic, err := h.consumer.ConsumePartition("system-topic", 0, sarama.OffsetNewest)
	if err != nil {
		newlog.LogInitFatal(logger2, err, "Init Consumer Error")
	}
	go func() {
		for {
			select {
			case err2 := <-errPool:
				newlog.Log(newlog.AddError(logger, err2, err2.StatusCode), err2.LogLevel, "Consumer error")
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

	go func() {
		for {
			var msg *sarama.ConsumerMessage
			select {
			case msg = <-GroupNoticeTopic.Messages():
			case msg = <-MessageTopic.Messages():
			case msg = <-SystemTopic.Messages():
			}
			if msg != nil {
				taskPool <- func() {
					traceID, err2 := ConsumerService.Consumer(msg)
					logger := newlog.AddTraceID(logger, traceID)
					if err2 != nil {
						select {
						case errPool <- err2:
						default:
							newlog.Log(newlog.AddError(logger, fmt.Errorf("ErrorPool Full, Drop Error"), -1), newerror.LevelError, "Consumer")
						}
					}
					logger = logger.With(zap.String("message_content", string(msg.Value)))
					newlog.Log(logger, newerror.LevelInfo, "Consumer")
				}
			}
		}
	}()
	ConsumerService.ClearClientOnTime()
}
