package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int

const (
	ERROR LogLevel = 0
	WARN  LogLevel = 1
	INFO  LogLevel = 2
	DEBUG LogLevel = 3
)

const logLevel LogLevel = DEBUG

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func Debug(msg string, fields ...zapcore.Field) {
	if logLevel >= DEBUG {
		logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zapcore.Field) {
	if logLevel >= INFO {
		logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zapcore.Field) {
	if logLevel >= WARN {
		logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zapcore.Field) {
	if logLevel >= ERROR {
		logger.Error(msg, fields...)
	}
}
