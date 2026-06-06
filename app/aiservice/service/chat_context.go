package service

import (
	"aim/app/aiservice/dao"
	"aim/app/aiservice/dao/aichatcontext"
	"aim/app/aiservice/model"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	"context"
)

type ChatContext struct {
	traceID       string
	dbContext     *model.DBContext
	serviceClient model.ServiceClient
}

func NewChatContext(traceID string, dbContext *model.DBContext) *ChatContext {
	return &ChatContext{
		traceID:   traceID,
		dbContext: dbContext,
	}
}
func (c *ChatContext) DeleteChatContext(ctx context.Context, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("chatcontext:DeleteChatContext")
	var finalErr error
	kitexReq := kitexgroupservice.GetGroupOrSessionRoleAndExistReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.traceID},
		UserId:     userID,
	}
	_, err = c.serviceClient.GroupService.GetGroupOrSessionRoleAndExist(ctx, &kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	aiChatContextStruct := aichatcontext.NewStruct(userID, 0, nil)
	err = dao.Delete(ctx, aiChatContextStruct, c.dbContext)
	if err != nil {
		return err
	}
	return nil
}
