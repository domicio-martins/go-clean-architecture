package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var ZapLogger *zap.Logger

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log_level"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.TimeKey = "timestamp_app"

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig:    encoderConfig,
	}

	var err error
	ZapLogger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = ZapLogger.Sugar()
	defer logger.Sync()
}

type LogContext map[string]interface{}

func Debug(msg string, ctx *LogContext) {
	logger.Debugw(msg, zap.Any("event", ctx))
}

func Info(msg string, ctx *LogContext) {
	logger.Infow(msg, zap.Any("event", ctx))
}

func Warn(msg string, ctx *LogContext) {
	logger.Warnw(msg, zap.Any("event", ctx))
}

func Error(msg string, err error, ctx *LogContext) {
	logger.Errorw(msg, zap.Error(err), zap.Any("event", ctx))
}

func Fatal(msg string, err error, ctx *LogContext) {
	logger.Fatalw(msg, zap.Error(err), zap.Any("event", ctx))
}
