package api

import (
	"godaemon/dao"
	"godaemon/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetServiceInfo
func SetServiceInfo(ctx *gin.Context) {
	var info model.ServiceInfo

	// get request
	if err := ctx.BindJSON(&info); err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : " + err.Error()})
		return
	}

	// valid request
	if info.Script.ExecPath == "" ||
		info.Script.ProgramFilePath == "" ||
		info.Script.StartCommand == "" ||
		info.Script.StopCommand == "" ||
		info.Name == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : 参数传递错误"})
		return
	}

	store := new(dao.DaoServiceInfo)
	store.Set(info)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS"})
}

// GetServiceInfo
func GetServiceInfo(ctx *gin.Context) {
	store := new(dao.DaoServiceInfo)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS", "data": store.Get()})
}

// FindServiceInfo
func FindServiceInfo(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok || name == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : 参数传递错误"})
		return
	}
	store := new(dao.DaoServiceInfo)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS", "data": store.Find(name)})
}

// DelServiceInfo
func DelServiceInfo(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok || name == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : 参数传递错误"})
		return
	}
	store := new(dao.DaoServiceInfo)
	store.Del(name)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS"})
}
