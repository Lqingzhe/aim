package config

import (
	"aim/app/api/model"
	"aim/commonmodel"
	"aim/pkg/config"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"reflect"
)

func InitConfig() *model.Config {
	data := commonconfig.OpenYaml()

	rawMap := make(map[string]interface{})
	rawData := `
common_config:
  service: "api-gateway"
  equip_id: 1
  service_info:
    user_service:
      kitex_time_out: 5s
    group_service:
      kitex_time_out: 5s
    message_service:
      kitex_time_out: 5s
    file_service:
      kitex_time_out: 5s
      service_addr:
  nacos_config:
    host: "localhost"
    port: 8848
  kafka_config:
    broker:
      - "localhost:9092"
      - "localhost:9093"

gateway_config:
  port: "8888"
  read_buffer_size: 1024
  write_buffer_size: 1024

  rout_time_out:
    /user/register: 10s
    /user/login: 10s
    /user/refresh-token: 10s
    /user/logout-all-device: 5s
    /user/logout-a-device: 5s
    /user/get-user-info: 5s
    /user/get-other-user-info: 5s
    /user/update-user-info: 5s
    /user/remark: 5s
    /ws: 30s
    /group/create-group: 10s
    /group/delete-group: 10s
    /group/search-group: 5s
    /group/set-group-apply: 10s
    /group/get-group-apply-list: 5s
    /group/agree-group-apply: 10s
    /group/refuse-group-apply: 10s
    /group/leave-group: 10s
    /group/get-group-info: 5s
    /group/change-group-info: 10s
    /group/get-group-info-with-user: 5s
    /group/update-group-info-with-user: 5s
    /group/get-group-and-session-id: 5s
    /group/transform-group-owner: 10s
    /group/set-manager: 10s
    /group/revoke-manager: 10s
    /group/get-last-visit-time: 5s
    /group/kick-out-group: 10s
    /group/apply-for-friend: 10s
    /group/get-friend-apply-list: 5s
    /group/refuse-friend-apply: 10s
    /group/creat-session: 10s
    /group/delete-session: 10s
    /group/get-friend-last-visit-time: 5s
    /group/set-mute: 10s
    /group/release-mute: 10s
    /message/send-message: 10s
    /message/send-file: 30s
    /message/send-voice: 30s
    /message/send-picture: 30s
    /message/withdraw-message: 10s
    /message/get-message-list: 10s
    /message/get-new-message: 10s
    /message/get-file-content: 30s
    /message/send-group-notice: 10s

limiter_config:
  max_token: 100
  generate_token: 10
  redis_expire_time: 60s

token_config:
  token_password: "your-jwt-secret-key-change-in-production"
  salt_byte_len: 16
  refresh_token_expire_time: 720h
  access_token_expire_time: 1h


redis_config:
  addr: "localhost:6379"
  password: "123456"
  pool_size: 100
  min_idle_conns: 10
  max_idle_conns: 50
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  conn_max_lifetime: 3600s
  conn_max_idle_time: 600s
`
	data = []byte(rawData)
	log.Printf("%s\n", data)
	newConfig2 := struct {
		CommonConfig commonmodel.CommonConfig `yaml:"common_config"`
	}{}
	if err := yaml.Unmarshal(data, &newConfig2); err != nil {
		log.Fatalf("Unmarshal Common Config Failed: %v", err)
	}
	log.Printf("newConfig2:%s\n", newConfig2)

	type LocalKafkaConfig struct {
		Broker []string `yaml:"broker"`
	}

	type LocalCommonConfig struct {
		Service     string                 `yaml:"service"`
		EquipID     int                    `yaml:"equip_id"`
		ServiceInfo map[string]interface{} `yaml:"service_info"`
		NacosConfig struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"nacos_config"`
		KafkaConfig LocalKafkaConfig `yaml:"kafka_config"`
	}

	// 反射对比两个结构体
	fmt.Println("=== 结构体对比开始 ===")

	localType := reflect.TypeOf(LocalCommonConfig{})
	remoteType := reflect.TypeOf(commonmodel.CommonConfig{})

	for i := 0; i < localType.NumField(); i++ {
		localField := localType.Field(i)
		remoteField, ok := remoteType.FieldByName(localField.Name)

		if !ok {
			fmt.Printf("❌ commonmodel 缺少字段：%s\n", localField.Name)
			continue
		}

		localTag := localField.Tag.Get("yaml")
		remoteTag := remoteField.Tag.Get("yaml")

		if localTag != remoteTag {
			fmt.Printf("❌ 字段 %s 标签不匹配：本地=%q，远程=%q\n", localField.Name, localTag, remoteTag)
		} else {
			fmt.Printf("✅ 字段 %s 标签匹配：%q\n", localField.Name, localTag)
		}
	}

	// 专门对比 KafkaConfig 结构体
	fmt.Println("\n=== KafkaConfig 对比 ===")
	localKafkaType := reflect.TypeOf(LocalKafkaConfig{})
	remoteKafkaType := reflect.TypeOf(commonmodel.KafkaConfig{})

	for i := 0; i < localKafkaType.NumField(); i++ {
		localField := localKafkaType.Field(i)
		remoteField, ok := remoteKafkaType.FieldByName(localField.Name)

		if !ok {
			fmt.Printf("❌ commonmodel.KafkaConfig 缺少字段：%s\n", localField.Name)
			continue
		}

		localTag := localField.Tag.Get("yaml")
		remoteTag := remoteField.Tag.Get("yaml")

		if localTag != remoteTag {
			fmt.Printf("❌ 字段 %s 标签不匹配：本地=%q，远程=%q\n", localField.Name, localTag, remoteTag)
		} else {
			fmt.Printf("✅ 字段 %s 标签匹配：%q\n", localField.Name, localTag)
		}
	}

	fmt.Println("=== 结构体对比结束 ===")
	Local := LocalCommonConfig{}
	yaml.Unmarshal(data, &Local)
	log.Printf("local:%s\n", Local)
	yaml.Unmarshal(data, rawMap)
	log.Printf("map:%s\n", rawMap)

	newConfig := &model.Config{}

	newConfig.CommonConfig = commonconfig.GetCommonConfig(data)
	log.Printf("newConfig:%s\n", newConfig.CommonConfig)
	_ = fmt.Sprintf("%p", &newConfig.CommonConfig.KafkaConfig.Broker)

	newConfig.GatewayConfig = commonconfig.GetGatewayConfig(data)
	newConfig.DBConfig = initDBConfig(data)
	newConfig.KafkaConfig = commonconfig.GetKafkaConfig(data)
	newConfig.LimiterConfig = commonconfig.GetLimitersConfig(data)
	newConfig.TokenConfig = commonconfig.GetTokenConfig(data)

	return newConfig
}
