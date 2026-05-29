package service

import (
	"aim/kitex_gen/kitexcommonmodel"
	kitexmessageservice2 "aim/kitex_gen/kitexmessageservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	newerror "aim/pkg/error"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	*WebSocketStruct
	MessageClient kitexmessageservice.Client
}

func NewConsumerService(websocketStruct *WebSocketStruct, MessageClient kitexmessageservice.Client) *Consumer {
	return &Consumer{
		WebSocketStruct: websocketStruct,
		MessageClient:   MessageClient,
	}
}
func (c *Consumer) Consumer(msg *sarama.ConsumerMessage) (err *newerror.Error) {
	defer func(trace string) {
		err = err.AddErrorTrace(trace)
	}("consumer:Consumer")
	var base struct {
		TraceID    string
		GoalUserID []int64         `json:"goal_user_id"`
		Data       json.RawMessage `json:"data"`
	}
	if err2 := json.Unmarshal(msg.Value, &base); err2 != nil {
		return newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelWarn)
	}
	goalUserAndDeviceIDList := make([]string, 0, len(base.GoalUserID)*3)
	for _, id := range base.GoalUserID {
		for deviceID := range c.hub.Client[id] {
			if !c.WebSocketStruct.PushToUser(id, deviceID, base.Data) {
				goalUserAndDeviceIDList = append(goalUserAndDeviceIDList, strconv.FormatInt(id, 10)+deviceID)
			}
		}
	}
	setOfflineMessageReq := kitexmessageservice2.SetOfflineMessageReq{
		CommonInfo:          &kitexcommonmodel.CommonInfo{Trace: base.TraceID},
		GoalUserAndDeviceId: goalUserAndDeviceIDList,
		JsonData:            base.Data,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err2 := c.MessageClient.SetOfflineMessage(ctx, &setOfflineMessageReq)
	if err2 != nil {
		if a, err := newerror.IsContextError(err2); a {
			return err
		}
		err = newerror.TranslateError(err2)
		return err
	}
	return nil
}

//	func (consumer *Consumer) ConsumeGroupNoticeTopic(msg *sarama.ConsumerMessage) (err *newerror.Error) {
//		defer func(trace string) {
//			err = err.AddErrorTrace(trace)
//		}("consumer:ConsumeGroupNoticeTopic")
//var receive = commonmodel.KafkaGroupNotice{}

//	err2 := sonic.Unmarshal(msg.Value, &receive)
//	if err2 != nil {
//		return newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelWarn)
//	}
//	switch receive.MessageCode {
//
//	}
//}
//func (consumer *Consumer) ConsumeSystemTopic(msg *sarama.ConsumerMessage)(err *newerror.Error)  {
//	defer func(trace string) {
//		err = err.AddErrorTrace(trace)
//	}("consumer:ConsumeSystemTopic")
//	receive := commonmodel.KafkaSystemMessage{}
//	err2 := sonic.Unmarshal(msg.Value, &receive)
//	if err2 != nil {
//		return newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelWarn)
//	}
//	switch receive.MessageCode {
//	//只需要userID的
//	case :
//	}
//}
//func (consumer *Consumer) ConsumeMessageTopic(msg *sarama.ConsumerMessage) (err *newerror.Error) {
//	defer func(trace string) {
//		err = err.AddErrorTrace(trace)
//	}("consumer:ConsumeMessageTopic")
//	receive := commonmodel.KafkaNewMessageNotice{}
//	err2 := sonic.Unmarshal(msg.Value, &receive)
//	if err2 != nil {
//		return newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelWarn)
//	}
//	switch receive.MessageCode {
//
//	}
//}
