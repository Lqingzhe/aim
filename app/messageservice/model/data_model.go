package model

type MessageInfo struct {
	GroupID   int64 `bson:"group_id"`
	MessageID int64 `bson:"_id"`
	UserID    int64 `bson:"user_id"`

	MessageContent string `bson:"message_content,omitempty"`

	FileStorageID       int64  `bson:"file_storage_id,omitempty"`
	ContentType         string `bson:"content_type,omitempty"`
	VoiceDurationSecond int64  `bson:"voice_duration_second,omitempty"`

	IsAI bool

	MessageType    string `bson:"message_type,omitempty"`
	SendTimeSecond int64  `bson:"send_time_second,omitempty"`
}
type OfflineMessageInfo struct {
	UserAndDeviceID []string `gorm:"-"`
	MessageID       int64    `gorm:"NOT NULL;primaryKey;comment:message_id"`
	JsonData        []byte   `gorm:"text;comment:json_data"`
	SendTimeSecond  int64    `gorm:"NOT NULL;comment:send_time_second"`
}

func (_ *OfflineMessageInfo) Data() {}
