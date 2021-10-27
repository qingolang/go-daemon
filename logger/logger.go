package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// 文件日志相关信息
type LogInfo struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      zapcore.Level
}

// 日志初始化
func Init(logInfo *LogInfo) {

	// 本地存储日志文件初始化
	hook := lumberjack.Logger{
		Filename:   logInfo.Filename,
		MaxSize:    logInfo.MaxSize, //megabytes
		MaxBackups: logInfo.MaxBackups,
		MaxAge:     logInfo.MaxAge,   //days
		Compress:   logInfo.Compress, //gzip, disabled by default
	}
	fileLogWriter := zapcore.AddSync(&hook)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "ms",
		LevelKey:       "lv",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "ca",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	logEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// >=logInfo.Level的日志才会被记录
	logPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logInfo.Level
	})

	// 生成本地存储文件接口
	fileCore := zapcore.NewCore(logEncoder, fileLogWriter, logPriority)
	var allCore []zapcore.Core
	allCore = append(allCore, fileCore)
	core := zapcore.NewTee(allCore...)
	logger = zap.New(core).WithOptions(zap.AddCaller())
}

// Log 返回日志对象
func Log() *zap.Logger {
	return logger
}
