package groupmember

import (
	"aim/app/groupservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"time"
)

type GroupMemberInfo struct {
	GroupID      int64
	UserID       int64
	Role         commonmodel.GroupRole
	LastReadTime time.Time
}
type GroupMember struct {
	GroupID         int64
	members         []int64
	Role            []commonmodel.GroupRole
	LastReadTime    []time.Time
	Info            []*GroupMemberInfo
	whereWithMember bool
}
type operate func(*GroupMember)

func WithWhereMemberID(info *GroupMember) {
	info.whereWithMember = true
}
func WithVisitTime(time []time.Time) operate {
	return func(info *GroupMember) {
		info.LastReadTime = time
	}
}
func NewStruct(groupID int64, members []int64, role []commonmodel.GroupRole, Operations ...operate) *GroupMember {
	newStruct := &GroupMember{
		GroupID:         groupID,
		members:         members,
		Role:            role,
		whereWithMember: false,
		LastReadTime:    make([]time.Time, len(members)),
		Info:            make([]*GroupMemberInfo, 0),
	}
	if role == nil {
		newStruct.Role = make([]commonmodel.GroupRole, len(members))
	}
	if len(Operations) > 0 {
		for _, Operate := range Operations {
			Operate(newStruct)
		}
	}
	return newStruct
}

func (g *GroupMember) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setRedis(ctx, DB, g)
	if err != nil {
		return err
	}
	return nil
}
func (g *GroupMember) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	if g.whereWithMember {
		exist, err = getRedisWithUser(ctx, DB, g)
	} else {
		exist, err = getRedisWithGroup(ctx, DB, g)
	}
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (g *GroupMember) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getRedisWithUser(ctx, DB, g)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	err = setRedis(ctx, DB, g)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (g *GroupMember) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	if g.whereWithMember {
		err = deleteRedisWithMember(ctx, DB, g)
	} else {
		err = deleteRedisWithGroup(ctx, DB, g)
	}
	return err
}
