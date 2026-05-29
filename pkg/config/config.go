package commonconfig

import (
	"aim/commonmodel"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func OpenYaml() (data []byte) {
	basePath, _ := os.Getwd()
	data, err := os.ReadFile(filepath.Join(basePath, "config.yaml"))
	if err != nil {
		log.Fatalf("Read File config.yaml Failed: %v", err)
	}
	return data
}
func GetCommonConfig(data []byte) commonmodel.CommonConfig {
	newConfig := struct {
		CommonConfig commonmodel.CommonConfig `yaml:"common_config"`
	}{}

	newConfig.CommonConfig.KafkaConfig.Broker = make([]string, 0, 100)

	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		log.Fatalf("Unmarshal Common Config Failed: %v", err)
	}
	log.Printf("newConfig:%s\n", newConfig)
	log.Printf("2newConfig:%+v\n", newConfig.CommonConfig)
	return newConfig.CommonConfig
}
func GetMysqlConfig(data []byte) commonmodel.MysqlConfig {
	newStruct := struct {
		MysqlConfig commonmodel.MysqlConfig `yaml:"mysql_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Mysql Config Failed: %v", err)
	}
	newStruct.MysqlConfig.Url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true&timeout=%s&readTimeout=%s&writeTimeout=%s",
		newStruct.MysqlConfig.Username,
		newStruct.MysqlConfig.Password,
		newStruct.MysqlConfig.Host,
		newStruct.MysqlConfig.Port,
		newStruct.MysqlConfig.DBName,
		newStruct.MysqlConfig.Timeout.String(),
		newStruct.MysqlConfig.ReadTimeout.String(),
		newStruct.MysqlConfig.WriteTimeout.String(),
	)
	return newStruct.MysqlConfig
}
func GetRedisConfig(data []byte) commonmodel.RedisConfig {
	newStruct := struct {
		RedisConfig commonmodel.RedisConfig `yaml:"redis_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Redis Config Failed: %v", err)
	}
	return newStruct.RedisConfig
}
func GetMongoDBConfig(data []byte) commonmodel.MongoDBConfig {
	newStruct := struct {
		MongoDBConfig commonmodel.MongoDBConfig `yaml:"mongo_db_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal MongoDB Config Failed: %v", err)
	}
	return newStruct.MongoDBConfig
}

func GetGatewayConfig(data []byte) commonmodel.GatewayConfig {
	newStruct := struct {
		GatewayConfig commonmodel.GatewayConfig `yaml:"gateway_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Gateway Config Failed: %v", err)
	}
	return newStruct.GatewayConfig
}
func GetLimitersConfig(data []byte) commonmodel.LimiterConfig {
	newStruct := struct {
		LimiterConfig commonmodel.LimiterConfig `yaml:"limiter_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Limiter Config Failed: %v", err)
	}
	return newStruct.LimiterConfig
}
func GetTokenConfig(data []byte) commonmodel.TokenConfig {
	newStruct := struct {
		TokenConfig commonmodel.TokenConfig `yaml:"token_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Token Config Failed: %v", err)
	}
	return newStruct.TokenConfig
}

func GetServiceConfig(data []byte) commonmodel.ServiceConfig {
	newStruct := struct {
		ServiceConfig commonmodel.ServiceConfig `yaml:"service_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Service Config Failed: %v", err)
	}
	return newStruct.ServiceConfig
}
func GetUserConfig(data []byte) commonmodel.UserConfig {
	newStruct := struct {
		UserConfig commonmodel.UserConfig `yaml:"user_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal User Config Failed: %v", err)
	}
	if newStruct.UserConfig.MaxUsernameLength > 255 {
		newStruct.UserConfig.MaxUsernameLength = 255
	}
	if newStruct.UserConfig.MaxPasswordLength > 255 {
		newStruct.UserConfig.MaxPasswordLength = 255
	}
	if newStruct.UserConfig.MaxUserNameLength > 255 {
		newStruct.UserConfig.MaxUserNameLength = 255
	}
	if newStruct.UserConfig.MaxIntroduceLength > 255 {
		newStruct.UserConfig.MaxIntroduceLength = 255
	}
	if newStruct.UserConfig.MaxNickNameLength > 255 {
		newStruct.UserConfig.MaxNickNameLength = 255
	}
	return newStruct.UserConfig
}
func GetGroupConfig(data []byte) commonmodel.GroupConfig {
	newStruct := struct {
		GroupConfig commonmodel.GroupConfig `yaml:"group_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Group Config Failed: %v", err)
	}
	if newStruct.GroupConfig.MaxGroupNameLength > 255 {
		newStruct.GroupConfig.MaxGroupNameLength = 255
	}
	if newStruct.GroupConfig.MaxGroupNickNameLength > 255 {
		newStruct.GroupConfig.MaxGroupNickNameLength = 255
	}
	if newStruct.GroupConfig.MaxGroupMuteReasonLength > 255 {
		newStruct.GroupConfig.MaxGroupMuteReasonLength = 255
	}
	return newStruct.GroupConfig
}
func GetMessageConfig(data []byte) commonmodel.MessageConfig {
	newStruct := struct {
		MessageConfig commonmodel.MessageConfig `yaml:"message_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Message Config Failed: %v", err)
	}
	if newStruct.MessageConfig.MaxMessageByteLength > 65535 {
		newStruct.MessageConfig.MaxMessageByteLength = 65535
	}
	return newStruct.MessageConfig
}
func GetFileConfig(data []byte) commonmodel.FileConfig {
	newStruct := struct {
		FileConfig commonmodel.FileConfig `yaml:"file_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal File Config Failed: %v", err)
	}
	return newStruct.FileConfig
}
func GetKafkaConfig(data []byte) commonmodel.KafkaConfig {
	newStruct := struct {
		KafkaConfig commonmodel.KafkaConfig `yaml:"kafka_config"`
	}{}
	if err := yaml.Unmarshal(data, &newStruct); err != nil {
		log.Fatalf("Unmarshal Kafka Config Failed: %v", err)
	}
	return newStruct.KafkaConfig
}
