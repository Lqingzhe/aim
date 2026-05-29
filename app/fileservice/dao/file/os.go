package file

import (
	newerror "aim/pkg/error"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func setOs(info *FileInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("os:SetOs")
	err = os.MkdirAll(filepath.Dir(info.StoragePath), 0755)
	if err != nil {
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeFileOperationFail, "Set File Failed", fmt.Errorf("Os Set File Error"), newerror.LevelError)
	}
	return os.WriteFile(info.StoragePath, info.DataStream, 0644)
}
func getOs(info *FileInfo) (exist bool, err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("os:GetOs")
	date, err := os.ReadFile(info.StoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, newerror.MakeError(http.StatusInternalServerError, newerror.CodeFileOperationFail, "Get File Failed", fmt.Errorf("Os Get File Error"), newerror.LevelError)
	}
	a := info.FileModel //隔离，防止修改原数据
	a.DataStream = date
	info.Info = append(info.Info, &a)
	return true, nil
}
func deleteOs(info *FileInfo) (err error) {
	defer func(trace string) {
		err = newerror.TranslateError(err).AddErrorTrace(trace)
	}("os:DeleteOs")
	err = os.Remove(info.StoragePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件已不存在，也算成功
			return nil
		}
		return newerror.MakeError(http.StatusInternalServerError, newerror.CodeFileOperationFail, "Delete File Failed", fmt.Errorf("Os Delete File Error"), newerror.LevelError)
	}
	return nil
}
