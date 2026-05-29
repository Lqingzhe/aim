package middleware

import (
	"aim/app/api/model"
	"aim/app/api/service"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Limiter(dbContext *model.DBContext, limiterConfig commonmodel.LimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		a, _ := c.Get("logger")
		logger := a.(*zap.Logger)
		limiterService := service.NewLimiter(dbContext, limiterConfig)
		whetherLimite, err := limiterService.Exceed(c.Request.Context(), c.GetInt64("user_id"), c.GetHeader("X-Device-ID"))
		if err != nil {
			err2 := newerror.TranslateError(err)
			c.JSON(err2.HttpCode, gin.H{
				"code":    err2.StatusCode,
				"message": err2.HttpMessage,
			})
			logger = newlog.AddError(logger, err, err2.StatusCode)
			logger = newlog.AddGateWayInfo(logger, err2.HttpCode, c.GetInt64("user_id"), c.ClientIP(), c.GetString("operation"))
			newlog.SetGinLog(c, logger, "Limiter", err2.LogLevel)
			c.Abort()
			return
		}
		if whetherLimite {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    newerror.CodeRateLimitExceeded,
				"message": "Too Many Requests",
			})
			logger = newlog.AddGateWayInfo(newlog.AddError(logger, fmt.Errorf("Frequent Access"), newerror.CodeRateLimitExceeded), http.StatusTooManyRequests, c.GetInt64("user_id"), c.ClientIP(), c.GetString("operation"))
			newlog.SetGinLog(c, logger, "Limiter", newerror.LevelInfo)
			c.Abort()
			return
		}
		c.Next()
	}
}
