package commonconfig

import (
	newlog "aim/pkg/log"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func GetKafkaProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.MaxMessageBytes = 1048576
	config.Producer.Timeout = 5 * time.Second

	config.Net.MaxOpenRequests = 10
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	config.Admin.Timeout = 10 * time.Second

	config.Consumer.Fetch.Min = 1
	config.Consumer.MaxWaitTime = 500 * time.Millisecond
	return config
}
func MakeKafkaProducer(brokers []string, config *sarama.Config, logger *zap.Logger) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Kafka Producer Failed")
		return nil
	}
	return producer
}

func GetKafkaConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()

	return config
}
func MakeKafkaConsumer(brokers []string, config *sarama.Config, logger *zap.Logger) sarama.Consumer {
	if brokers == nil {
		newlog.LogInitFatal(logger, fmt.Errorf("Kafka Consumer Broker Is nil"), "Init Kafka Consumer Failed")
	}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		newlog.LogInitFatal(logger, err, "Init Kafka Consumer Failed")
	}
	return consumer
}
