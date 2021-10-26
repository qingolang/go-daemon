package tasks

import (
	"godaemon/model"
	"sync"
	"time"
)

// lock
var lock *sync.RWMutex

// tasks
var taskData []*task

// InitManageTask
func InitManageTask(serviceNum int) {
	lock = new(sync.RWMutex)
	// 额外加20个空间
	taskData = make([]*task, 0, serviceNum+20)
}

// SetTask
func SetTask(s model.ServiceInfo) {
	lock.Lock()
	defer lock.Unlock()
	for i, v := range taskData {
		if v.ServiceInfo.Name == s.Name {
			taskData[i].ServiceInfo = s

			// 重启当前任务
			taskData[i].restart()
			return
		}
	}

	t := new(task)
	t.ServiceInfo = s
	// 启动
	t.State = 1
	// 启动当前任务
	t.start()

	// 守护进程
	if t.ServiceInfo.IsDaemon {
		t.Done = make(chan struct{}, 1)
		go t.daemon()
	}

	taskData = append(taskData, t)
}

// DelTask
func DelTask(serviceName string) {
	lock.Lock()
	defer lock.Unlock()

	var tmpTaskList []*task
	for _, v := range taskData {
		if v.ServiceInfo.Name == serviceName {
			// 关闭守护进程
			if v.ServiceInfo.IsDaemon {
				v.Done <- struct{}{}
			}

			// 停止进程运行
			time.Sleep(DAEMON_DETECTION + 1)
			v.stop()

			continue
		}

		tmpTaskList = append(tmpTaskList, v)
	}
	taskData = tmpTaskList
}

// FindTasks
func FindTasks(serviceName string) *task {
	lock.RLock()
	defer lock.RUnlock()
	for _, v := range taskData {
		if v.ServiceInfo.Name == serviceName {
			return v
		}
	}
	return nil
}

// GetTasks
func GetTasks() []*task {
	lock.RLock()
	defer lock.RUnlock()
	return taskData
}
