package logger

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *Logger

type Logger struct {
	Core *zap.Logger
}

// 模块相关信息
type ModuleInfo struct {
	ModuleName string
	Airlines   string
	TailNo     string
	Version    string
}

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
func Init(logInfo *LogInfo) error {

	logger = new(Logger)

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
	logger.Core = zap.New(core).WithOptions(zap.AddCaller())

	return nil
}

// Panic 极端错误
func (ll *Logger) Panic(format string, v ...interface{}) {
	ll.Core.Panic(fmt.Sprintf(format, v...))
}

// Error 错误
func (ll *Logger) Error(format string, v ...interface{}) {
	ll.Core.Error(fmt.Sprintf(format, v...))
}

// Warning 警告
func (ll *Logger) Warning(format string, v ...interface{}) {
	ll.Core.Warn(fmt.Sprintf(format, v...))
}

// Info 信息
func (ll *Logger) Info(format string, v ...interface{}) {
	ll.Core.Info(fmt.Sprintf(format, v...))
}

// Debug 校验
func (ll *Logger) Debug(format string, v ...interface{}) {
	ll.Core.Debug(fmt.Sprintf(format, v...))
}

// Log 返回日志对象
func Log() *Logger {
	return logger
}
