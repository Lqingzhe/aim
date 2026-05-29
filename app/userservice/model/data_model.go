package model

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID        int64  `gorm:"primaryKey;comment:user_id"`
	UserName      string `gorm:"NOT NULL;comment:user_name"`
	Introduction  string `gorm:"size:255;comment:introduction"`
	BirthdayYear  int64  `gorm:"comment:birthday_year"`
	BirthdayMonth int64  `gorm:"comment:birthday_month"`
	BirthdayDay   int64  `gorm:"comment:birthday_day"`
}

func (_ *UserInfo) Data() {}

type UserLoginInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID   int64  `gorm:"primaryKey;comment:user_id"`
	Password string `gorm:"size:127;NOT NULL;comment:password"`
	Salt     string `gorm:"size:64;NOT NULL;comment:salt"`
}

func (_ *UserLoginInfo) Data() {}

type RemarkInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID     int64  `gorm:"uniqueIndex:idx_user_id;NOT NULL;comment:user_id"`
	GoalUserID int64  `gorm:"uniqueIndex:idx_user_id;NOT NULL;comment:goal_user_id"`
	NickName   string `gorm:"size:64;comment:nick_name"`
}

func (_ *RemarkInfo) Data() {}
