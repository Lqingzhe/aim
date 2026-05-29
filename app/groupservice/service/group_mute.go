package service

import (
	"aim/app/groupservice/dao"
	"aim/app/groupservice/dao/groupmember"
	"aim/app/groupservice/dao/groupmuteinfo"
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
)

type Mute struct {
	traceID          string
	groupNoticeTopic sarama.SyncProducer
	dbContext        *model.DBContext
	groupConfig      commonmodel.GroupConfig
}

func NewMute(traceID string, groupNoticeTopic sarama.SyncProducer, dbContext *model.DBContext, groupConfig commonmodel.GroupConfig) *Mute {
	return &Mute{
		traceID:          traceID,
		groupNoticeTopic: groupNoticeTopic,
		dbContext:        dbContext,
		groupConfig:      groupConfig,
	}
}

func (m *Mute) SetMute(ctx context.Context, userID int64, groupID int64, goalUserID int64, MuteTimeSecond int64, MuteReason string) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group_mute:SetMute")
	if MuteTimeSecond > int64(m.groupConfig.MaxGroupMuteTime.Seconds()) {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Mute Time Too Long", fmt.Errorf("Mute Time Too Long"), newerror.LevelInfo)
	}
	if tool.CalculateLength(MuteReason) > m.groupConfig.MaxGroupMuteReasonLength {
		return newerror.MakeError(http.StatusBadRequest, newerror.CodeParamValueInvalid, "Mute Reason Too Long", fmt.Errorf("Mute Reason Too Long"), newerror.LevelInfo)
	}
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID, goalUserID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, m.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Set Mute By Uninjion User"), newerror.LevelInfo)
	}
	userRole := groupMemberStruct.Info[0].Role
	goalUserRole := groupMemberStruct.Info[1].Role
	if userRole == commonmodel.Member || (userRole == commonmodel.Manager && goalUserRole != commonmodel.Member) {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not The Manager", fmt.Errorf("Try To Set Mute Without Enough Permission"), newerror.LevelInfo)
	}
	groupMuteStruct := groupmuteinfo.NewStruct(groupID, goalUserID, time.Now().Add(time.Duration(MuteTimeSecond)*time.Second), MuteReason, groupmuteinfo.WithWhereUserID)
	exist, err = dao.Update(ctx, groupMuteStruct, m.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		err = dao.Add(ctx, groupMuteStruct, m.dbContext)
		if err != nil {
			return err
		}
	}
	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, m.dbContext)
	if err != nil {
		return err
	}
	memberList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberList = append(memberList, info.UserID)
	}
	groupNoticeMessage := commonmodel.KafkaGroupNotice{
		TraceID:        m.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberList,
		SessionID:      groupID,
		Data:           map[string]any{"user_id": userID, "goal_user_id": goalUserID, "mute_time": MuteTimeSecond, "mute_reason": MuteReason},
		MessageCode:    commonmodel.MessageCode_GroupSetMute,
	}
	_, _, err = tool.SendKafkaGroupNotice(m.groupNoticeTopic, groupNoticeMessage)
	if err != nil {
		return err
	}
	return nil
}
func (m *Mute) ReleaseMute(ctx context.Context, userID int64, groupID int64, goalUserID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group_mute:ReleaseMute")
	groupMemberStruct := groupmember.NewStruct(groupID, []int64{userID}, nil, groupmember.WithWhereMemberID)
	exist, err := dao.Get(ctx, groupMemberStruct, m.dbContext)
	if err != nil {
		return err
	}
	if !exist {
		return newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Group Is Not Exist", fmt.Errorf("Try To Set Mute By Uninjion User"), newerror.LevelInfo)
	}
	if groupMemberStruct.Info[0].Role == commonmodel.Member {
		return newerror.MakeError(http.StatusForbidden, newerror.CodePermissionDenied, "You Are Not Manager", fmt.Errorf("Try To Realse Mute Without Enough Permission"), newerror.LevelInfo)
	}
	groupMuteStruct := groupmuteinfo.NewStruct(groupID, goalUserID, time.Time{}, "")
	err = dao.Delete(ctx, groupMuteStruct, m.dbContext)
	if err != nil {
		return err
	}
	groupMemberStruct = groupmember.NewStruct(groupID, nil, nil)
	_, err = dao.Get(ctx, groupMemberStruct, m.dbContext)
	if err != nil {
		return err
	}
	memberList := make([]int64, 0, len(groupMemberStruct.Info))
	for _, info := range groupMemberStruct.Info {
		memberList = append(memberList, info.UserID)
	}
	groupNoticeMessage := commonmodel.KafkaGroupNotice{
		TraceID:        m.traceID,
		SendTimeSecond: time.Now().Unix(),
		GoalUserID:     memberList,
		SessionID:      groupID,
		Data:           map[string]any{"user_id": userID, "goal_user_id": goalUserID},
		MessageCode:    commonmodel.MessageCode_GroupReleaseMute,
	}
	_, _, err = tool.SendKafkaGroupNotice(m.groupNoticeTopic, groupNoticeMessage)
	if err != nil {
		return err
	}
	return nil
}
func (m *Mute) GetMuteStatus(ctx context.Context, userID int64, groupID int64) (muteReason string, muteEndTime string, isMute bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("group_mute:GetMuteStatus")
	var finalErr error
	groupMuteStruct := groupmuteinfo.NewStruct(groupID, userID, time.Time{}, "", groupmuteinfo.WithWhereUserID)
	exist, err := dao.Get(ctx, groupMuteStruct, m.dbContext)
	if newerror.WhetherInterrupt(err, &finalErr) {
		return "", "", false, err
	}
	if !exist {
		return "", "", false, finalErr
	}
	if groupMuteStruct.Info[0].MuteEndTime.Unix() > time.Now().Unix() {
		return groupMuteStruct.Info[0].MuteReason, groupMuteStruct.Info[0].MuteEndTime.String(), true, finalErr
	} else {
		return "", "", false, finalErr
	}
}
