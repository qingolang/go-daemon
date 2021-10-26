package api

import (
	"godaemon/model"
	"godaemon/tasks"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetTask
func SetTask(ctx *gin.Context) {
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
	tasks.SetTask(info)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS"})
}

// GetTask
func GetTask(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS", "data": tasks.GetTasks()})
}

// FindTask
func FindTask(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok || name == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : 参数传递错误"})
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS", "data": tasks.FindTasks(name)})
}

// DelTask
func DelTask(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok || name == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ERROR : 参数传递错误"})
		return
	}
	tasks.DelTask(name)
	ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "SUCCESS"})
}
