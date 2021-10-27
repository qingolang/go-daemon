package tasks

import (
	"godaemon/model"

	"godaemon/logger"
)

// UpTask 服务初始化使用
func UpTask(data model.Data) {
	for _, service := range data.ServiceList {
		t := new(task)
		t.ServiceInfo = service
		// 启动
		t.State = 1

		startState, err := t.start()
		if err != nil {
			logger.Log().Error(t.ServiceInfo.Name + " func start  ERR: " + err.Error())
		} else {
			if !startState {
				logger.Log().Error(t.ServiceInfo.Name + " process start fault ")
			}
		}

		// 守护
		if t.ServiceInfo.IsDaemon {
			t.Done = make(chan struct{}, 1)
			go t.daemon()
		}

		taskData = append(taskData, t)
	}
}
