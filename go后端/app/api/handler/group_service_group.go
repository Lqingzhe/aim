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

func (h *HandlerConfig) GetGroupInfo(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfo", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfo", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.GetGroupInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		GroupId: req.GroupID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetGroupInfo(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfo", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"group_id":   strconv.FormatInt(req.GroupID, 10),
				"group_name": kitexResp.GroupName,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetGroupInfo", newerror.LevelInfo)
}

func (h *HandlerConfig) ChangeGroupInfo(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID   int64  `json:"group_id,string"`
		GroupName string `json:"group_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ChangeGroupInfo", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GroupName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ChangeGroupInfo", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.ChangeGroupInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GroupName:  req.GroupName,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.ChangeGroupInfo(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "ChangeGroupInfo", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "ChangeGroupInfo", newerror.LevelInfo)
}

func (h *HandlerConfig) SearchGroup(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupName string `json:"group_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SearchGroup", newerror.LevelInfo)
		return
	}
	if req.GroupName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SearchGroup", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.SearchGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupName:  req.GroupName,
	}
	kitexResp, err := h.serviceClient.GroupClient.SearchGroup(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SearchGroup", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"group_id_list": kitexResp.GroupIdList,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SearchGroup", newerror.LevelInfo)
}

func (h *HandlerConfig) CreateGroup(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupName string `json:"group_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreateGroup", newerror.LevelInfo)
		return
	}
	if req.GroupName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreateGroup", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.CreateGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		GroupName:  req.GroupName,
	}
	kitexResp, err := h.serviceClient.GroupClient.CreateGroup(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "CreateGroup", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"group_id": strconv.FormatInt(kitexResp.GroupId, 10),
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "CreateGroup", newerror.LevelInfo)
}

func (h *HandlerConfig) DeleteGroup(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteGroup", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteGroup", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.DeleteGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.DeleteGroup(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "DeleteGroup", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "DeleteGroup", newerror.LevelInfo)
}

func (h *HandlerConfig) LeaveGroup(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "LeaveGroup", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "LeaveGroup", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.LeaveGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.LeaveGroup(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "LeaveGroup", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "LeaveGroup", newerror.LevelInfo)
}

func (h *HandlerConfig) SetGroupApply(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetGroupApply", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetGroupApply", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.SetGroupApplyReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.SetGroupApply(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetGroupApply", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": newerror.CodeSuccess, "message": "success"})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SetGroupApply", newerror.LevelInfo)
}

func (h *HandlerConfig) GetGroupApplyList(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupApplyList", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupApplyList", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.GetGroupApplyListReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		GroupId:    req.GroupID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetGroupApplyList(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupApplyList", err2.LogLevel)
		return
	}
	idList := make([]string, 0, len(kitexResp.ApplyUserIdList))
	for _, id := range kitexResp.ApplyUserIdList {
		idList = append(idList, strconv.FormatInt(id, 10))
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"group_id_list": idList,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetGroupApplyList", newerror.LevelInfo)
}

func (h *HandlerConfig) GetLastVisitTime(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetLastVisitTime", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetLastVisitTime", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.GetLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetLastVisitTime(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetLastVisitTime", err2.LogLevel)
		return
	}
	lastVisitTimeMap := make(map[string]int64)
	for i := range kitexResp.LastVisitTimeList {
		lastVisitTimeMap[strconv.FormatInt(kitexResp.UserIdList[i], 10)] = kitexResp.LastVisitTimeList[i]
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"last_visit_time": lastVisitTimeMap,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetLastVisitTime", newerror.LevelInfo)
}

func (h *HandlerConfig) AgreeGroupApply(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "AgreeGroupApply", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "AgreeGroupApply", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.AgreeGroupApplyReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.AgreeGroupApply(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "AgreeGroupApply", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "AgreeGroupApply", newerror.LevelInfo)
}
func (h *HandlerConfig) RefuseGroupApply(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseGroupApply", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseGroupApply", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.RefuseGroupApplyReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.RefuseGroupApply(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefuseGroupApply", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "RefuseGroupApply", newerror.LevelInfo)
}

func (h *HandlerConfig) TransformGroupOwner(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "TransformGroupOwner", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "TransformGroupOwner", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.TransformGroupOwnerReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.TransformGroupOwner(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "TransformGroupOwner", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "TransformGroupOwner", newerror.LevelInfo)
}

func (h *HandlerConfig) KickOutGroup(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "KickOutGroup", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "KickOutGroup", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.KickOutGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.KickOutGroup(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "KickOutGroup", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "KickOutGroup", newerror.LevelInfo)
}

func (h *HandlerConfig) SetManager(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetManager", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetManager", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.SetManagerReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.SetManager(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetManager", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "SetManager", newerror.LevelInfo)
}

func (h *HandlerConfig) RevokeManager(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID    int64 `json:"group_id,string"`
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RevokeManager", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RevokeManager", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.RevokeManagerReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
		GroupId:    req.GroupID,
		GoalUserId: req.GoalUserID,
	}
	_, err := h.serviceClient.GroupClient.RevokeManager(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RevokeManager", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "RevokeManager", newerror.LevelInfo)
}

func (h *HandlerConfig) GetGroupInfoWithUser(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfoWithUser", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfoWithUser", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.GetGroupInfoWithUserReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetGroupInfoWithUser(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupInfoWithUser", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"group_info": gin.H{
				"group_id":          strconv.FormatInt(kitexResp.GroupId, 10),
				"group_remark_name": kitexResp.GroupRemarkName,
				"group_role":        kitexResp.Role,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetGroupInfoWithUser", newerror.LevelInfo)
}

func (h *HandlerConfig) UpdateGroupInfoWithUser(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID         int64  `json:"group_id,string"`
		GroupRemarkName string `json:"group_remark_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateGroupInfoWithUser", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.GroupRemarkName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateGroupInfoWithUser", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.UpdateGroupInfoWithUserReq{
		CommonInfo:      &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:         req.GroupID,
		GroupRemarkName: req.GroupRemarkName,
		UserId:          userID,
	}
	_, err := h.serviceClient.GroupClient.UpdateGroupInfoWithUser(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateGroupInfoWithUser", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "UpdateGroupInfoWithUser", newerror.LevelInfo)
}
func (h *HandlerConfig) SetLastVisitTime(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id,string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetLastVisitTime", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetLastVisitTime", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexgroupservice.SetLastVisitTimeReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	_, err := h.serviceClient.GroupClient.SetLastVisitTime(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SetLastVisitTime", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "SetLastVisitTime", newerror.LevelInfo)
}
