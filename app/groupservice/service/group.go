package service

import (
	"aim/app/groupservice/dao"
	"aim/app/groupservice/dao/groupapply"
	"aim/app/groupservice/dao/groupinfo"
	"aim/app/groupservice/dao/groupmember"
	"aim/app/groupservice/dao/groupwithuser"
	"aim/app/groupservice/model"
	"aim/commonmodel"
	"aim/kitex_gen/kitexcommonmodel"
	"aim/kitex_gen/kitexmessageservice"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
)

type ServiceGroup struct {
	traceID          string
	groupNoticeTopic sarama.SyncProducer
	systemTopic      sarama.SyncProducer
	dbContext        *model.DBContext
	groupConfig      commonmodel.GroupConfig
	snowFlack        *snowflake.Node
	serviceClient    model.ServiceClient
}

func NewGroup(traceID string, groupNoticeTopic sarama.SyncProducer, systemTopic sarama.SyncProducer, dbContext *model.DBContext, groupConfig commonmodel.GroupConfig, snowFlack *snowflake.Node, ServiceClient model.ServiceClient) *ServiceGroup {
	return &ServiceGroup{
		traceID:          traceID,
		groupNoticeTopic: groupNoticeTopic,
		systemTopic:      systemTopic,
		dbContext:        dbContext,
		groupConfig:      groupConfig,
		snowFlack:        snowFlack,
		serviceClient:    ServiceClient,
	}
}

func (s *ServiceGroup) GetGroupInfo(ctx context.Context, groupID int64) (groupInfoResp *model.GroupInfo, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:GetGroupInfo")
	groupStruct := groupinfo.NewStruct(groupID, "")
	exist, err := dao.Get(ctx, groupStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Get Unexist Group Info"), newerror.LevelInfo)
	}
	return groupStruct.Info[0], nil
}
func (s *ServiceGroup) ChangeGroupInfo(ctx context.Context, groupID int64, userID int64, groupName string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:ChangeGroupInfo")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Update Group Info About Unexist Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role == commonmodel.Member {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Do Not Have Enough Permission", fmt.Errorf("Group Member Try To Change Group Info"), newerror.LevelInfo)
	}
	if tool.CalculateLength(groupName) > s.groupConfig.MaxGroupNameLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Group Name Too Long", fmt.Errorf("Group Name Too Long"), newerror.LevelInfo)
	}
	groupStruct := groupinfo.NewStruct(groupID, groupName)
	_, err = dao.Update(ctx, groupStruct, s.dbContext)
	if err != nil {
		return err
	}
	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberList = append(memberList, info.UserID)
	}
	groupNoticeMessage := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberList,
		SessionID:      groupID,
		Data:           map[string]any{"user_id": userID},
		MessageCode:    commonmodel.MessageCode_GroupInfoChange,
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupNoticeMessage)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) SearchGroup(ctx context.Context, groupName string) (groupID []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:SearchGroup")
	if tool.CalculateLength(groupName) > s.groupConfig.MaxGroupNameLength {
		return nil, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Group Name Too Long", fmt.Errorf("Group Name Too Long"), newerror.LevelInfo)
	}
	groupStruct := groupinfo.NewStruct(0, groupName, groupinfo.WithGroupName)
	exist, err := dao.Get(ctx, groupStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "Do Not Find Group Info", fmt.Errorf("Do Not Search The Groups"), newerror.LevelInfo)
	}
	groupID = make([]int64, 0, len(groupStruct.Info))
	for _, info := range groupStruct.Info {
		groupID = append(groupID, info.GroupID)
	}
	return groupID, nil
}

