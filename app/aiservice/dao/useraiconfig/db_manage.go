package useraiconfig

import (
	"aim/app/aiservice/model"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
)

type UserAiConfig struct {
	model.BotConfig
	Info *model.BotConfig
}
type operation func(*UserAiConfig)

func NewStruct(UserID int64, Operations ...operation) *UserAiConfig {
	newStruct := &UserAiConfig{
		BotConfig: model.BotConfig{
			UserID: UserID,
		},
		Info: &model.BotConfig{},
	}
	for _, Opration := range Operations {
		Opration(newStruct)
	}
	return newStruct
}
func SetWithModelConfig(ModelName string, BaseUrl string, ApiKey string) operation {
	return func(u *UserAiConfig) {
		u.ModelName = ModelName
		u.BaseUrl = BaseUrl
		u.ApiKey = ApiKey
	}
}
func SetWithModelRoleAndPrompt(Role string, Prompt string) operation {
	return func(u *UserAiConfig) {
		u.Role = Role
		u.Prompt = Prompt
	}
}
func (u *UserAiConfig) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, u)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserAiConfig) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, u)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (u *UserAiConfig) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = updateMysql(ctx, DB, u)
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (u *UserAiConfig) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, u)
	if err != nil {
		return err
	}
	return nil
}
