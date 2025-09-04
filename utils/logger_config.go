package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Sugar *zap.SugaredLogger
)

func InitLogger() *zap.Logger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewTee(zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel))
	logger := zap.New(core, zap.AddCaller()).Named("[taskManager]")
	Sugar = logger.Sugar()
	return logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/task_manager.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
