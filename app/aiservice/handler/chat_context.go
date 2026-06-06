package handler

import (
	"aim/app/aiservice/service"
	"aim/kitex_gen/kitexaiservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

// DeleteChatContext implements the AiServiceImpl interface.
func (s *AiServiceImpl) DeleteChatContext(ctx context.Context, req *kitexaiservice.DeleteChatContextReq) (resp *kitexaiservice.DeleteChatContextResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	aiContextService := service.NewChatContext(req.CommonInfo.Trace, s.dbContext)
	err = aiContextService.DeleteChatContext(ctx, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "DeleteChatContext")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "DeleteChatContext")
	return &kitexaiservice.DeleteChatContextResp{}, nil
}
