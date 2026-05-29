package handler

import (
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (h *HandlerConfig) SendMessage(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID        int64  `json:"group_id"`
		MessageContent string `json:"message_content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendMessage", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendMessage", newerror.LevelInfo)
		return
	}
	kitexReq := kitexmessageservice.SendMessageReq{
		CommonInfo:     &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:        req.GroupID,
		UserId:         userID,
		MessageContent: req.MessageContent,
		IsAi:           false,
	}
	kitexResp, err := h.serviceClient.MessageClient.SendMessage(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{"code": err2.StatusCode, "message": err2.HttpMessage})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendMessage", err2.LogLevel)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_info": gin.H{
				"message_id": kitexResp.MessageId,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SendMessage", newerror.LevelInfo)
}
func (h *HandlerConfig) SendFile(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)

	GroupIDString := c.PostForm("group_id")
	GroupID, err := strconv.ParseInt(GroupIDString, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	FileName := c.PostForm("file_name")
	if GroupID == 0 || FileName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
			return
		}
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    newerror.CodeUnsupportedMedia,
				"message": "content type not supported",
			})
			logger = newlog.AddError(logger, err, newerror.CodeUnsupportedMedia)
			logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	defer file.Close()
	contentType := header.Header.Get("Content-Type")

	maxSize := int64(10 << 20)
	limitedReader := io.LimitReader(file, maxSize)
	dataStream, err := io.ReadAll(limitedReader)
	if err != nil {
		if errors.Is(err, bytes.ErrTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    newerror.CodeRequestBodyTooBig,
				"message": "Request Entity Too Large",
			})
			logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
			logger = newlog.AddGateWayInfo(logger, http.StatusRequestEntityTooLarge, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
			return
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	if int64(len(dataStream)) > maxSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeRequestBodyTooBig,
			"message": "File too large",
		})
		logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	kitexReq := kitexmessageservice.SendFileReq{
		CommonInfo:  &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:     GroupID,
		UserId:      userID,
		FileName:    FileName,
		ContentType: contentType,
		DataStream:  dataStream,
	}
	kitexResp, err := h.serviceClient.MessageClient.SendFile(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_info": gin.H{
				"message_id": kitexResp.MessageId,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SendFile", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) SendVoice(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	GroupIDString := c.PostForm("group_id")
	GroupID, err := strconv.ParseInt(GroupIDString, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	voiceTimeString := c.PostForm("voice_time")
	voiceTime, err := strconv.ParseInt(voiceTimeString, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	if GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
			return
		}
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    newerror.CodeUnsupportedMedia,
				"message": "content type not supported",
			})
			logger = newlog.AddError(logger, err, newerror.CodeUnsupportedMedia)
			logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	defer file.Close()
	contentType := header.Header.Get("Content-Type")

	maxSize := int64(10 << 20)
	limitedReader := io.LimitReader(file, maxSize)
	dataStream, err := io.ReadAll(limitedReader)
	if err != nil {
		if errors.Is(err, bytes.ErrTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    newerror.CodeRequestBodyTooBig,
				"message": "Request Entity Too Large",
			})
			logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
			logger = newlog.AddGateWayInfo(logger, http.StatusRequestEntityTooLarge, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
			return
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	if int64(len(dataStream)) > maxSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeRequestBodyTooBig,
			"message": "File too large",
		})
		logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexmessageservice.SendVoiceReq{
		CommonInfo:      &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:         GroupID,
		UserId:          userID,
		ContentType:     contentType,
		VoiceTimeSecond: voiceTime,
		DataStream:      dataStream,
	}
	kitexResp, err := h.serviceClient.MessageClient.SendVoice(ctx, kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_info": gin.H{
				"message_id": kitexResp.MessageId,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SendVoice", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) SendPicture(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)

	GroupIDString := c.PostForm("group_id")
	GroupID, err := strconv.ParseInt(GroupIDString, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	if GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	file, header, err := c.Request.FormFile("picture")
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
			return
		}
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    newerror.CodeUnsupportedMedia,
				"message": "content type not supported",
			})
			logger = newlog.AddError(logger, err, newerror.CodeUnsupportedMedia)
			logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	defer file.Close()
	contentType := header.Header.Get("Content-Type")

	maxSize := int64(10 << 20)
	limitedReader := io.LimitReader(file, maxSize)
	dataStream, err := io.ReadAll(limitedReader)
	if err != nil {
		if errors.Is(err, bytes.ErrTooLarge) {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    newerror.CodeRequestBodyTooBig,
				"message": "Request Entity Too Large",
			})
			logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
			logger = newlog.AddGateWayInfo(logger, http.StatusRequestEntityTooLarge, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
			return
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger = newlog.AddError(logger, err, newerror.CodeConnectionInterrupted)
			logger = newlog.AddGateWayInfo(logger, -1, userID, c.ClientIP(), c.FullPath())
			newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidParam,
			"message": "Formate Error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	if int64(len(dataStream)) > maxSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeRequestBodyTooBig,
			"message": "File too large",
		})
		logger = newlog.AddError(logger, err, newerror.CodeRequestBodyTooBig)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexmessageservice.SendPictureReq{
		CommonInfo:  &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:     GroupID,
		UserId:      userID,
		ContentType: contentType,
		DataStream:  dataStream,
	}
	kitexResp, err := h.serviceClient.MessageClient.SendPicture(ctx, kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_info": gin.H{
				"message_id": kitexResp.MessageId,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SendPicture", newerror.LevelInfo)
	return
}
func (h *HandlerConfig) WithdrawMessage(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID   int64 `json:"group_id"`
		MessageID int64 `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "WithdrawMessage", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.MessageID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "WithdrawMessage", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexmessageservice.WithdrawMessageReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
		MessageId:  req.MessageID,
	}
	_, err := h.serviceClient.MessageClient.WithdrawMessage(ctx, kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "WithdrawMessage", newerror.LevelInfo)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "WithdrawMessage", newerror.LevelInfo)
}
func (h *HandlerConfig) GetMessageList(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID         int64 `json:"group_id"`
		StartTimeSecond int64 `json:"start_time_second"`
		EndTimeSecond   int64 `json:"end_time_second"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetMessageList", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetMessageList", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexmessageservice.GetMessageListReq{
		CommonInfo:      &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:         req.GroupID,
		StartTimeSecond: req.StartTimeSecond,
		EndTimeSecond:   req.EndTimeSecond,
	}
	kitexResp, err := h.serviceClient.MessageClient.GetMessageList(ctx, kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetMessageList", newerror.LevelInfo)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_info": gin.H{
				"message_list": kitexResp.MessageInfo,
			},
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetMessageList", newerror.LevelInfo)
}
func (h *HandlerConfig) GetNewMessage(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID int64 `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetNewMessage", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetNewMessage", newerror.LevelInfo)
		return
	}
	kitexReq := &kitexmessageservice.GetNewMessageReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.MessageClient.GetNewMessage(ctx, kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetNewMessage", newerror.LevelInfo)
		return
	}
	type MessageInfo struct {
		MessageId      int64  `json:"message_id"`
		MessageContent string `json:"message_content"`
		MessageType    string `json:"message_type"`
		SendTimeSecond int64  `json:"send_time_second"`
	}
	MessageList := make([]*MessageInfo, 0, len(kitexResp.MessageId))
	for i := range kitexResp.MessageId {
		MessageList = append(MessageList, &MessageInfo{
			MessageId:      kitexResp.MessageId[i],
			MessageContent: kitexResp.MessageContent[i],
			MessageType:    kitexResp.MessageType[i],
			SendTimeSecond: kitexResp.SendTimeSecond[i],
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
		"data": gin.H{
			"message_list": MessageList,
		},
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetNewMessage", newerror.LevelInfo)
}
func (h *HandlerConfig) GetFileContent(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID   int64 `json:"group_id"`
		MessageID int64 `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFileContent", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.MessageID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFileContent", newerror.LevelInfo)
		return
	}
	kitexReq := kitexmessageservice.GetFileContentReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:    req.GroupID,
		MessageId:  req.MessageID,
		UserId:     userID,
	}
	kitexResp, err := h.serviceClient.MessageClient.GetFileContent(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "GetFileContent", newerror.LevelInfo)
		return
	}
	c.Data(http.StatusOK, kitexResp.ContentType, kitexResp.DataStream)
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "GetFileContent", newerror.LevelInfo)
}
func (h *HandlerConfig) SendGroupNotice(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	userID := c.GetInt64("user_id")
	a, _ := c.Get("logger")
	logger := a.(*zap.Logger)
	var req struct {
		GroupID        int64  `json:"group_id"`
		MessageContent string `json:"message_content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeInvalidJSON,
			"message": "JSON unmarshal error",
		})
		logger = newlog.AddError(logger, err, newerror.CodeInvalidJSON)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendGroupNotice", newerror.LevelInfo)
		return
	}
	if req.GroupID == 0 || req.MessageContent == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    newerror.CodeMissingParam,
			"message": "Lack Necessary Param",
		})
		logger = newlog.AddError(logger, fmt.Errorf("Lack Necessary Param"), newerror.CodeMissingParam)
		logger = newlog.AddGateWayInfo(logger, http.StatusBadRequest, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendGroupNotice", newerror.LevelInfo)
		return
	}
	kitexReq := kitexmessageservice.SendGroupNoticeReq{
		CommonInfo:     &kitexcommonmodel.CommonInfo{Trace: c.GetString("trace")},
		GroupId:        req.GroupID,
		UserId:         userID,
		MessageContent: req.MessageContent,
	}
	_, err := h.serviceClient.MessageClient.SendGroupNotice(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		c.AbortWithStatusJSON(err2.HttpCode, gin.H{
			"code":    err2.StatusCode,
			"message": err2.HttpMessage,
		})
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, c.ClientIP(), c.FullPath())
		newlog.SetGinLog(c, logger, "SendGroupNotice", newerror.LevelInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    newerror.CodeSuccess,
		"message": "success",
	})
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, c.ClientIP(), c.FullPath())
	newlog.SetGinLog(c, logger, "SendGroupNotice", newerror.LevelInfo)
}

