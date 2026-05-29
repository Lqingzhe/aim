package service

import (
	"aim/app/api/dao"
	"aim/app/api/dao/refreshtoken"
	"aim/app/api/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	dbContext   *model.DBContext
	tokenConfig commonmodel.TokenConfig
}

func NewToken(dbContext *model.DBContext, tokenConfig commonmodel.TokenConfig) *Token {
	return &Token{
		dbContext:   dbContext,
		tokenConfig: tokenConfig,
	}
}
func (t *Token) MakeTokens(ctx context.Context, userID int64, deviceID string) (accessToken string, newRefreshToken string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("token:MakeTokens")
	//生成RefreshToken
	salt, err := tool.AddSaltByByteLength(t.tokenConfig.SaltByteLen)
	if err != nil {
		return "", "", err
	}
	refreshToken := fmt.Sprintf("%x%s%s", userID, deviceID, salt)
	TokenStruct := refreshtoken.NewStruct(refreshToken, userID, deviceID, t.tokenConfig.RefreshTokenExpireTime)
	err = dao.Add(ctx, TokenStruct, t.dbContext)
	if err != nil {
		return "", "", err
	}
	//生成AccessToken
	claim := jwt.MapClaims{
		"user_id":   userID,
		"device_id": deviceID,
		"exp":       time.Now().Add(t.tokenConfig.AccessTokenExpireTime).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(t.tokenConfig.TokenPassword))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func (t *Token) RefreshToken(ctx context.Context, oldRefreshToken string) (userID int64, accessToken string, newRefreshToken string, err error) {
	defer func() {
		err = newerror.TranslateError(err).AddErrorTrace("token:RefreshToken")
	}()
	tokenStruct := refreshtoken.NewStruct(oldRefreshToken, 0, "", 0)
	exist, err := dao.Get(ctx, tokenStruct, t.dbContext)
	if err != nil {
		return 0, "", "", err
	}
	if !exist {
		return 0, "", "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeRefreshTokenInvalid, "Useless Refresh Tokens", fmt.Errorf("%s", "Refresh Token Is Not Exist"), newerror.LevelInfo)
	}
	salt, err := tool.AddSaltByByteLength(t.tokenConfig.SaltByteLen)
	if err != nil {
		return 0, "", "", err
	}
	err = dao.Delete(ctx, tokenStruct, t.dbContext)
	if err != nil {
		return 0, "", "", err
	}
	tokenStruct.TokenInfo = *tokenStruct.Info
	tokenStruct.RefreshTokenID = fmt.Sprintf("%x%s%s", tokenStruct.UserID, tokenStruct.DeviceID, salt)
	tokenStruct.ExpireTime = t.tokenConfig.RefreshTokenExpireTime
	err = dao.Add(ctx, tokenStruct, t.dbContext)
	if err != nil {
		return 0, "", "", err
	}
	claim := jwt.MapClaims{
		"user_id":   tokenStruct.UserID,
		"device_id": tokenStruct.DeviceID,
		"exp":       time.Now().Add(t.tokenConfig.AccessTokenExpireTime).Unix(),
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(t.tokenConfig.TokenPassword))
	if err != nil {
		return 0, "", "", err
	}
	return tokenStruct.UserID, accessToken, tokenStruct.RefreshTokenID, nil
}
func (t *Token) ReleaseAllToken(ctx context.Context, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("token:ReleaseAllToken")
	tokenStruct := refreshtoken.NewStruct(fmt.Sprintf("%x", userID), 0, "", 0)
	err = tokenStruct.DeleteInfo(ctx, t.dbContext)
	return err
}
func (t *Token) ReleaseOneTokenWithDeviceID(ctx context.Context, userID int64, deviceID string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("token:ReleaseOneTokenWithDeviceID")
	tokenStruct := refreshtoken.NewStruct(fmt.Sprintf("%x%s", userID, deviceID), 0, "", 0)
	err = tokenStruct.DeleteInfo(ctx, t.dbContext)
	return err
}
func (t *Token) AnalysisAccessToken(accessToken string) (userID int64, deviceID string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("token:AnalysisAccessToken")
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.tokenConfig.TokenPassword), nil
	})
	if err != nil {
		var errMsg string
		// 精准区分JWT库定义的标准错误类型
		switch {
		// 1. Token已过期（最常见，优先提示）
		case errors.Is(err, jwt.ErrTokenExpired):
			errMsg = "Token Is Expired"
			return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenExpired, "Invalid Token", fmt.Errorf("%s", errMsg), newerror.LevelInfo)
		// 2. Token签名错误（被篡改/秘钥错误）
		case errors.Is(err, jwt.ErrSignatureInvalid):
			errMsg = "Signature Error"
		// 3. Token尚未生效（nbf字段设置了生效时间）
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			errMsg = "Not Valid Yet"
		// 4. Token格式错误（不是3段式/缺少.分隔符）
		case strings.Contains(err.Error(), "malformed"):
			errMsg = "Malformed Format"
		// 5. 算法不匹配（你自定义的算法错误）
		case strings.Contains(err.Error(), "unexpected signing method"):
			errMsg = "Unsupported Algorithm"
		// 6. 其他未知错误（保留原标识，便于定位）
		default:
			errMsg = err.Error()
		}
		return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenInvalid, "Invalid Token", fmt.Errorf("%s", errMsg), newerror.LevelWarn)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenInvalid, "Invalid Token", fmt.Errorf("%s", `Type Assertion To "jwt.MapClaims" Error`), newerror.LevelWarn)
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenInvalid, "Invalid Token", fmt.Errorf("%s", `Type Assertion To "expiretime as float64" Error`), newerror.LevelWarn)
	}
	exp := int64(expFloat)
	if time.Now().Unix() > exp {
		return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenInvalid, "Invalid Token", fmt.Errorf("%s", "Token Is Expired"), newerror.LevelInfo)
	}
	if !parsedToken.Valid || claims["user_id"] == nil || claims["device_id"] == nil {
		return 0, "", newerror.MakeError(http.StatusUnauthorized, newerror.CodeAccessTokenInvalid, "Invalid Token", fmt.Errorf("%s", "Lack Useful Info"), newerror.LevelWarn)
	}
	userID = int64(claims["user_id"].(float64))
	deviceID = claims["device_id"].(string)
	return userID, deviceID, nil
}
