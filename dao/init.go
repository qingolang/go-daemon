package dao

import (
	"godaemon/model"
	"sync"
)

// lock
var lock *sync.RWMutex

// serviceInfoData
var serviceInfoData []model.ServiceInfo

// Init
func InitServiceInfoData(serviceInfo []model.ServiceInfo) {
	lock = new(sync.RWMutex)
	serviceInfoData = serviceInfo
}
