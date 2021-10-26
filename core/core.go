package core

import (
	"time"

	"godaemon/initialization"
	"godaemon/logger"
	"godaemon/routers"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

// RunServer 跑起服务
func RunServer() {
	// 取出路由
	router := routers.Routers()
	// 取出注册端口
	address := initialization.GetBaseConfig().SERVER_PORT
	// 初始化服务
	s := initServer(address, router)

	logger.Log().Info("server run success on %s \n", address)
	logger.Log().Info(s.ListenAndServe().Error())
}

// initServer 初始化服务
func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}

// server 服务
type server interface {
	ListenAndServe() error
}
