package handler

import (
	"aim/app/aiservice/service"
	"aim/kitex_gen/kitexaiservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

// GetAiConfig implements the AiServiceImpl interface.
func (s *AiServiceImpl) GetAiConfig(ctx context.Context, req *kitexaiservice.GetAiConfigReq) (resp *kitexaiservice.GetAiConfigResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	userAiConfigService := service.NewUserAiConfig(s.dbContext)
	modelName, baseUrl, ApiKey, Role, Prompt, err := userAiConfigService.GetAiConfig(ctx, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetAiConfig")
		return nil, err
	}
	kitesResp := &kitexaiservice.GetAiConfigResp{
		ModelName: modelName,
		BaseUrl:   baseUrl,
		ApiKey:    ApiKey,
		Role:      Role,
		Prompt:    Prompt,
	}
	newlog.Log(logger, newerror.LevelInfo, "GetAiConfig")
	return kitesResp, nil
}

// UpdateAiConfig implements the AiServiceImpl interface.
func (s *AiServiceImpl) UpdateAiConfig(ctx context.Context, req *kitexaiservice.UpdateAiConfigReq) (resp *kitexaiservice.UpdateAiConfigResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	userAiConfigService := service.NewUserAiConfig(s.dbContext)
	err = userAiConfigService.UpdateAiConfig(ctx, req.UserId, req.ModelName, req.BaseUrl, req.ApiKey, req.Role, req.Prompt)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "UpdateAiConfig")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "UpdateAiConfig")
	return &kitexaiservice.UpdateAiConfigResp{}, nil
}

// DeleteAiConfig implements the AiServiceImpl interface.
func (s *AiServiceImpl) DeleteAiConfig(ctx context.Context, req *kitexaiservice.DeleteAiConfigReq) (resp *kitexaiservice.DeleteAiConfigResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	userAiConfigService := service.NewUserAiConfig(s.dbContext)
	err = userAiConfigService.DeleteAiConfig(ctx, req.UserId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "DeleteAiConfig")
		return nil, err
	}
	newlog.Log(logger, newerror.LevelInfo, "DeleteAiConfig")
	return &kitexaiservice.DeleteAiConfigResp{}, nil
}
