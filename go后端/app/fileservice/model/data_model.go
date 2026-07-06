package model

import (
	"aim/commonmodel"
	"time"

	"gorm.io/gorm"
)

type FileModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FileID      int64                `gorm:"primaryKey;autoIncrement:false;NOT NULL;comment: file_id"`
	FileName    string               `gorm:"size:255;comment: file_name"`
	FileType    commonmodel.FileType `gorm:"NOT NULL;size:64;comment: file_type"`
	ContentType string               `gorm:"size:255;NOT NULL;comment: content_type"`

	VoiceDurationSecond int64  `gorm:"comment: voice_duration_second"`
	StoragePath         string `gorm:"size:255;NOT NULL;comment: storage_path"`

	DataStream []byte `gorm:"-"`
}

func (_ *FileModel) Data() {}
