package model

import "github.com/cloudwego/eino/schema"

type MessageContext struct {
	UserID        int64 `bson:"_id"`
	SumByteLength int64
	Messages      []*Message
}
type Message struct {
	IsAi      bool
	ToolCalls []schema.ToolCall
	Msg       string
}
type BotConfig struct {
	UserID int64

	ModelName string
	BaseUrl   string
	ApiKey    string

	Role   string
	Prompt string
}

func (_ *BotConfig) Data() {}

type UserProfile struct {
	UserID  int64 `bson:"_id"`
	Profile string
}
