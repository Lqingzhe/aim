package handler

import (
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HandlerConfig) SetMute(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupId         int64  `json:"group_id,string"`
		GoalUserId      int64  `json:"goal_user_id,string"`
		MuteTimeSeconds int64  `json:"mute_time_seconds"`
		MuteReason      string `json:"mute_reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetMute", newerror.LevelInfo)
		return
	}
	if req.MuteTimeSeconds == 0 || req.GroupId == 0 || req.GoalUserId == 0 || req.MuteReason == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetMute", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.SetMuteReq{
		CommonInfo:     &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:        req.GroupId,
		UserId:         userID,
		GoalUserId:     req.GoalUserId,
		MuteReason:     req.MuteReason,
		MuteTimeSecond: req.MuteTimeSeconds,
	}
	_, err := h.serviceClient.GroupClient.SetMute(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetMute", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "SetMute", newerror.LevelInfo)
}

func (h *HandlerConfig) ReleaseMute(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserId int64 `json:"goal_user_id,string"`
		GroupId    int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ReleaseMute", newerror.LevelInfo)
		return
	}
	if req.GroupId == 0 || req.GoalUserId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ReleaseMute", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.ReleaseMuteReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GoalUserId: req.GoalUserId,
		GroupId:    req.GroupId,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.ReleaseMute(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ReleaseMute", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "ReleaseMute", newerror.LevelInfo)
}
