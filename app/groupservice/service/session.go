package service

import (
	"aim/app/groupservice/dao"
	"aim/app/groupservice/dao/groupapply"
	"aim/app/groupservice/dao/groupmember"
	"aim/app/groupservice/dao/sessioninfo"
	"aim/app/groupservice/model"
	"aim/commonmodel"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	"aim/kitex_gen/kitexuserservice"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
)

type ServiceSession struct {
	traceID       string
	systemTopic   sarama.SyncProducer
	dbContext     *model.DBContext
	snowFlack     *snowflake.Node
	serviceClient model.ServiceClient
}

func NewSession(traceID string, systemTopic sarama.SyncProducer, dbContext *model.DBContext, snowFlack *snowflake.Node, serviceClient model.ServiceClient) *ServiceSession {
	return &ServiceSession{
		traceID:       traceID,
		systemTopic:   systemTopic,
		dbContext:     dbContext,
		snowFlack:     snowFlack,
		serviceClient: serviceClient,
	}
}
func (s *ServiceSession) CreatSession(ctx context.Context, userID int64, goalUserID int64) (SessionID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:CreatSession")
	SessionID = s.snowFlack.Generate().Int64()
	sessionStruct1 := sessioninfo.NewStruct(SessionID, userID, goalUserID, sessioninfo.WithGoalUserID, sessioninfo.WithUserID)
	sessionStruct2 := sessioninfo.NewStruct(SessionID, goalUserID, userID, sessioninfo.WithGoalUserID, sessioninfo.WithUserID)
	exist, err := dao.Get(ctx, sessionStruct1, s.dbContext)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, newerror.MakeError(http.StatusTooManyRequests, newerror.CodeResourceDuplicate, "You Are Already Being Friends", fmt.Errorf("Try To Repeat Make Friend"), newerror.LevelInfo)
	}

	groupMemberStruct := groupmember.NewStruct(SessionID, []int64{userID, goalUserID}, nil)
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	err = dao.Add(ctx, sessionStruct1, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	err = dao.Add(ctx, sessionStruct2, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	err = dao.Add(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	groupApplyStruct1 := groupapply.NewStruct(userID, goalUserID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	err = dao.Delete(ctx, groupApplyStruct1, s.dbContext)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	groupApplyStruct2 := groupapply.NewStruct(goalUserID, userID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	err = dao.Delete(ctx, groupApplyStruct2, s.dbContext)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	result := DB.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return 0, err2
	}
	//向s.GoalUserID发送消息
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": strconv.FormatInt(userID, 10)},
		MessageCode:    commonmodel.MessageCode_FriendRequest_Success,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return 0, err
	}
	return SessionID, nil
}

func (s *ServiceSession) DeleteSession(ctx context.Context, sessionID int64, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:DeleteSession")
	sessionStruct := sessioninfo.NewStruct(sessionID, userID, 0, sessioninfo.WithSessionID, sessioninfo.WithUserID)
	exist, err := dao.Get(ctx, sessionStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeResourceNotFound, "He IS Not Your Friend", fmt.Errorf("Delete Unexist SessionID"), newerror.LevelInfo)
	}
	goalUserID := sessionStruct.Info[0].GoalUserID
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupMemberStruct := groupmember.NewStruct(sessionID, nil, nil)
	sessionStruct = sessioninfo.NewStruct(sessionID, 0, 0, sessioninfo.WithSessionID)
	err = dao.Delete(ctx, sessionStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = dao.Delete(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	deleteMessageAllGroupReq := kitexmessageservice.DeleteMessageAllGroupReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: s.traceID},
		GroupId:    sessionID,
	}
	_, err = s.serviceClient.MessageClient.DeleteMessageAllGroup(ctx, &deleteMessageAllGroupReq)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return newerror.UnMarshalError(err)
	}
	err = DB.Mysql.Client.Commit().Error
	if err != nil {
		return err
	}
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": strconv.FormatInt(userID, 10)},
		MessageCode:    commonmodel.MessageCode_FriendDelete,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceSession) GetFriendLastVisitTime(ctx context.Context, sessionID int64, goalUserID int64) (lastVisitTime string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:GetFriendLastVisitTime")
	groupMemberStruct := groupmember.NewStruct(sessionID, []int64{goalUserID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", newerror.MakeError(http.StatusBadRequest, newerror.CodeResourceNotFound, "He Is Not Your Friend", fmt.Errorf("Do Not Get The Info Whith SessionID and UserID"), newerror.LevelInfo)
	}
	return groupMemberStruct.Info[0].LastReadTime.String(), nil
}
func (s *ServiceSession) ApplyForFriend(ctx context.Context, userID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:ApplyForFriend")
	var finalErr error
	getOtherUserInfoReq := kitexuserservice.GetOtherUserInfoReq{
		CommonInfo: &kitexcommonmodel.CommonInfo{Trace: s.traceID},
		GoalUserId: goalUserID,
	}
	_, err = s.serviceClient.UserClient.GetOtherUserInfo(ctx, &getOtherUserInfoReq)
	if newerror.WhetherInterrupt(newerror.UnMarshalError(err), &finalErr) {
		return finalErr
	}
	groupApplyStruct := groupapply.NewStruct(goalUserID, userID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	exist, err := dao.Get(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	if exist {
		return newerror.MakeError(http.StatusTooManyRequests, newerror.CodeResourceDuplicate, "Friend Apply Is Already Exist", fmt.Errorf("Repeat Send Friend Apply"), newerror.LevelInfo)
	}
	err = dao.Add(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": strconv.FormatInt(userID, 10)},
		MessageCode:    commonmodel.MessageCode_FriendRequest,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceSession) GetFriendApplyList(ctx context.Context, userID int64) (applyUserID []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:GetApplyList")
	groupApplyStruct := groupapply.NewStruct(userID, 0, groupapply.WithGoalID)
	exist, err := dao.Get(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	applyUserID = make([]int64, len(groupApplyStruct.Info))
	for i, v := range groupApplyStruct.Info {
		applyUserID[i] = v.ApplyUserID
	}
	return applyUserID, nil
}
func (s *ServiceSession) RefuseFriendApply(ctx context.Context, userID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("session:RefuseFriendApply")
	groupApplyStruct := groupapply.NewStruct(userID, goalUserID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	err = dao.Delete(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	groupApplyStruct = groupapply.NewStruct(goalUserID, userID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	err = dao.Delete(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": strconv.FormatInt(userID, 10)},
		MessageCode:    commonmodel.MessageCode_FriendRequest_Refuse,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}
