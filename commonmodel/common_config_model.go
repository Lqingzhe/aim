package commonmodel

import "time"

type CommonConfig struct {
	Service     string                 `yaml:"service"`
	EquipID     int                    `yaml:"equip_id"`
	ServiceInfo map[string]ServiceInfo `yaml:"service_info"`
	NacosConfig NacosConfig            `yaml:"nacos_config"`
}
type GatewayConfig struct {
	Port            string                   `yaml:"port"`
	ReadBufferSize  int                      `yaml:"read_buffer_size"`
	WriteBufferSize int                      `yaml:"write_buffer_size"`
	RoutTimeOut     map[string]time.Duration `yaml:"rout_time_out"`
}
type ServiceInfo struct {
	KitexTimeOut time.Duration `yaml:"kitex_time_out"`
}
type ServiceConfig struct {
	Timeout     time.Duration `yaml:"timeout"`
	ServiceAddr ServiceAddr   `yaml:"service_addr"`
}
type ServiceAddr struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
}
type LimiterConfig struct {
	MaxToken        int64         `yaml:"max_token"`
	GenerateToken   int64         `yaml:"generate_token"`
	RedisExpireTime time.Duration `yaml:"redis_expire_time"`
}

type TokenConfig struct {
	TokenPassword          string        `yaml:"token_password"`
	SaltByteLen            int64         `yaml:"salt_byte_len"`
	RefreshTokenExpireTime time.Duration `yaml:"refresh_token_expire_time"`
	AccessTokenExpireTime  time.Duration `yaml:"access_token_expire_time"`
}
type UserConfig struct {
	SaltByteLen       int64 `yaml:"salt_byte_len"`
	MaxPasswordLength int64 `yaml:"max_password_length"`
	MinPasswordLength int64 `yaml:"min_password_length"`
	MaxUsernameLength int64 `yaml:"max_username_length"`

	MaxUserNameLength  int64 `yaml:"max_user_name_length"`
	MaxIntroduceLength int64 `yaml:"max_introduce_length"`

	MaxNickNameLength int64 `yaml:"max_nick_name_length"`
}
type GroupConfig struct {
	MaxGroupNameLength       int64         `yaml:"max_group_name_length"`
	MaxGroupNickNameLength   int64         `yaml:"max_group_nick_name_length"`
	MaxGroupMuteTime         time.Duration `yaml:"max_group_mute_time"`
	MaxGroupMuteReasonLength int64         `yaml:"max_group_mute_reason_length"`
}
type MessageConfig struct {
	MaxMessageByteLength int64 `yaml:"max_message_byte_length"`
	MaxGroupNoticeLength int64 `yaml:"max_group_notice_length"`
	MaxFileByteLength    int64 `yaml:"max_file_byte_length"`
	MaxVoiceTimeSecond   int64 `yaml:"max_voice_time_second"`
}
type FileConfig struct {
	FileStoragePath    string `yaml:"file_storage_path"`
	MaxPictureDay      int64  `yaml:"max_picture_day"`
	MaxVoiceKeepDay    int64  `yaml:"max_voice_keep_day"`
	MaxFileKeepDay     int64  `yaml:"max_file_keep_day"`
	MaxFileKeepSize    int64  `yaml:"max_file_keep_size"`
	MaxVoiceTimeSecond int64  `yaml:"max_voice_time_second"`
}
type KafkaConfig struct {
	Broker []string `yaml:"broker"`
}
type NacosConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
