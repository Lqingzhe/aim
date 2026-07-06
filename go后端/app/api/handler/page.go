package handler

import (
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginPage 登录/注册页面
func (h *HandlerConfig) LoginPage(c *gin.Context) {
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "IM 即时通讯",
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, -1, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "LoginPage", newerror.LevelInfo)
}

// ChatPage 聊天主页面（需要登录）
func (h *HandlerConfig) ChatPage(c *gin.Context) {
	userID := c.GetInt64("user_id")
	deviceID := c.GetHeader("X-Device-ID")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "ChatPage", newerror.LevelInfo)
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"title":    "IM 聊天",
		"userID":   userID,
		"deviceID": deviceID,
		"wsURL":    "/ws",
	})
}

// IndexPage 根路径重定向
func (h *HandlerConfig) IndexPage(c *gin.Context) {
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	c.Redirect(http.StatusFound, "/login")
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, -1, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "IndexPage", newerror.LevelInfo)
}
