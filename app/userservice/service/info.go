package service

import (
	"aim/app/userservice/dao"
	"aim/app/userservice/dao/nickname"
	"aim/app/userservice/dao/userinfo"
	"aim/app/userservice/model"
	"aim/commonmodel"
	"aim/kitex_gen/kitexcommonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"fmt"
	"net/http"

	"context"
)

type Info struct {
	dbContext  *model.DBContext
	UserInfo   *model.UserInfo
	RemarkInfo *model.RemarkInfo
}

func NewUserInfo(dbContext *model.DBContext, userInfo *model.UserInfo, remarkInfo *model.RemarkInfo) *Info {
	return &Info{
		dbContext:  dbContext,
		UserInfo:   userInfo,
		RemarkInfo: remarkInfo,
	}
}
func TranslateUserInfoModel(result *model.UserInfo) *kitexcommonmodel.UserInfo {
	return &kitexcommonmodel.UserInfo{
		UserID:        result.UserID,
		UserName:      result.UserName,
		Introduction:  result.Introduction,
		BirthdayYear:  result.BirthdayYear,
		BirthdayMonth: result.BirthdayMonth,
		BirthdayDay:   result.BirthdayDay,
	}
}
func (i *Info) GetUserInfo(ctx context.Context) (respUserInfo *model.UserInfo, respRemarkInfo []*model.RemarkInfo, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("info:GetUserInfo")
	remarkInfo := nickname.NewStruct(i.UserInfo.UserID, 0, "")
	userInfo := userinfo.NewStruct(i.UserInfo.UserID)
	exist, err := dao.Get(ctx, remarkInfo, i.dbContext)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		respRemarkInfo = nil
	} else {
		respRemarkInfo = remarkInfo.Info
	}
	exist, err = dao.Get(ctx, userInfo, i.dbContext)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		err = dao.Add(ctx, userInfo, i.dbContext)
		if err != nil {
			return nil, nil, err
		}
		respUserInfo = nil
	} else {
		respUserInfo = userInfo.Info
	}
	return respUserInfo, respRemarkInfo, nil
}
func (i *Info) GetOtherUserInfo(ctx context.Context) (respUserInfo *model.UserInfo, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("info:GerOtherUserInfo")
	userInfo := userinfo.NewStruct(i.UserInfo.UserID)
	exist, err := dao.Get(ctx, userInfo, i.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeUserNotFound, "User Is Not Exist", fmt.Errorf("Get Unexist User"), newerror.LevelInfo)
	}
	return userInfo.Info, nil
}
func (i *Info) UpdateUserInfo(ctx context.Context, userConfig commonmodel.UserConfig) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("info:UpdateUserInfo")

	if day, month, year := i.UserInfo.BirthdayDay, i.UserInfo.BirthdayMonth, i.UserInfo.BirthdayYear; (day != 0 && month != 0 && year != 0) && (day > 31 || day < 0 || month < 0 || month > 12 || year > 1900) {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Invalid Parameter", fmt.Errorf("Unexpect Birthday"), newerror.LevelInfo)
	}
	if tool.CalculateLength(i.UserInfo.UserName) > userConfig.MaxUsernameLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "User Name Too Long", fmt.Errorf("User Name Too Long"), newerror.LevelInfo)
	}
	if tool.CalculateLength(i.UserInfo.Introduction) > userConfig.MaxIntroduceLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Introduce Too Long", fmt.Errorf("Introduce Too Long"), newerror.LevelInfo)
	}

	userInfo := userinfo.NewStruct(i.UserInfo.UserID,
		userinfo.WithUserName(i.UserInfo.UserName),
		userinfo.WithBirthday(i.UserInfo.BirthdayYear, i.UserInfo.BirthdayMonth, i.UserInfo.BirthdayDay),
		userinfo.WithIntroduction(i.UserInfo.Introduction),
	)
	exist, err := dao.Update(ctx, userInfo, i.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		err = dao.Add(ctx, userInfo, i.dbContext)
		if err != nil {
			return err
		}
	}
	return nil
}
func (i *Info) Remark(ctx context.Context, userConfig commonmodel.UserConfig) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("remark:Remark")
	remarkInfo := nickname.NewStruct(i.RemarkInfo.UserID, i.RemarkInfo.GoalUserID, i.RemarkInfo.NickName)
	if i.RemarkInfo.NickName == "" {
		err = dao.Delete(ctx, remarkInfo, i.dbContext)
		if err != nil {
			return err
		}
		return nil
	}
	if tool.CalculateLength(i.RemarkInfo.NickName) > userConfig.MaxNickNameLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Remark Name Too Long", fmt.Errorf("Remark Name Too Long"), newerror.LevelInfo)
	}
	exist, err := dao.Update(ctx, remarkInfo, i.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		err = dao.Add(ctx, remarkInfo, i.dbContext)
		if err != nil {
			return err
		}
	}
	return nil

}
