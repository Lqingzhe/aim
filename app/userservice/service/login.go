package service

import (
	"aim/app/userservice/dao"
	"aim/app/userservice/dao/userinfo"
	"aim/app/userservice/dao/userlogin"
	"aim/app/userservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"crypto/hmac"
	"fmt"
	"net/http"

	"github.com/bwmarrin/snowflake"
)

type LoginInfo struct {
	dbContext *model.DBContext

	LoginInfo *model.UserLoginInfo
	UserInfo  *model.UserInfo
}

func NewLoginInfo(dbContext *model.DBContext, userLogin *model.UserLoginInfo, userInfo *model.UserInfo) *LoginInfo {
	return &LoginInfo{
		dbContext: dbContext,
		LoginInfo: userLogin,
		UserInfo:  userInfo,
	}
}

func (l *LoginInfo) Register(ctx context.Context, userConfig commonmodel.UserConfig, snowFlack *snowflake.Node) (userID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("info:AddUserInfo")
	passwordLength := tool.CalculateLength(l.LoginInfo.Password)
	if passwordLength < userConfig.MinPasswordLength {
		return 0, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Password Too Short", fmt.Errorf("Password Input Too Shout"), newerror.LevelInfo)
	}
	if passwordLength > userConfig.MaxPasswordLength {
		return 0, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Password Too Long", fmt.Errorf("Password Input Too Long"), newerror.LevelInfo)
	}
	salt, err := tool.AddSaltByByteLength(userConfig.SaltByteLen)
	if err != nil {
		return 0, err
	}
	dbContext := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(l.dbContext.Mysql),
	}
	userID = snowFlack.Generate().Int64()
	userLogin := userlogin.NewStruct(userID, tool.Encrypt(salt+l.LoginInfo.Password), salt)
	err = dao.Add(ctx, userLogin, dbContext)
	if err != nil {
		dbContext.Mysql.Client.Rollback()
		return 0, err
	}
	userInfo := userinfo.NewStruct(userID)
	err = dao.Add(ctx, userInfo, dbContext)
	if err != nil {
		dbContext.Mysql.Client.Rollback()
		return 0, err
	}
	result := dbContext.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return 0, err2
	}
	return userID, nil
}

func (l *LoginInfo) Login(ctx context.Context) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("info:Login")
	userLogin := userlogin.NewStruct(l.LoginInfo.UserID, "", "")
	exist, err := dao.Get(ctx, userLogin, l.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeUserNotFound, "User Not Found", fmt.Errorf("User Not Found"), newerror.LevelInfo)
	}
	if !hmac.Equal([]byte(userLogin.Info.Password), []byte(tool.Encrypt(userLogin.Info.Salt+l.LoginInfo.Password))) {
		return newerror.MakeError(http.StatusUnauthorized, newerror.CodePasswordWrong, "Password Error", fmt.Errorf("Password Error"), newerror.LevelInfo)
	}
	return nil
}
