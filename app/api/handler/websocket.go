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
	userID := c.GetInt64("user_id")
	deviceID := c.GetHeader("X-Device-ID")
	rawCtx, _ := c.Get("context")
	traceID := c.GetString("trace_id")
	ip := c.ClientIP()
	fullPath := c.FullPath()
	ctx := rawCtx.(context.Context)
	webStruct := service.NewWebSocket(h.hub)
	conn, err := h.websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger = newlog.AddError(logger, err, newerror.CodeInternalError)
		logger = newlog.AddGateWayInfo(logger, -1, userID, ip, fullPath)
		newlog.AddLatencyAndTime(logger, beginTime)
		newlog.Log(logger, newerror.LevelInfo, "ConnectWebsocket")
		return
	}

	webStruct.ConnectWebsocket(conn, userID, deviceID)
	logger, message, logLevel := h.GetOfflineMessages(ctx, logger, userID, deviceID, traceID, ip, fullPath)
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