func (h *HandlerConfig) GetOfflineMessages(ctx context.Context, rawLogger *zap.Logger, userID int64, deviceID string, traceID string, ip string, fullPath string) (logger *zap.Logger, message string, logLevel zapcore.Level) {
	var Data []byte
	defer func() {
		h.hub.Mu.RLock()
		if h.hub.Client[userID] != nil && h.hub.Client[userID][deviceID] != nil && h.hub.Client[userID][deviceID].IsConnected {
			h.hub.Client[userID][deviceID].Send <- Data
		}
		h.hub.Mu.RUnlock()
	}()
	type Result struct {
		HttpCode int
		gin.H
	}
	logger = rawLogger
	kitexReq := kitexmessageservice.GetOfflineMessageListReq{
		CommonInfo:      &kitexcommonmodel.CommonInfo{Trace: traceID},
		UserAndDeviceId: strconv.FormatInt(userID, 10) + deviceID,
	}
	kitexResp, err := h.serviceClient.MessageClient.GetOfflineMessageList(ctx, &kitexReq)
	if err != nil {
		err2 := newerror.TranslateError(err)
		R := Result{
			HttpCode: err2.HttpCode,
			H: gin.H{
				"Code":    err2.StatusCode,
				"Message": err2.HttpMessage,
			},
		}
		Data, _ = sonic.Marshal(R)
		logger = newlog.AddError(logger, err, err2.StatusCode)
		logger = newlog.AddGateWayInfo(logger, err2.HttpCode, userID, ip, fullPath)
		return logger, "GetOfflineMessages", newerror.LevelInfo
	}
	if !kitexResp.Exist {
		R := Result{
			HttpCode: http.StatusOK,
			H: gin.H{
				"Code":    newerror.CodeSuccess,
				"Message": "success",
			},
		}
		Data, _ = sonic.Marshal(R)
	} else {
		R := Result{
			HttpCode: http.StatusOK,
			H: gin.H{
				"Code":    newerror.CodeSuccess,
				"Message": "success",
				"data": gin.H{
					"user_info": gin.H{
						"offline_message_list": kitexResp.JsonData,
					},
				},
			},
		}
		Data, _ = sonic.Marshal(R)
	}
	logger = newlog.AddGateWayInfo(logger, http.StatusOK, userID, ip, fullPath)
	return logger, "GetOfflineMessages", newerror.LevelInfo
}
