package file

import (
	"aim/app/fileservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
)

type FileInfo struct {
	model.FileModel
	Info []*model.FileModel
}

func NewStruct(FileID int64, FileName string, FileType commonmodel.FileType, DataStream []byte, StoragePath string, ContentType string, VoiceTimeSecond int64) *FileInfo {
	return &FileInfo{
		FileModel: model.FileModel{
			FileID:              FileID,
			FileName:            FileName,
			FileType:            FileType,
			DataStream:          DataStream,
			StoragePath:         StoragePath,
			ContentType:         ContentType,
			VoiceDurationSecond: VoiceTimeSecond,
		},
		Info: make([]*model.FileModel, 0),
	}
}
func NewStructWithFileID(FileID int64) *FileInfo {
	return &FileInfo{
		FileModel: model.FileModel{
			FileID: FileID,
		},
		Info: make([]*model.FileModel, 0),
	}
}
func (f *FileInfo) AddInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:AddInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = setOs(f)
	if err != nil {
		return err
	}
	err = setMysql(ctx, DB, f)
	if err != nil {
		_ = deleteOs(f)
		return err
	}
	return nil
}
func (f *FileInfo) GetInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:GetInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return false, err
	}
	exist, err = getMysql(ctx, DB, f)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}
	newStruct := FileInfo{
		FileModel: *f.Info[0],
		Info:      make([]*model.FileModel, 1),
	}
	exist, err = getOs(&newStruct)
	if err != nil {
		return false, err
	}
	if !exist {
		_ = deleteMysql(ctx, DB, &newStruct)
		return false, nil
	}
	f.Info = append(f.Info, newStruct.Info[0])
	return true, nil
}
func (f *FileInfo) UpdateInfo(ctx context.Context, dbContext any) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:UpdateInfo")
	return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeInternalError, "Useless Module Unexpectedly Used", fmt.Errorf("%s", "Useless Module Unexpectly Used"), newerror.LevelFatal)
}
func (f *FileInfo) DeleteInfo(ctx context.Context, dbContext any) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("db_manage:DeleteInfo")
	DB, err := tool.TypeAssert[model.DBContext](dbContext)
	if err != nil {
		return err
	}
	err = deleteOs(f)
	if err != nil {
		return err
	}
	err = deleteMysql(ctx, DB, f)
	if err != nil {
		return err
	}
	return nil
}
