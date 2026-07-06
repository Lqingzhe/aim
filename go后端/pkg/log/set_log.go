package newlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetGinLog[t int8 | zapcore.Level](c *gin.Context, logger *zap.Logger, message string, logLevel t) {
	c.Set("log_message", message)
	c.Set("log_level", logLevel)
	c.Set("logger", logger)
}
