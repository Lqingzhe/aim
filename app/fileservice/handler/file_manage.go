package handler

import (
	"aim/app/fileservice/service"
	"aim/commonmodel"
	"aim/kitex_gen/kitexfileservice"
	newerror "aim/pkg/error"
	"context"
)

func (s *KitexFileServiceImpl) CreateFile(ctx context.Context, req *kitexfileservice.CreateFileReq) (resp *kitexfileservice.CreateFileResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	fileID, err := serviceStruct.CreateFile(ctx, req.FileName, req.DataStream, commonmodel.FileType(req.FileType), req.ContentType, req.VoiceDurationTimeSecond)
	if err != nil {
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
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	err = serviceStruct.DeleteFile(ctx, req.FileId)
	if err != nil {
		return nil, err
	}
	return &kitexfileservice.DeleteFileResp{}, nil
}

func (s *KitexFileServiceImpl) GetFile(ctx context.Context, req *kitexfileservice.GetFileReq) (resp *kitexfileservice.GetFileResp, err error) {
	defer func() {
		err = newerror.TranslateError(err).MarshalError()
	}()
	serviceStruct := service.NewFileService(s.snowFlake, s.dbContext, s.fileConfig)
	dataStream, contentType, err := serviceStruct.GetFile(ctx, req.FileId)
	if err != nil {
		return nil, err
	}
	resp = &kitexfileservice.GetFileResp{
		DataStream:  dataStream,
		ContentType: contentType,
	}
	return resp, nil
}
