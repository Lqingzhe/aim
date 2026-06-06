package handler

import "C"
import (
	"aim/kitex_gen/kitexaiservice"
	"aim/kitex_gen/kitexcommonmodel"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HandlerConfig) DeleteChatContext(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	kitexReq := &kitexaiservice.DeleteChatContextReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
	}
	_, err := h.serviceClient.AiClient.DeleteChatContext(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteChatContext", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "DeleteChatContext", newerror.LevelInfo)
}
func (h *HandlerConfig) GetAiConfig(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	kitexReq := &kitexaiservice.GetAiConfigReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.AiClient.GetAiConfig(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetAiConfig", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"ai_info": gin.H{
				"model_name": kitexResp.ModelName,
				"base_url":   kitexResp.BaseUrl,
				"api_key":    kitexResp.ApiKey,
				"role":       kitexResp.Role,
				"prompt":     kitexResp.Prompt,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetAiConfig", newerror.LevelInfo)
}
func (h *HandlerConfig) UpdateAiConfig(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		ModelName string `json:"model_name"`
		BaseUrl   string `json:"base_url"`
		ApiKey    string `json:"api_key"`
		Role      string `json:"role"`
		Prompt    string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateAiConfig", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexaiservice.UpdateAiConfigReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		ModelName:  req.ModelName,
		BaseUrl:    req.BaseUrl,
		ApiKey:     req.ApiKey,
		Role:       req.Role,
		Prompt:     req.Prompt,
	}
	_, err := h.serviceClient.AiClient.UpdateAiConfig(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateAiConfig", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "UpdateAiConfig", newerror.LevelInfo)
}
func (h *HandlerConfig) DeleteAiConfig(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	kitexReq := &kitexaiservice.DeleteAiConfigReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
	}
	_, err := h.serviceClient.AiClient.DeleteAiConfig(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteAiConfig", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "DeleteAiConfig", newerror.LevelInfo)
}
