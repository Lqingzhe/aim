package tool

import (
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
)

func SendKafkaNewMessageNotice(producer sarama.SyncProducer, message commonmodel.KafkaNewMessageNotice) (partition int32, offset int64, err error) {
	message.MessageType = commonmodel.KafkaMessageType_Message
	result, _ := sonic.Marshal(message)
	kafkaMessage := &sarama.ProducerMessage{
		Topic:     "message-topic",
		Key:       sarama.ByteEncoder(strconv.FormatInt(message.SessionID, 10)),
		Value:     sarama.ByteEncoder(result),
		Timestamp: time.Now(),
	}
	partition, offset, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		err = newerror.MakeError(http.StatusOK, newerror.CodeMessageQueueError, "Send Message Error", err, newerror.LevelError, newerror.WithContinueError)
	}
	return partition, offset, err
}
func SendKafkaGroupNotice(producer sarama.SyncProducer, message commonmodel.KafkaGroupNotice) (partition int32, offset int64, err error) {
	message.MessageType = commonmodel.KafkaMessageType_Notice
	result, _ := sonic.Marshal(message)
	kafkaMessage := &sarama.ProducerMessage{
		Topic:     "group-notice-topic",
		Key:       sarama.ByteEncoder(strconv.FormatInt(message.SessionID, 10)),
		Value:     sarama.ByteEncoder(result),
		Timestamp: time.Now(),
	}
	partition, offset, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		err = newerror.MakeError(http.StatusOK, newerror.CodeMessageQueueError, "Send Message Error", err, newerror.LevelError, newerror.WithContinueError)
	}
	return partition, offset, err
}
func SendKafkaSystemMessage(producer sarama.SyncProducer, message commonmodel.KafkaSystemMessage) (partition int32, offset int64, err error) {
	message.MessageType = commonmodel.KafkaMessageType_System
	result, _ := sonic.Marshal(message)
	kafkaMessage := &sarama.ProducerMessage{
		Topic:     "system-topic",
		Value:     sarama.ByteEncoder(result),
		Timestamp: time.Now(),
	}
	partition, offset, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		err = newerror.MakeError(http.StatusOK, newerror.CodeMessageQueueError, "Send Message Error", err, newerror.LevelError, newerror.WithContinueError)
	}
	return partition, offset, err
}
func SendKafkaAiMessage(producer sarama.SyncProducer, message commonmodel.KafkaAiMessage) (partition int32, offset int64, err error) {
	result, _ := sonic.Marshal(message)
	kafkaMessage := &sarama.ProducerMessage{
		Topic:     "ai-topic",
		Key:       sarama.ByteEncoder(strconv.FormatInt(message.UserID%10, 10)),
		Value:     sarama.ByteEncoder(result),
		Timestamp: time.Now(),
	}
	partition, offset, err = producer.SendMessage(kafkaMessage)
	if err != nil {
		err = newerror.MakeError(http.StatusOK, newerror.CodeMessageQueueError, "Send Message Error", err, newerror.LevelError, newerror.WithContinueError)
	}
	return partition, offset, err
}
