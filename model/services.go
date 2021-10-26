package model

// ServiceInfo
type ServiceInfo struct {
	Name        string      `json:"name"`
	Priority    uint        `json:"priority"`
	IsInit      bool        `json:"isInit"`
	IsDaemon    bool        `json:"isDaemon"`
	Script      Script      `json:"script"`
	HealthCheck HealthCheck `json:"healthCheck"`
}

// Script
type Script struct {
	ExecPath        string `json:"execPath"`
	ProgramFilePath string `json:"programFilePath"`
	StartCommand    string `json:"startCommand"`
	StopCommand     string `json:"stopCommand"`
}

// HealthCheck
type HealthCheck struct {
	IsHealthCheck bool   `json:"isHealthCheck"`
	Interval      uint64 `json:"interval"`
	StartPeriod   uint64 `json:"start_period"`
	Retries       uint64 `json:"Retries"`
	CMD           string `json:"cmd"`
	IsRe          bool   `json:"-"`
}

type Data struct {
	StrategyPriority uint // 真实运行以策略优先级为准转 ， isInit = true 复用 RealPriority进行排序 如果 isInit = false StrategyPriority = len(isInit = true) + RealPriority
	IsInit           bool `json:"isInit"`
	RealPriority     uint
	ServiceList      []ServiceInfo
}
