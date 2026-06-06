package service

import (
	"aim/kitex_gen/kitexcommonmodel"
	kitexmessageservice2 "aim/kitex_gen/kitexmessageservice"
	"aim/kitex_gen/kitexmessageservice/kitexmessageservice"
	newerror "aim/pkg/error"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
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
func (c *Consumer) Consumer(msg *sarama.ConsumerMessage) (traceID string, err *newerror.Error) {
	defer func(trace string) {
		if err != nil {
			err = err.AddErrorTrace(trace).(*newerror.Error)
		}
		err2 := recover()
		if err2 != nil {
			err = newerror.TranslateError(newerror.MakeKafkaError(newerror.CodeInternalError, fmt.Errorf("%v", err2), newerror.LevelError)).AddErrorTrace(trace).(*newerror.Error)
		}
	}("consumer:Consumer")
	base := make(map[string]any)
	if err2 := sonic.Unmarshal(msg.Value, &base); err2 != nil {
		return "", newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelWarn)
	}
	var ok bool
	goalUserIDInterface, ok := base["goal_user_id"].([]interface{})
	if !ok {
		return "", newerror.MakeKafkaError(newerror.CodeInternalError, fmt.Errorf("goal_user_id Is Not Exist Or Not Array"), newerror.LevelError)
	}
	delete(base, "goal_user_id")

	// 转换为 []int64
	goalUserID := make([]int64, 0, len(goalUserIDInterface))
	for _, v := range goalUserIDInterface {
		switch val := v.(type) {
		case float64:
			goalUserID = append(goalUserID, int64(val))
		case int64:
			goalUserID = append(goalUserID, val)
		case int:
			goalUserID = append(goalUserID, int64(val))
		default:
			return "", newerror.MakeKafkaError(newerror.CodeInternalError, fmt.Errorf("goal_user_id element not a number"), newerror.LevelError)
		}
	}
	traceID, ok = base["trace_id"].(string)
	if !ok {
		return "", newerror.MakeKafkaError(newerror.CodeInternalError, fmt.Errorf("trace_id Is Not Exist Or []int64"), newerror.LevelError)
	}
	delete(base, "goal_user_id")
	delete(base, "trace_id")
	goalUserAndDeviceIDList := make([]string, 0, len(goalUserID)*3)
	Data, err2 := sonic.Marshal(base)
	if err2 != nil {
		return traceID, newerror.MakeKafkaError(newerror.CodeInvalidJSON, err2, newerror.LevelError)
	}
	for _, id := range goalUserID {
		for deviceID := range c.hub.Client[id] {
			if !c.WebSocketStruct.PushToUser(id, deviceID, Data) {
				goalUserAndDeviceIDList = append(goalUserAndDeviceIDList, strconv.FormatInt(id, 10)+deviceID)
			}
		}
	}
	if len(goalUserAndDeviceIDList) > 0 {
		setOfflineMessageReq := kitexmessageservice2.SetOfflineMessageReq{
			CommonInfo:          &kitexcommonmodel.CommonInfo{Trace: traceID},
			GoalUserAndDeviceId: goalUserAndDeviceIDList,
			JsonData:            Data,
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_, err2 := c.MessageClient.SetOfflineMessage(ctx, &setOfflineMessageReq)
		if err2 != nil {
			if a, err := newerror.IsContextError(newerror.UnMarshalError(err2)); a {
				return traceID, err
			}
			err = newerror.TranslateError(err2)
			return traceID, err
		}
	}
	return traceID, nil
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
