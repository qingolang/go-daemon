package dao

import (
	"encoding/json"
	"godaemon/model"
	"os"

	"godaemon/logger"
)

// DaoServiceInfo
type DaoServiceInfo struct{}

// Del
func (d *DaoServiceInfo) Del(name string) {
	d.del(name)
	d.writeFile()
}

// del
func (d *DaoServiceInfo) del(name string) {
	lock.Lock()
	defer lock.Unlock()
	var tmpData []model.ServiceInfo
	for _, v := range serviceInfoData {
		if v.Name != name {
			tmpData = append(tmpData, v)
		}
	}
	serviceInfoData = tmpData
}

// Set
func (d *DaoServiceInfo) Set(res model.ServiceInfo) {
	d.set(res)
	d.writeFile()
}

// set
func (d *DaoServiceInfo) set(res model.ServiceInfo) {
	lock.Lock()
	defer lock.Unlock()
	for i, v := range serviceInfoData {
		if v.Name == res.Name {
			serviceInfoData[i] = res
			return
		}
	}
	serviceInfoData = append(serviceInfoData, res)
}

// Find
func (d *DaoServiceInfo) Find(name string) (res model.ServiceInfo) {
	lock.RLock()
	defer lock.RUnlock()
	for _, v := range serviceInfoData {
		if v.Name == name {
			res = v
			break
		}
	}
	return
}

// Get
func (d *DaoServiceInfo) Get() []model.ServiceInfo {
	lock.RLock()
	defer lock.RUnlock()
	return serviceInfoData
}

// writeFile
func (d *DaoServiceInfo) writeFile() {
	f, err := os.OpenFile("conf/services.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Log().Error(` os.OpenFile("conf/services.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644) ERR : ` + err.Error())
		return
	}
	defer f.Close()

	lock.RLock()
	defer lock.RUnlock()
	byteData, err := json.Marshal(serviceInfoData)
	if err != nil {
		logger.Log().Error(` data json.Marshal(data) ERR : ` + err.Error())
		return
	}
	n, _ := f.Seek(0, os.SEEK_END)
	_, err = f.WriteAt(byteData, n)
	if err != nil {
		logger.Log().Error(` data f.WriteAt(byteData, n) ERR : ` + err.Error())
		return
	}
}
