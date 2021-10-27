package core

import (
	"time"

	"godaemon/initialization"
	"godaemon/routers"

	"godaemon/logger"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

// RunServer
func RunServer() {
	router := routers.Routers()
	address := initialization.GetBaseConfig().SERVER_PORT
	s := initServer(address, router)

	logger.Log().Info("server run success on " + address)
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
