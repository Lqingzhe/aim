package handler

import (
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexgroupservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HandlerConfig) CreatSession(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreatSession", newerror.LevelInfo)
		return
	}
	if req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreatSession", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.CreatSessionReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		GoalUserId: req.GoalUserID,
	}
	kitexResp, err := h.serviceClient.GroupClient.CreatSession(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreatSession", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"session_info": gin.H{
				"session_id": strconv.FormatInt(kitexResp.SessionId, 10),
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "CreatSession", newerror.LevelInfo)
}

func (h *HandlerConfig) DeleteSession(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		SessionID int64 `json:"session_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteSession", newerror.LevelInfo)
		return
	}
	if req.SessionID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteSession", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.DeleteSessionReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		SessionId:  req.SessionID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.DeleteSession(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteSession", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "DeleteSession", newerror.LevelInfo)
}

func (h *HandlerConfig) GetFriendLastVisitTime(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		SessionID  int64 `json:"session_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFriendLastVisitTime", newerror.LevelInfo)
		return
	}
	if req.SessionID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFriendLastVisitTime", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.GetFriendLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		SessionId:  req.SessionID,
		GoalUserId: req.GoalUserID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetFriendLastVisitTime(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFriendLastVisitTime", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"session_info": gin.H{
				"goal_user_id":    strconv.FormatInt(req.GoalUserID, 10),
				"last_visit_time": kitexResp.LastVisitTime,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetFriendLastVisitTime", newerror.LevelInfo)
}

func (h *HandlerConfig) ApplyForFriend(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ApplyForFriend", newerror.LevelInfo)
		return
	}
	if req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ApplyForFriend", newerror.LevelInfo)
		return
	}

	kitexReq := &kitexgroupservice.ApplyForFriendReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.ApplyForFriend(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ApplyForFriend", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "ApplyForFriend", newerror.LevelInfo)
}

func (h *HandlerConfig) GetFriendApplyList(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	kitexReq := &kitexgroupservice.GetFriendApplyListReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetFriendApplyList(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFriendApplyList", err2.LogLevel)
		return
	}
	applyUserIDList := make([]string, 0, len(kitexResp.ApplyUserIdList))
	for _, applyUserID := range kitexResp.ApplyUserIdList {
		applyUserIDList = append(applyUserIDList, strconv.FormatInt(applyUserID, 10))
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"session_info": gin.H{
				"apply_user_list": applyUserIDList,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetFriendApplyList", newerror.LevelInfo)
}

func (h *HandlerConfig) RefuseFriendApply(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserId int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseFriendApply", newerror.LevelInfo)
		return
	}
	if req.GoalUserId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseFriendApply", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.RefuseFriendApplyReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		GoalUserId: req.GoalUserId,
	}
	_, err := h.serviceClient.GroupClient.RefuseFriendApply(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseFriendApply", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "RefuseFriendApply", newerror.LevelInfo)
}
