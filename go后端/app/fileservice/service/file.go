package service

import (
	"aim/app/fileservice/dao"
	"aim/app/fileservice/dao/file"
	"aim/app/fileservice/model"
	"aim/commonmodel"
	newerror "aim/pkg/error"
	"aim/tool"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/bwmarrin/snowflake"
)

type FileService struct {
	snowFlake  *snowflake.Node
	dbContext  *model.DBContext
	fileConfig commonmodel.FileConfig
}

func NewFileService(snowFlake *snowflake.Node, dbContext *model.DBContext, fileConfig commonmodel.FileConfig) *FileService {
	return &FileService{
		snowFlake:  snowFlake,
		dbContext:  dbContext,
		fileConfig: fileConfig,
	}
}
func (f *FileService) CreateFile(ctx context.Context, fileName string, dataStream []byte, fileType commonmodel.FileType, contentType string, voiceDurationTimeSecond int64) (fileID int64, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("file:CreateFile")
	fileID = f.snowFlake.Generate().Int64()
	storagePath := tool.MakeFileStoragePath(f.fileConfig.FileStoragePath, fileID, filepath.Ext(fileName))
	fileStruct := file.NewStruct(fileID, fileName, fileType, dataStream, storagePath, contentType, voiceDurationTimeSecond)
	err = dao.Add(ctx, fileStruct, f.dbContext)
	if err != nil {
		return 0, err
	}

	return fileID, nil
}

func (f *FileService) DeleteFile(ctx context.Context, fileID int64) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("file:DeleteFile")
	fileStruct := file.NewStructWithFileID(fileID)
	err = dao.Delete(ctx, fileStruct, f.dbContext)
	if err != nil {
		return err
	}
	return nil
}
func (f *FileService) GetFile(ctx context.Context, fileID int64) (dataStream []byte, contentType string, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("file:GetFile")
	fileStruct := file.NewStructWithFileID(fileID)
	exist, err := dao.Get(ctx, fileStruct, f.dbContext)
	if err != nil {
		return nil, "", err
	}
	if !exist {
		return nil, "", newerror.MakeError(http.StatusNotFound, newerror.CodeResourceNotFound, "The Rescource Is Not Exist", fmt.Errorf("Try To Get Unexist File"), newerror.LevelInfo)
	}
	getInfo := fileStruct.Info[0]
	return getInfo.DataStream, getInfo.ContentType, nil
}
