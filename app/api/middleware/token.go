package middleware

import (
	"aim/app/api/model"
	"aim/app/api/service"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AnalyseToken(tokenConfig commonmodel.TokenConfig, dbContext *model.DBContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		a, _ := c.Get("logger")
		logger := a.(*zap.Logger)
		tokenStruct := service.NewToken(dbContext, tokenConfig)
		rawAccessToken := c.GetHeader("Authorization")
		if !strings.HasPrefix(rawAccessToken, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    newerror.CodeMissingParam,
				"message": "Lack Access Token",
			})
			logger = newlog.AddGateWayInfo(newlog.AddError(logger, fmt.Errorf("Lack Access Token"), newerror.CodeMissingParam), http.StatusUnauthorized, 0, c.ClientIP(), c.GetString("operation"))
			newlog.SetGinLog(c, logger, "Analyse Token", newerror.LevelInfo)
			c.Abort()
			return
		}
		accessToken := strings.TrimPrefix(rawAccessToken, "Bearer ")
		userID, deviceID, err := tokenStruct.AnalysisAccessToken(accessToken)
		if err != nil {
			err2 := newerror.TranslateError(err)
			c.JSON(err2.HttpCode, gin.H{
				"code":    err2.StatusCode,
				"message": err2.HttpMessage,
			})
			logger = newlog.AddError(logger, err, err2.StatusCode)
			logger = newlog.AddGateWayInfo(logger, err2.HttpCode, 0, c.ClientIP(), c.GetString("operation"))
			newlog.SetGinLog(c, logger, "Analyse Token", err2.LogLevel)
			c.Abort()
			return
		}
		if deviceID != c.GetHeader("X-Device-ID") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    newerror.CodeAccessTokenInvalid,
				"message": "Useless Token",
			})
			logger = newlog.AddError(logger, fmt.Errorf("Useless Token"), newerror.CodeAccessTokenInvalid)
			logger = newlog.AddGateWayInfo(logger, http.StatusUnauthorized, userID, c.ClientIP(), c.GetString("operation"))
			newlog.SetGinLog(c, logger, "Analyse Token", newerror.LevelInfo)
			c.Abort()
			return
		}
		c.Set("user_id", userID)
	}
}
