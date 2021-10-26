package routers

import (
	"godaemon/api"

	"github.com/gin-gonic/gin"
)

// Routers 创建路由
func Routers() *gin.Engine {
	var Router = gin.Default()
	Router.POST("/serviceInfo/set", api.SetServiceInfo)
	Router.GET("/serviceInfo/get", api.GetServiceInfo)
	Router.GET("/serviceInfo/find", api.FindServiceInfo)
	Router.DELETE("/serviceInfo/del", api.DelServiceInfo)

	Router.POST("/task/set", api.SetTask)
	Router.GET("/task/get", api.GetTask)
	Router.GET("/task/find", api.FindTask)
	Router.DELETE("/task/del", api.DelTask)

	return Router
}
