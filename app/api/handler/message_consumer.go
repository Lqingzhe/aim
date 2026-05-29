package handler

import (
	"aim/app/api/service"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func (h *HandlerConfig) Consumer(logger *zap.Logger, equipID int64, poolLimit int64) {
	logger = newlog.AddTraceAndEquipID(logger, "-1", equipID)
	ConsumerService := service.NewConsumerService(service.NewWebSocket(h.hub), h.serviceClient.MessageClient)
	errPool := make(chan *newerror.Error, poolLimit)
	taskPool := make(chan func(), poolLimit*2)
	defer close(errPool)
	defer close(taskPool)
	GroupNoticeTopic, err := h.consumer.ConsumePartition("group-notice-topic", 0, sarama.OffsetNewest)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Consumer Error")
	}
	MessageTopic, err := h.consumer.ConsumePartition("message-topic", 0, sarama.OffsetOldest)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Consumer Error")
	}
	SystemTopic, err := h.consumer.ConsumePartition("system-topic", 0, sarama.OffsetNewest)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Consumer Error")
	}
	go func() {
		for {
			select {
			case err := <-errPool:
				newlog.Log(newlog.AddError(logger, err, err.StatusCode), err.LogLevel, "Consumer error")
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
			case msg = <-GroupNoticeTopic.Messages():
			case msg = <-MessageTopic.Messages():
			case msg = <-SystemTopic.Messages():
			}
			if msg != nil {
				taskPool <- func() {
					err := ConsumerService.Consumer(msg)
					if err != nil {
						select {
						case errPool <- err:
						default:
							newlog.Log(newlog.AddError(logger, fmt.Errorf("ErrorPool Full, Drop Error"), -1), newerror.LevelError, "Consumer")
						}
					}
				}
			}
		}
	}()
	ConsumerService.ClearClientOnTime()
}
