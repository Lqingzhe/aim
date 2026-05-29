package middleware

import (
	newerror "aim/pkg/error"
	"aim/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(rawLogger *zap.Logger, equipID int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()
		trace := uuid.NewString()
		logger := newlog.AddTraceAndEquipID(rawLogger, trace, equipID)
		c.Set("logger", logger)
		c.Set("trace", trace)

		c.Next()

		message := c.GetString("log_message")
		loglevel := c.GetInt("log_level")
		RawLogger, exist := c.Get("logger")
		if !exist {
			newlog.Log(rawLogger, newerror.LevelFatal, "Middleware Can't Get Logger")
		}
		logger, ok := RawLogger.(*zap.Logger)
		if !ok {
			newlog.Log(rawLogger, newerror.LevelFatal, "Middleware Can't Get Logger")
		}
		newlog.AddLatencyAndTime(logger, beginTime)
		newlog.Log(logger, zapcore.Level(loglevel), message)
	}
}
