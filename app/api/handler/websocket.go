package handler

import (
	"aim/app/api/service"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HandlerConfig) ConnectWebsocket(c *gin.Context) {
	beginTime := time.Now()
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	accessToken := c.Query("token")
	rawCtx, _ := c.Get("ctx")
	traceID := c.GetString("trace")
	ip := c.ClientIP()
	fullPath := c.FullPath()
	ctx := rawCtx.(context.Context)
	tokenStruct := service.NewToken(h.dbContext, h.tokenConfig)
	userID, deviceID, err := tokenStruct.AnalysisAccessToken(accessToken)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, ip, fullPath)
		logger = newlog.AddLatencyAndTime(logger, beginTime)
		newlog.Log(logger, newerror.LevelInfo, "ConnectWebsocket")
		return
	}
	webStruct := service.NewWebSocket(h.hub)
	conn, err := h.websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger = newlog.AddError(logger, err, newerror.CodeInternalError)
		logger = newlog.AddGateWayInfo(logger, -1, userID, ip, fullPath)
		logger = newlog.AddLatencyAndTime(logger, beginTime)
		newlog.Log(logger, newerror.LevelInfo, "ConnectWebsocket")
		return
	}

	webStruct.ConnectWebsocket(conn, userID, deviceID)
	logger, message, logLevel := h.GetOfflineMessages(ctx, logger, userID, deviceID, traceID, ip, fullPath)
	logger = newlog.AddGateWayInfo(logger, -1, userID, ip, fullPath)
	logger = newlog.AddLatencyAndTime(logger, beginTime)
	newlog.Log(logger, logLevel, message)
}

func (h *HandlerConfig) DisconnectWebsocketAllDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		webStruct := service.NewWebSocket(h.hub)
		webStruct.Unregister(userID, "")
	}
}
func (h *HandlerConfig) DisconnectWebsocketOneDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		deviceID := c.GetHeader("X-Device-ID")
		webStruct := service.NewWebSocket(h.hub)
		webStruct.Unregister(userID, deviceID)
	}
}
