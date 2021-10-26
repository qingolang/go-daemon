package initialization

import (
	"godaemon/model"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// config 服务配置
var config *model.BaseConfig

// initBaseConfig 初始化服务配置
func initBaseConfig() (err error) {
	config = new(model.BaseConfig)
	fd, err := ioutil.ReadFile("conf/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fd, config)
	return
}

// GetBaseConfig 获取配置信息
func GetBaseConfig() *model.BaseConfig {
	return config
}
