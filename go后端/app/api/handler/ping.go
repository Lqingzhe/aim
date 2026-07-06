package handler

import (
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HandlerConfig) Ping(c *gin.Context) {
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	ip := c.ClientIP()
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, 0, ip, c.FullPath())
	newlog.SetGinLog(c, logger, "Ping", newerror.LevelInfo)
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "pong",
	})
}
