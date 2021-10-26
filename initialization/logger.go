package initialization

import (
	"godaemon/logger"

	"go.uber.org/zap/zapcore"
)

// initLogger 初始化日志
func initLogger() {

	// 设置日志
	logInfo := logger.LogInfo{
		Filename: "./godaemon.log",
		//单位MB,
		MaxSize: 100,
		//备份文件数
		//假设配置log文件名:"kafka.log",则生成的备份文件名格式："kafka-2021-01-21T03-40-40.660.log"
		//MaxBackups:2 => 则最多会存在3个log文件，每个大小100MB
		MaxBackups: 2,
		//log保存的天数s
		MaxAge: 100,
		//是否压缩成gzip，备份文件压缩成"kafka-2021-01-21T03-43-48.117.log.gz"
		Compress: false,
		// >= Level的log才会被记录
		Level: zapcore.DebugLevel,
	}
	logger.Init(&logInfo)
}
