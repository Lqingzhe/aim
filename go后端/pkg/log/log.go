package newlog

import (
	newerror "aim/pkg/error"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(service string, equipID int) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, _ := cfg.Build()
	logger = logger.With(zap.String("service", service), zap.Int("equip_id", equipID)) // 全局使用
	return logger                                                                      //设置默认值
} //注册logger,gateway就传入gateway,调用完成后defer logger.Sync()
func AddLatencyAndTime(zapConfig *zap.Logger, beginTime time.Time) *zap.Logger {
	now := time.Now()
	return zapConfig.With(zap.Duration("latency_time", now.Sub(beginTime)))
}
func AddTraceID(zapConfig *zap.Logger, traceID string) *zap.Logger {
	return zapConfig.With(zap.String("trace_id", traceID))
} //在每个handler开头调用
func AddError[t int | newerror.ErrorStatue](zapConfig *zap.Logger, err error, statueCode t) *zap.Logger {
	return zapConfig.With(zap.Int("status_code", int(statueCode)), zap.Error(err))
} //用来在发生错误时添加错误，在AddInfo前调用
func AddGateWayInfo(zapConfig *zap.Logger, httpStatue int, userID int64, ip string, operation string) *zap.Logger {
	return zapConfig.With(zap.Int("http_code", httpStatue), zap.Int64("user_id", userID), zap.String("ip", ip), zap.String("operation", operation))
} //用于在gateway的handler最后调用
func Log[t int8 | zapcore.Level](logger *zap.Logger, rawLogLevel t, message string) {
	logState := zapcore.Level(rawLogLevel)
	switch logState {
	case zapcore.DebugLevel:
		logger.Debug(message)
	case zapcore.InfoLevel:
		logger.Info(message)
	case zapcore.WarnLevel:
		logger.Warn(message)
	case zapcore.ErrorLevel:
		logger.Error(message)
	case zapcore.FatalLevel:
		logger.Fatal(message)
	default:
		logger.Info(message)
	}
} //日志等级见newerror包

func LogInitInfo(logger *zap.Logger, message string) {
	logger.Info(message)
}
func LogInitWarn(logger *zap.Logger, err error, message string) {
	logger.Warn(message, zap.Error(err))
} //打印初始化阶段的Warn级别日志
func LogInitError(logger *zap.Logger, err error, message string) {
	logger.Error(message, zap.Error(err))
} //打印初始化阶段的Error级别日志
func LogInitFatal(logger *zap.Logger, err error, message string) {
	logger.Fatal(message, zap.Error(err))
} //打印初始化阶段的Fatal级别日志
