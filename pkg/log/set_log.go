package newlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetGinLog(c *gin.Context, logger *zap.Logger, message string, logLevel zapcore.Level) {
	c.Set("log_message", message)
	c.Set("log_level", logLevel)
	c.Set("logger", logger)
}
