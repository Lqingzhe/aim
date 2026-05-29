package model

import (
	"aim/commonmodel"
	"time"

	"gorm.io/gorm"
)

type GroupWithUserInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	GroupID         int64                 `gorm:"NOT NULL;uniqueIndex:group_and_user_id"`
	UserID          int64                 `gorm:"NOT NULL;uniqueIndex:group_and_user_id;index:user_id"`
	GroupRemarkName string                `gorm:"size:255;comment: group_remark_name"`
	Role            commonmodel.GroupRole `gorm:"NOT NULL;size:32;comment: role"`
}

func (_ *GroupWithUserInfo) Data() {}

type GroupInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	GroupID   int64  `gorm:"NOT NULL;primaryKey;index;comment:group_id"`
	GroupName string `gorm:"size:255;NOT NULL;comment:group_name"`
}

func (_ *GroupInfo) Data() {}

type SessionInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	SessionID  int64 `gorm:"NOT NULL;uniqueIndex:session_and_user_id;comment:session_id;index:session_and_goal_user_id"`
	UserID     int64 `gorm:"NOT NULL;uniqueIndex:session_and_user_id;comment:user_id"`
	GoalUserID int64 `gorm:"NOT NULL;index:session_and_goal_user_id;comment:goal_user_id"`
}

func (_ *SessionInfo) Data() {}

type GroupApplyInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	GoalID      int64 `gorm:"NOT NULL;uniqueIndex:group_and_user_id;comment:goal_id"`
	ApplyUserID int64 `gorm:"NOT NULL;uniqueIndex:group_and_user_id;comment:apply_user_id"`
}

func (_ *GroupApplyInfo) Data() {}

type GroupMuteInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	GroupID     int64     `gorm:"NOT NULL;uniqueIndex:group_and_user_id;comment:group_id"`
	UserID      int64     `gorm:"NOT NULL;uniqueIndex:group_and_user_id;comment:user_id"`
	MuteEndTime time.Time `gorm:"NOT NULL;comment:mute_end_time"`
	MuteReason  string    `gorm:"size:255;comment:mute_reason"`
}

func (_ *GroupMuteInfo) Data() {}