func (s *ServiceGroup) CreateGroup(ctx context.Context, userID int64, groupName string) (groupID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:CreateGroup")
	if tool.CalculateLength(groupName) > s.groupConfig.MaxGroupNameLength {
		return 0, newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Group Name Too Long", fmt.Errorf("Group Name Too Long"), newerror.LevelInfo)
	}
	groupID = s.snowFlack.Generate().Int64()
	groupStruct := groupinfo.NewStruct(groupID, groupName)
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", commonmodel.GroupOwner)
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, []commonmodel.GroupRole{commonmodel.GroupOwner}, nil)
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	err = dao.Add(ctx, groupStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	err = dao.Add(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	err = dao.Add(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return 0, err
	}
	result := DB.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return 0, err2
	}
	return groupID, nil
}
func (s *ServiceGroup) DeleteGroup(ctx context.Context, userID int64, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:DeleteGroup")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusUnauthorized, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Disband Uninjoin Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role != commonmodel.GroupOwner {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not The Group Owner", fmt.Errorf("Unauthorized Operat Group"), newerror.LevelInfo)
	}
	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberID := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberID = append(memberID, info.UserID)
	}
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, 0, "", "", groupwithuser.WithGroupID)
	err = dao.Delete(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	groupInfo := groupinfo.NewStruct(groupID, "")
	err = dao.Delete(ctx, groupInfo, DB)
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
		GroupId:    groupID,
	}
	_, err = s.serviceClient.MessageClient.DeleteMessageAllGroup(ctx, &deleteMessageAllGroupReq)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	result := DB.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberID,
		Data:           map[string]any{"user_id": userID},
		MessageCode:    commonmodel.MessageCode_GroupDisband,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceGroup) LeaveGroup(ctx context.Context, groupID int64, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:LeaveGroup")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusUnauthorized, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Leave Uninjoin Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role == commonmodel.GroupOwner {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "Please Transfer The Role Of Group Owner", fmt.Errorf("Group Owner Try Leave Group"), newerror.LevelInfo)
	}
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	err = dao.Delete(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = dao.Delete(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = DB.Mysql.Client.Commit().Error
	if err != nil {
		return err
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberID := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberID = append(memberID, info.UserID)
	}
	messageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberID,
		SessionID:      groupID,
		Data:           map[string]any{"user_id": userID},
		MessageCode:    commonmodel.MessageCode_GroupLeave,
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) SetGroupApply(ctx context.Context, groupID int64, userID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:PullGroupApply")
	groupWithUserStruct := groupwithuser.NewStruct(groupID, 0, "", "", groupwithuser.WithGroupID)
	exist, err := dao.Get(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "Group Is Not Exist", fmt.Errorf("Pull Apply To Unexist Group"), newerror.LevelInfo)
	}
	groupApplyStruct := groupapply.NewStruct(groupID, userID)
	err = dao.Add(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	managerIDs := make([]int64, 0, len(groupWithUserStruct.Info))
	for _, info := range groupWithUserStruct.Info {
		if info.Role == commonmodel.GroupOwner || info.Role == commonmodel.Manager {
			managerIDs = append(managerIDs, info.UserID)
		}
	}
	messageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     managerIDs,
		Data:           map[string]any{"user_id": userID, "group_id": groupID},
		MessageCode:    commonmodel.MessageCode_GroupApply,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, messageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) GetGroupApplyList(ctx context.Context, groupID int64, userID int64) (applyUserID []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:GetGroupApplyList")
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err := dao.Get(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "Group Is Not Exist", fmt.Errorf("Get Apply List From Not Exist Group"), newerror.LevelInfo)
	}
	if groupWithUserStruct.Info[0].Role != commonmodel.GroupOwner && groupWithUserStruct.Info[0].Role != commonmodel.Manager {
		return nil, newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Do Not Have Enough Permission", fmt.Errorf("Group Member Try To Get Group Apply List"), newerror.LevelInfo)
	}
	groupApplyStruct := groupapply.NewStruct(groupID, 0, groupapply.WithGoalID)
	exist, err = dao.Get(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	applyUserID = make([]int64, len(groupApplyStruct.Info))
	for i, info := range groupApplyStruct.Info {
		applyUserID[i] = info.ApplyUserID
	}
	return applyUserID, nil
}
func (s *ServiceGroup) GetLastVisitTime(ctx context.Context, groupID int64, userID int64) (userIDList []int64, lastVisitTime []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:GetLastVisitTime")
	groupMemberStruct := groupmember.NewStruct(groupID, nil, nil)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return nil, nil, err
	}
	if !exist {
		return nil, nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Get Last Visit Time From Unexist Group"), newerror.LevelInfo)
	}
	whetherUserInGroup := false
	LastVisitTimeList := make([]int64, 0, len(groupMemberStruct.Info)-1)
	memberIDList := make([]int64, 0, len(groupMemberStruct.Info)-1)
	for _, info := range groupMemberStruct.Info {
		if info.UserID == userID {
			whetherUserInGroup = true
		}
		LastVisitTimeList = append(LastVisitTimeList, info.LastReadTime.Unix())
		memberIDList = append(memberIDList, info.UserID)
	}
	if !whetherUserInGroup {
		return nil, nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Uninjoin User Try To Get Group Info "), newerror.LevelInfo)
	}
	return memberIDList, LastVisitTimeList, nil
}
func (s *ServiceGroup) AgreeGroupApply(ctx context.Context, groupID int64, userID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:AgreeGroupApply")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Get Last Visit Time From Unexist Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role == commonmodel.Member {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Do Not Have Enough Permission", fmt.Errorf("Group Member Try To Agree Group Apply"), newerror.LevelInfo)
	}

	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupApplyStruct := groupapply.NewStruct(groupID, goalUserID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	groupWithUserStruct := groupwithuser.NewStruct(groupID, goalUserID, "", commonmodel.Member)
	groupMemberStruct = groupmember.NewStruct(groupID, []int64{goalUserID}, []commonmodel.GroupRole{commonmodel.Member}, groupmember.WithWhereMemberID)
	err = dao.Add(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = dao.Delete(ctx, groupApplyStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = dao.Add(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	result := DB.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	exist, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Agree Group Apply From Unexist Group"), newerror.LevelInfo)
	}
	memberID := make([]int64, len(groupMemberStruct.Info))
	for i, info := range groupMemberStruct.Info {
		memberID[i] = info.UserID
	}

	groupMessageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberID,
		SessionID:      groupID,
		Data:           map[string]interface{}{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_GroupJoin,
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"group_id": groupID},
		MessageCode:    commonmodel.MessageCode_GroupJoin,
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) RefuseGroupApply(ctx context.Context, groupID int64, userID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:RefuseGroupApply")
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err := dao.Get(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Transform Group Owner From Unexist Group"), newerror.LevelInfo)
	}
	if groupWithUserStruct.Info[0].Role != commonmodel.GroupOwner {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not The Group Owner", fmt.Errorf("Transform Group Owner By Member"), newerror.LevelInfo)
	}
	groupApplyStruct := groupapply.NewStruct(groupID, goalUserID, groupapply.WithGoalID, groupapply.WithApplyUserID)
	err = dao.Delete(ctx, groupApplyStruct, s.dbContext)
	if err != nil {
		return err
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"group_id": groupID},
		MessageCode:    commonmodel.MessageCode_GroupRefuse,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) TransformGroupOwner(ctx context.Context, groupID int64, userID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:TransformGroupOwner")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Set Manager With Unexist Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role != commonmodel.GroupOwner {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not The Group Owner", fmt.Errorf("Transform Group Owner By Member"), newerror.LevelInfo)
	}

	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, goalUserID, "", commonmodel.GroupOwner)
	exist, err = dao.Update(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	if !exist {
		DB.Mysql.Client.Rollback()
		return newerror.MakeError(http.StatusNotFound, newerror.CodeUserNotFound, "The User Did Not Join The Group", fmt.Errorf("Transform Group Owner To Unexpected Member"), newerror.LevelInfo)
	}
	groupWithUserStruct = groupwithuser.NewStruct(groupID, userID, "", commonmodel.Member)
	_, err = dao.Update(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	groupMemberStruct = groupmember.NewStruct(groupID, []int64{userID, goalUserID}, []commonmodel.GroupRole{commonmodel.Member, commonmodel.GroupOwner}, groupmember.WithWhereMemberID)
	_, err = dao.Update(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	result := DB.Mysql.Client.Commit()
	if err2 := newerror.IsMysqlError(result); err2 != nil {
		return err2
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberIDList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberIDList = append(memberIDList, info.UserID)
	}
	groupMessageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberIDList,
		SessionID:      groupID,
		Data:           map[string]interface{}{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_TransformGroupOwner,
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"group_id": groupID},
		MessageCode:    commonmodel.MessageCode_TransformGroupOwner,
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) KickOutGroup(ctx context.Context, userID int64, goalUserID int64, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:KickOutGroup")
	if userID == goalUserID {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeInvalidParam, "You Can not Kick Yourself", fmt.Errorf("Try To Kick User's ownself"), newerror.LevelInfo)
	}
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID, goalUserID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Set Manager With Unexist Group"), newerror.LevelInfo)
	}
	userRole := groupMemberStruct.Info[0].Role
	goalUserRole := groupMemberStruct.Info[1].Role
	if userRole == commonmodel.Member || (userRole == commonmodel.Manager && goalUserRole != commonmodel.Member) {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Only Can Kick Out The One Under You", fmt.Errorf("Do Not Have Enough Permission To Kick Out Member"), newerror.LevelInfo)
	}

	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, goalUserID, "", "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	groupMemberStruct = groupmember.NewStruct(groupID, []int64{goalUserID}, nil, groupmember.WithWhereMemberID)
	err = dao.Delete(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = dao.Delete(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = DB.Mysql.Client.Commit().Error
	if err != nil {
		return err
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberIDList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberIDList = append(memberIDList, info.UserID)
	}
	groupMessageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberIDList,
		SessionID:      groupID,
		Data:           map[string]interface{}{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_GroupKick,
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": userID, "group_id": groupID},
		MessageCode:    commonmodel.MessageCode_GroupKick,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) SetManager(ctx context.Context, userID int64, goalUserID int64, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:SetManager")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Set Manager With Unexist Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role == commonmodel.Member {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not Manager", fmt.Errorf("Member Try To Set Manager"), newerror.LevelInfo)
	}
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, goalUserID, "", commonmodel.Manager, groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err = dao.Update(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	if !exist {
		DB.Mysql.Client.Rollback()
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeUserNotFound, "The User Is Not In Group", fmt.Errorf("Set Manager With Unexist User In Group"), newerror.LevelInfo)
	}
	groupMemberStruct = groupmember.NewStruct(groupID, []int64{goalUserID}, []commonmodel.GroupRole{commonmodel.Manager}, groupmember.WithWhereMemberID)
	_, err = dao.Update(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = DB.Mysql.Client.Commit().Error
	if err != nil {
		return err
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberIDList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberIDList = append(memberIDList, info.UserID)
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": userID, "group_id": groupID},
		MessageCode:    commonmodel.MessageCode_SetGroupManager,
	}
	groupMessageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberIDList,
		SessionID:      groupID,
		Data:           map[string]interface{}{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_SetGroupManager,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) RevokeManager(ctx context.Context, userID int64, goalUserID int64, groupID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:RevokeManager")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Revoke Manager With Unexist Group"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role != commonmodel.GroupOwner {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not The Group Owner", fmt.Errorf("Try To Revoke Manager Without Enough Permission"), newerror.LevelInfo)
	}
	DB := &model.DBContext{
		Mysql: tool.BeginMysqlTransaction(s.dbContext.Mysql),
		Redis: s.dbContext.Redis,
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", commonmodel.Member, groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err = dao.Update(ctx, groupWithUserStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeUserNotFound, "The Member Is Not In The Group", fmt.Errorf("Revoke Manager With Unexist Group"), newerror.LevelInfo)
	}
	groupMemberStruct = groupmember.NewStruct(groupID, []int64{goalUserID}, []commonmodel.GroupRole{commonmodel.Member}, groupmember.WithWhereMemberID)
	_, err = dao.Update(ctx, groupMemberStruct, DB)
	if err != nil {
		DB.Mysql.Client.Rollback()
		return err
	}
	err = DB.Mysql.Client.Commit().Error
	if err != nil {
		return err
	}

	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return err
	}
	memberIDList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberIDList = append(memberIDList, info.UserID)
	}
	systemMessageStruct := commonmodel.KafkaSystemMessage{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     []int64{goalUserID},
		Data:           map[string]any{"user_id": userID, "group_id": groupID},
		MessageCode:    commonmodel.MessageCode_RevokeGroupManager,
	}
	groupMessageStruct := commonmodel.KafkaGroupNotice{
		TraceID:        s.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberIDList,
		SessionID:      groupID,
		Data:           map[string]interface{}{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_RevokeGroupManager,
	}
	_, _, err = tool.SendKafkaSystemMessage(s.systemTopic, systemMessageStruct)
	if err != nil {
		return err
	}
	_, _, err = tool.SendKafkaGroupNotice(s.groupNoticeTopic, groupMessageStruct)
	if err != nil {
		return err
	}
	return nil
}
func (s *ServiceGroup) GetGroupInfoWithUser(ctx context.Context, groupID int64, userID int64) (groupWithUserInfo *model.GroupWithUserInfo, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:GetGroupInfoWithUser")
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, "", "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err := dao.Get(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "You Did Not Join The Group", fmt.Errorf("User Try To Get User Info About Group Without Joining The Group"), newerror.LevelInfo)
	}
	return groupWithUserStruct.Info[0], nil
}
func (s *ServiceGroup) UpdateGroupInfoWithUser(ctx context.Context, userID int64, groupID int64, groupRemarkName string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:UpdateGroupInfoWithUser")
	if tool.CalculateLength(groupRemarkName) > s.groupConfig.MaxGroupNickNameLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Group Remark Name Too Long", fmt.Errorf("Group Remark Name Too Long"), newerror.LevelInfo)
	}
	groupWithUserStruct := groupwithuser.NewStruct(groupID, userID, groupRemarkName, "", groupwithuser.WithGroupID, groupwithuser.WithUserID)
	exist, err := dao.Update(ctx, groupWithUserStruct, s.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Update Group With User Info Without Joining The Group"), newerror.LevelInfo)
	}
	return nil
}
func (s *ServiceGroup) GetGroupUserID(ctx context.Context, groupID int64, userID int64) (userIDList []int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group:GetGroupUserID")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Get Group User ID List Without Joining The Group"), newerror.LevelInfo)
	}
	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, s.dbContext)
	if err != nil {
		return nil, err
	}
	userIDList = make([]int64, 0, len(groupMemberStruct.Info)-1)
	for _, info := range groupMemberStruct.Info {
		if info.UserID != userID {
			userIDList = append(userIDList, info.UserID)
		}
	}
	return userIDList, nil
}
