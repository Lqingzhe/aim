package aichatcontext

import (
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type AiChatContext struct {
	model.MessageContext
	Info *model.MessageContext
}

func NewStruct(userID int64, sumByteLength int64, context []*model.Message) *AiChatContext {
	return &AiChatContext{
		MessageContext: model.MessageContext{
			UserID:        userID,
			SumByteLength: sumByteLength,
			Messages:      context,
		},
		Info: &model.MessageContext{},
	}
}

func (a *AiChatContext) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMongo(ctx, DB, a)
	return err
}
func (a *AiChatContext) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMongo(ctx, DB, a)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (a *AiChatContext) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMongo(ctx, DB, a)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	err = setMongo(ctx, DB, a)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (a *AiChatContext) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMongo(ctx, DB, a)
	if err != nil {
		return err
	}
	return nil
}
