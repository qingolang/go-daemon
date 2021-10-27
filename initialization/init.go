package initialization

import "godaemon/logger"

// Init 初始化数据
func Init() (err error) {

	// 加载数据
	if err = initServiceInfoData(); err != nil {
		logger.Log().Panic(err.Error())
		return err
	}

	// 加载服务配置
	if err = initBaseConfig(); err != nil {
		logger.Log().Panic(err.Error())
		return
	}

	// 重新初始化日志
	initLogger()

	// 启动Init任务
	upInitTask()

	return
}
