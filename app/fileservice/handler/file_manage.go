package handler

import (
	"aim/app/fileservice/service"
	"aim/commonmodel"
	"aim/kitex_gen/kitexfileservice"
	newerror "aim/pkg/error"
	newlog "aim/pkg/log"
	"context"
)

func (s *KitexFileServiceImpl) CreateFile(ctx context.Context, req *kitexfileservice.CreateFileReq) (resp *kitexfileservice.CreateFileResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	fileID, err := serviceStruct.CreateFile(ctx, req.FileName, req.DataStream, commonmodel.FileType(req.FileType), req.ContentType, req.VoiceDurationTimeSecond)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "CreateFile")
		return nil, err
	}
	resp = &kitexfileservice.CreateFileResp{
		FileId: fileID,
	}
	return resp, nil
}

func (s *KitexFileServiceImpl) DeleteFile(ctx context.Context, req *kitexfileservice.DeleteFileReq) (resp *kitexfileservice.DeleteFileResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	err = serviceStruct.DeleteFile(ctx, req.FileId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "DeleteFile")
		return nil, err
	}
	return &kitexfileservice.DeleteFileResp{}, nil
}

func (s *KitexFileServiceImpl) GetFile(ctx context.Context, req *kitexfileservice.GetFileReq) (resp *kitexfileservice.GetFileResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	logger := newlog.AddTraceID(s.logger, req.CommonInfo.Trace)
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	dataStream, contentType, err := serviceStruct.GetFile(ctx, req.FileId)
	if err != nil {
		err2 := newerror.TranslateError(err)
		logger = newlog.AddError(logger, err2, err2.StatusCode)
		newlog.Log(logger, err2.LogLevel, "GetFile")
		return nil, err
	}
	resp = &kitexfileservice.GetFileResp{
		DataStream:  dataStream,
		ContentType: contentType,
	}
	return resp, nil
}
