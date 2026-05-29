package message

import (
	"aim/app/messageservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	model.MessageInfo
	whereWithGroupID     bool
	whereWithUserID      bool
	whereWithMessageID   bool
	whereWithMessageTime bool
	findStartTimeSecond  time.Time
	findEndTimeSecond    time.Time
	Info                 []*model.MessageInfo
}
type operate func(*Message)

func NewStruct(groupID int64, messageID int64, userID int64, sendTime time.Time, Operates ...operate) *Message {
	newStruct := &Message{
		MessageInfo: model.MessageInfo{
			GroupID:        groupID,
			UserID:         userID,
			MessageID:      messageID,
			SendTimeSecond: sendTime.Unix(),
		},
	}
	for _, Operate := range Operates {
		Operate(newStruct)
	}
	return newStruct
}
func WithMessageInfo(messageContent string) operate {
	return func(info *Message) {
		info.MessageType = "message"
		info.MessageContent = messageContent
	}
}
func WithFileInfo(fileStorageID int64, ContentType string) operate {
	return func(info *Message) {
		info.MessageType = "file"
		info.FileStorageID = fileStorageID
		info.ContentType = ContentType
	}
}
func WithPictureInfo(fileStorageID int64, ContentType string) operate {
	return func(info *Message) {
		WithFileInfo(fileStorageID, ContentType)(info)
		info.MessageType = "picture"
	}
}
func WithVoiceInfo(VoiceDurationSecond int64, fileStorageID int64, ContentType string) operate {
	return func(info *Message) {
		WithFileInfo(fileStorageID, ContentType)(info)
		info.MessageType = "voice"
		info.VoiceDurationSecond = VoiceDurationSecond
	}
}
func WithAI(isAi bool) operate {
	return func(info *Message) {
		info.IsAI = isAi
	}
}
func GetWithMessageID(info *Message) {
	info.whereWithMessageID = true
}
func GetWithGroupID(info *Message) {
	info.whereWithGroupID = true
}
func GetWithUserID(info *Message) {
	info.whereWithUserID = true
}
func GetWithStartAndEndTime(startTime time.Time, endTime time.Time) operate {
	return func(info *Message) {
		info.findStartTimeSecond = startTime
		info.findEndTimeSecond = endTime
		info.whereWithMessageTime = true
	}
}
func (m *Message) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMongo(ctx, DB, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Useless Module Unexpectedly Used", fmt.Errorf("%s", "Useless Module Unexpectly Used"), newerror.LevelFatal)
}

func (m *Message) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMongo(ctx, DB, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMongo(ctx, DB, m)
	if err != nil {
		return false, err
	}
	return exist, nil
}
