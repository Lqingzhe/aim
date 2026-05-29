package commonconfig

import (
	newlog "aim/pkg/log"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func GetKafkaProducerConfig() *sarama.Config {
	config := sarama.Config{}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	return &config
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
