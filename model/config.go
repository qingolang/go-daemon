package model

// BaseConfig
type BaseConfig struct {
	SERVER      string `yaml:"SERVER"`
	PUBLIC_CONF string `yaml:"PUBLIC_CONF"`
	SERVER_PORT string `yaml:"SERVER_PORT"`
}

// PublicConf
type PublicConf struct {
	Rootpath         string   `json:"rootpath"`
	OtaURL           string   `json:"ota_url"`
	CloudURL         string   `json:"cloud_url"`
	LoggingURL       string   `json:"logging_url"`
	Database         Database `json:"database"`
	ServerAddress    string   `json:"server_address"`
	ServerPort       string   `json:"server_port"`
	DomainName       string   `json:"domain_name"`
	DrmPort          string   `json:"drm_port"`
	VersionName      string   `json:"version_name"`
	Deployment       int      `json:"deployment"`
	CacheAddr        string   `json:"cacheAddr"`
	PtcAddr          string   `json:"ptcAddr"`
	CloudGateWayAddr string   `json:"cloudGateWayAddr"`
	SystemService    string   `json:"systemService"`
}

// Database
type Database struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
}
