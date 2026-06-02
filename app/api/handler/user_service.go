package handler

import (
	"aim/app/api/service"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexgroupservice"
	"aim/kitex_gen/kitexuserservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//

func (h *HandlerConfig) Register(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Register", newerror.LevelInfo)
		return
	}
	if req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Register", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexuserservice.RegisterReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		Password: req.Password,
	}
	kitexResp, err := h.serviceClient.UserClient.Register(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Register", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"user_info": gin.H{
				"user_id": strconv.FormatInt(kitexResp.UserId, 10),
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, kitexResp.UserId, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "Register", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) Login(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		UserID   int64  `json:"user_id,string"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Login", newerror.LevelInfo)
		return
	}
	if req.Password == "" || req.UserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Login", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexuserservice.LoginReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     req.UserID,
		Password:   req.Password,
	}
	_, err := h.serviceClient.UserClient.Login(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, req.UserID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Login", err2.LogLevel)
		return
	}

	serviceStruct := service.NewToken(h.dbContext, h.tokenConfig)
	accessToken, refreshToken, err := serviceStruct.MakeTokens(ctx, req.UserID, c.GetHeader("X-Device-ID"))
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, req.UserID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Login", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"token_info": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, req.UserID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "Login", newerror.LevelInfo)
	return
}

func (h *HandlerConfig) LogoutAll(c *gin.Context) {
	userID := c.GetInt64("user_id")
	tokenStruct := service.NewToken(h.dbContext, h.tokenConfig)
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	err := tokenStruct.ReleaseAllToken(c.Request.Context(), userID)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Logout", err2.LogLevel)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		return
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "Logout", newerror.LevelInfo)
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "Logout Success",
	})
}
func (h *HandlerConfig) LogoutOne(c *gin.Context) {
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	tokenStruct := service.NewToken(h.dbContext, h.tokenConfig)
	err := tokenStruct.ReleaseOneTokenWithDeviceID(c.Request.Context(), userID, c.GetHeader("X-Device-ID"))
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Logout", err2.LogLevel)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		return
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "Logout", newerror.LevelInfo)
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "Logout Success",
	})
}

func (h *HandlerConfig) RefreshToken(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefreshToken", newerror.LevelInfo)
		return
	}
	if req.RefreshToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefreshToken", newerror.LevelInfo)
		return
	}
	tokenStruct := service.NewToken(h.dbContext, h.tokenConfig)
	userID, accessToken, refreshToken, err := tokenStruct.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, 0, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "RefreshToken", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"token_info": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "RefreshToken", newerror.LevelInfo)
	return
}

//

func (h *HandlerConfig) GetUserInfo(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	userID := c.GetInt64("user_id")
	kitexReq := kitexuserservice.GetUserInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		UserId: userID,
	}
	kitexResp, err := h.serviceClient.UserClient.GetUserInfo(ctx, &kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
	}
	type userInfo struct {
		UserID        int64  `json:"user_id,string"`
		UserName      string `json:"user_name"`
		Introduction  string `json:"introduction"`
		BirthdayYear  int64  `json:"birthday_year"`
		BirthdayMonth int64  `json:"birthday_month"`
		BirthdayDay   int64  `json:"birthday_day"`
	}

	type remarkInfo struct {
		GoalUserID int64  `json:"goal_user_name,string"`
		NickName   string `json:"nick_name"`
	}
	RemarkInfoList := make([]remarkInfo, len(kitexResp.RemarkInfo))
	for i, j := range kitexResp.RemarkInfo {
		RemarkInfoList[i] = remarkInfo{
			GoalUserID: j.GoalUserID,
			NickName:   j.NickName,
		}
	}
	var resp = struct {
		UserInfo    userInfo
		RemarkInfos []remarkInfo
	}{
		UserInfo: userInfo{
			UserID:        userID,
			UserName:      kitexResp.UserInfo.UserName,
			Introduction:  kitexResp.UserInfo.Introduction,
			BirthdayYear:  kitexResp.UserInfo.BirthdayYear,
			BirthdayMonth: kitexResp.UserInfo.BirthdayMonth,
			BirthdayDay:   kitexResp.UserInfo.BirthdayDay,
		},
		RemarkInfos: RemarkInfoList,
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"user_info": resp,
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetUserInfo", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) GetOtherUserInfo(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserID int64 `json:"goal_user_id,string"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetOtherUserInfo", newerror.LevelInfo)
		return
	}
	if req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetOtherUserInfo", newerror.LevelInfo)
		return
	}
	kitexReq := kitexuserservice.GetOtherUserInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		GoalUserId: req.GoalUserID,
	}
	kitexResp, err := h.serviceClient.UserClient.GetOtherUserInfo(ctx, &kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetOtherUserInfo", newerror.LevelInfo)
	}
	h.hub.Mu.RLock()
	deviceMap := h.hub.Client[req.GoalUserID]
	IsConnect := false
	for _, client := range deviceMap {
		if client.IsConnected {
			IsConnect = true
			break
		}
	}
	h.hub.Mu.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"user_info": gin.H{
				"user_name":      kitexResp.UserInfo.UserName,
				"introduction":   kitexResp.UserInfo.Introduction,
				"birthday_year":  kitexResp.UserInfo.BirthdayYear,
				"birthday_month": kitexResp.UserInfo.BirthdayMonth,
				"birthday_day":   kitexResp.UserInfo.BirthdayDay,
				"is_connect":     IsConnect,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetOtherUserInfo", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) UpdateUserInfo(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		UserName      string `json:"user_name"`
		Introduction  string `json:"introduction"`
		BirthdayYear  int64  `json:"birthday_year"`
		BirthdayMonth int64  `json:"birthday_month"`
		BirthdayDay   int64  `json:"birthday_day"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateUserInfo", newerror.LevelInfo)
		return
	}
	if req.UserName == "" && req.Introduction == "" && (req.BirthdayYear&req.BirthdayMonth&req.BirthdayDay == 0) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateUserInfo", newerror.LevelInfo)
		return
	}
	kitexReq := kitexuserservice.UpdateUserInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		UserInfo: &kitexcommonmodel.UserInfo{
			UserID:        userID,
			UserName:      req.UserName,
			Introduction:  req.Introduction,
			BirthdayYear:  req.BirthdayYear,
			BirthdayMonth: req.BirthdayMonth,
			BirthdayDay:   req.BirthdayDay,
		},
	}
	_, err := h.serviceClient.UserClient.UpdateUserInfo(ctx, &kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "UpdateUserInfo", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "UpdateUserInfo", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) Remark(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GoalUserID int64  `json:"goal_user_id,string"`
		NickName   string `json:"nick_name"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Remark", newerror.LevelInfo)
		return
	}
	if req.GoalUserID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Remark", newerror.LevelInfo)
		return
	}
	kitexReq := kitexuserservice.RemarkReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{
			Trace: c.GetString("trace"),
		},
		RemarkInfo: &kitexcommonmodel.RemarkInfo{
			UserID:     userID,
			GoalUserID: req.GoalUserID,
			NickName:   req.NickName,
		},
	}
	_, err := h.serviceClient.UserClient.Remark(ctx, &kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "Remark", err2.LogLevel)
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
	newlog.SetGinLog(c, logger, "Remark", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) GetGroupAndSessionID(c *gin.Context) {
	var finalErr error
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	kitexReq := &kitexgroupservice.GetGroupAndSessionIDReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.GroupClient.GetGroupAndSessionID(ctx, kitexReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		err2 := newerror.TranslateError(finalErr)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{"code": err2.StatusCode, "message": err2.HttpMessage})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetGroupAndSessionID", err2.LogLevel)
		return
	}
	sessionIDList := make([]string, 0, len(kitexResp.SessionIdList))
	userOfSessionIDList := make([]string, 0, len(kitexResp.SessionIdList))
	for i := range kitexResp.SessionIdList {
		sessionIDList = append(sessionIDList, strconv.FormatInt(kitexResp.SessionIdList[i], 10))
		userOfSessionIDList = append(userOfSessionIDList, strconv.FormatInt(kitexResp.UserOfSessionIdList[i], 10))
	}
	groupIDList := make([]string, 0, len(kitexResp.GroupIdList))
	for i := range kitexResp.GroupIdList {
		groupIDList = append(groupIDList, strconv.FormatInt(kitexResp.GroupIdList[i], 10))
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"session_info": gin.H{
				"session_id_list":         sessionIDList,
				"user_of_session_id_list": userOfSessionIDList,
			},
			"group_info": gin.H{
				"group_id_list": groupIDList,
			},
		},
	})
	if finalErr != nil {
		err2 := newerror.TranslateError(finalErr)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetGroupAndSessionID", newerror.LevelInfo)
}
