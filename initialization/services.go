package initialization

import (
	"encoding/json"
	"godaemon/dao"
	"godaemon/model"
	"godaemon/tasks"
	"io/ioutil"
	"os"

	"godaemon/logger"
)

// serviceData
var serviceData []model.Data

// upInitTask
func upInitTask() {
	for _, v := range serviceData {
		if v.IsInit {
			tasks.UpTask(v)
		}
	}
}

// UpTask
func UpTask() {
	for _, v := range serviceData {
		if !v.IsInit {
			tasks.UpTask(v)
		}
	}
}

// initServiceInfoData
func initServiceInfoData() (err error) {
	// 读取服务配置文件
	data, err := readJSONFile()
	if err != nil {
		return err
	}

	// 渲染数据
	var tmpServiceList []model.ServiceInfo
	if err = json.Unmarshal(data, &tmpServiceList); err != nil {
		return err
	}

	// 去重或不合格数据
	var serviceList = make([]model.ServiceInfo, 0, len(tmpServiceList))
	for _, tServiceInfo := range tmpServiceList {
		isExit := false
		for _, serviceInfo := range serviceList {
			if tServiceInfo.Name == serviceInfo.Name ||
				tServiceInfo.Script.ExecPath == "" ||
				tServiceInfo.Script.ProgramFilePath == "" ||
				tServiceInfo.Script.StartCommand == "" ||
				tServiceInfo.Script.StopCommand == "" ||
				tServiceInfo.Name == "" {
				logger.Log().Error("服务配表错误或重复数据")
				isExit = true
				break
			}

		}
		if !isExit {
			serviceList = append(serviceList, tServiceInfo)
		}
	}

	// 排序
	bubbleAsortServiceList(serviceList)

	// 取出Init=true 并重新计算优先级
	maxInitPriority := uint(0)
	for _, info := range serviceList {
		if info.IsInit {
			priorityExit := false
			if maxInitPriority < info.Priority {
				maxInitPriority = info.Priority
			}
			for i, data := range serviceData {
				if info.Priority == data.StrategyPriority {
					serviceData[i].ServiceList = append(serviceData[i].ServiceList, info)
					priorityExit = true
					break
				}
			}
			if !priorityExit {
				var data model.Data
				data.IsInit = info.IsInit
				data.RealPriority = info.Priority
				data.StrategyPriority = info.Priority
				data.ServiceList = append(data.ServiceList, info)
				serviceData = append(serviceData, data)
			}
		}
	}

	// 初始化serviceInfo持久层
	dao.InitServiceInfoData(serviceList)

	// 初始化 Manage Task
	tasks.InitManageTask(len(serviceList))

	// 取出Init=false 并重新计算优先级
	for _, info := range serviceList {
		if !info.IsInit {
			priorityExit := false
			strategyPriority := info.Priority + maxInitPriority

			for i, data := range serviceData {
				if strategyPriority == data.StrategyPriority {
					serviceData[i].ServiceList = append(serviceData[i].ServiceList, info)
					priorityExit = true
					break
				}
			}
			if !priorityExit {
				var data model.Data
				data.IsInit = info.IsInit
				data.RealPriority = info.Priority
				data.StrategyPriority = strategyPriority
				data.ServiceList = append(data.ServiceList, info)
				serviceData = append(serviceData, data)
			}
		}
	}

	// 排序
	bubbleAsortServiceData(serviceData)
	return
}

// readJSONFile
func readJSONFile() (data []byte, err error) {
	f, err := os.Open("conf/services.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err = ioutil.ReadAll(f)
	return
}

// bubbleAsortServiceList
func bubbleAsortServiceList(serviceList []model.ServiceInfo) {
	for i := 0; i < len(serviceList)-1; i++ {
		for j := i + 1; j < len(serviceList); j++ {
			if serviceList[i].Priority > serviceList[j].Priority {
				serviceList[i], serviceList[j] = serviceList[j], serviceList[i]
			}
		}
	}
}

// bubbleAsortServiceData
func bubbleAsortServiceData(serviceData []model.Data) {
	for i := 0; i < len(serviceData)-1; i++ {
		for j := i + 1; j < len(serviceData); j++ {
			if serviceData[i].StrategyPriority > serviceData[j].StrategyPriority {
				serviceData[i], serviceData[j] = serviceData[j], serviceData[i]
			}
		}
	}
}
