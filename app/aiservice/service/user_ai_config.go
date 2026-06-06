package service

import (
	"aim/app/aiservice/agent"
	"aim/app/aiservice/dao"
	"aim/app/aiservice/dao/useraiconfig"
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"context"
	"fmt"
	"net/http"
)

type UserAiConfig struct {
	dbContext *model.DBContext
}

func NewUserAiConfig(dbContext *model.DBContext) *UserAiConfig {
	return &UserAiConfig{
		dbContext: dbContext,
	}
}
func (u *UserAiConfig) GetAiConfig(ctx context.Context, userID int64) (ModelName string, BaseUrl string, ApiKey string, Role string, Prompt string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("useraiconfig:SetAiConfig")
	var finalErr error
	userAiConfigStruct := useraiconfig.NewStruct(userID)
	exist, err := dao.Get(ctx, userAiConfigStruct, u.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return "", "", "", "", "", finalErr
	}
	if !exist {
		err = dao.Add(ctx, userAiConfigStruct, u.dbContext)
		if newerror.WhetherInterrupt(err, &finalErr) {
			return "", "", "", "", "", finalErr
		}
		return "", "", "", "", "", nil
	}
	Info := userAiConfigStruct.Info
	return Info.ModelName, Info.BaseUrl, Info.ApiKey, Info.Role, Info.Prompt, nil
}
func (u *UserAiConfig) UpdateAiConfig(ctx context.Context, userID int64, ModelName string, BaseUrl string, ApiKey string, Role string, Prompt string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("useraiconfig:UpdateAiConfig")
	var finalErr error
	if ModelName == "" || BaseUrl == "" || ApiKey == "" || Role == "" || Prompt == "" {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeMissingParam, "Lack Necessary Param", fmt.Errorf("Lack Necessary Ai Config"), newerror.LevelInfo)
	}
	_, err = agent.CreateAiAgent(ctx, ModelName, BaseUrl, ApiKey, nil, 0)
	if err != nil {
		return err
	}
	userAiConfigStruct := useraiconfig.NewStruct(userID, useraiconfig.SetWithModelConfig(ModelName, BaseUrl, ApiKey), useraiconfig.SetWithModelRoleAndPrompt(Role, Prompt))
	exist, err := dao.Update(ctx, userAiConfigStruct, u.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return finalErr
	}
	if !exist {
		err = dao.Add(ctx, userAiConfigStruct, u.dbContext)
		if newerror.WhetherInterrupt(err, &finalErr) {
			return finalErr
		}
	}
	return nil
}
func (u *UserAiConfig) DeleteAiConfig(ctx context.Context, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("useraiconfig:DeleteAiConfig")
	userAiConfigStruct := useraiconfig.NewStruct(userID)
	err = dao.Delete(ctx, userAiConfigStruct, u.dbContext)
	if err != nil {
		return err
	}
	return nil
}
