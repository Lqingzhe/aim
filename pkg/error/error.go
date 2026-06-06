package newerror

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = zapcore.DebugLevel // -1  调试
	LevelInfo  = zapcore.InfoLevel  //  0  一般信息
	LevelWarn  = zapcore.WarnLevel  //  1  警告
	LevelError = zapcore.ErrorLevel //  2  错误
	LevelFatal = zapcore.FatalLevel //  5  致命（谨慎使用）
)

type Error struct {
	HttpCode    int    `json:"http_code"`
	HttpMessage string `json:"http_message"`

	StatusCode int    `json:"status_code"`
	LogMessage string `json:"log_message"`
	LogLevel   int8   `json:"log_level"`

	IsNeedInterrupt bool `json:"is_need_interrupt"`
}
type Option struct {
	id        uint64
	operation string
	ip        string
}
type operate func(*Error)

func (e *Error) MarshalError() error {
	if e == nil {
		return nil
	}
	return fmt.Errorf(`{"http_code":%d,"http_message":"%s","status_code":%d,"log_level":%d,"log_message":"%s","is_need_interrupt":%t}`, e.HttpCode, e.HttpMessage, e.StatusCode, e.LogLevel, e.LogMessage, e.IsNeedInterrupt)
}
func UnMarshalError(err error) error {
	if err == nil {
		return nil
	}
	errStr := err.Error()

	// 查找 JSON 开始位置
	startIdx := strings.Index(errStr, "{")
	if startIdx == -1 {
		return MakeError(http.StatusInternalServerError, CodeInternalError, "Unmarshal Error Failed 1",
			fmt.Errorf("no JSON found in error: %s", errStr), LevelError)
	}

	jsonPart := errStr[startIdx:]
	err2 := &Error{}
	if err3 := json.Unmarshal([]byte(jsonPart), err2); err3 != nil {
		return MakeError(http.StatusInternalServerError, CodeInternalError, "Unmarshal Error Failed 2", fmt.Errorf("%s : %s", `raw error:`, err.Error()), LevelFatal).AddErrorTrace("error:Unmarshal Error").(*Error)
	}
	return err2
}
func WithContinueError(err *Error) {
	err.IsNeedInterrupt = false
}
func MakeError(httpCode int, statueCode ErrorStatue, httpMessage string, err error, logLevel zapcore.Level, Operates ...operate) *Error {
	newStruct := &Error{
		HttpCode:        httpCode,
		StatusCode:      int(statueCode),
		HttpMessage:     httpMessage,
		LogMessage:      fmt.Sprintf("-> %s", err.Error()),
		LogLevel:        int8(logLevel),
		IsNeedInterrupt: true,
	}
	for _, Operate := range Operates {
		Operate(newStruct)
	}
	return newStruct
}
func MakeKafkaError(statueCode ErrorStatue, err error, logLevel zapcore.Level) *Error {
	return &Error{
		StatusCode: int(statueCode),
		LogMessage: fmt.Sprintf("-> %s", err.Error()),
		LogLevel:   int8(logLevel),
	}
}
func (e *Error) Error() string {
	if e == nil || e.LogMessage == "" {
		return ""
	}
	return e.LogMessage
}
func TranslateError(err error) *Error {
	if err == nil {
		return nil
	}
	var err2 *Error
	err2, ok := errors.AsType[*Error](err)
	if !ok {
		return MakeError(http.StatusInternalServerError, CodeInternalError, "Type Assertion Error", fmt.Errorf("%s%w", `Type Assertion To "*newerror.Error" Error`, err), LevelError).AddErrorTrace("error:TranslateError").(*Error)
	}
	return err2
}
func WhetherInterrupt(err error, finalErr *error) bool {
	e := TranslateError(err)
	if e == nil {
		return false
	}
	*finalErr = e
	return e.IsNeedInterrupt
}
func (e *Error) AddErrorTrace(trace string) error {
	if e == nil {
		return nil
	}
	e.LogMessage = trace + " / " + e.LogMessage
	return e
}
func (o *Option) OptionInfo() (uint64, string, string) {
	return o.id, o.operation, o.ip
}
func IsContextError(err error) (bool, *Error) {
	if errors.Is(err, context.DeadlineExceeded) {
		return true, MakeError(http.StatusGatewayTimeout, CodeNetworkTimeout, "Time Out", err, LevelWarn)
	}
	if errors.Is(err, context.Canceled) {
		return true, MakeError(http.StatusGatewayTimeout, CodeNetworkTimeout, "Time Out", err, LevelWarn)
	}
	return false, nil
}
