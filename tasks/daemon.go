package tasks

import (
	"fmt"
	"godaemon/logger"
	"godaemon/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// 最大失败次数 -1 表示不检测
	FAULT_MAX_START_NUM = -1
	FAULT_MAX_STOP_NUM  = -1

	// 守护进程检测频率
	DAEMON_DETECTION = 3 * time.Second
)

// 可执行文件根目录
var EXEC_BASE_PATH = "/opt/portal/media/tencent/portal"

type task struct {
	ServiceInfo   model.ServiceInfo
	State         uint          // 0 已停止 1 正在启动 2 运行中 3 停止中 4 服务异常
	StartFaultNum int           // 启动失败次数
	StopFaultNum  int           // 停止失败次数
	Pids          []int         // 进程ID
	Done          chan struct{} `json:"-"`
}

// Start 启动
func (t *task) start() (state bool, err error) {

	exit, err := t.isProcessExist()
	if err != nil {
		return false, err
	}
	if exit {
		return true, err
	}
	logger.Log().Warning(t.ServiceInfo.Name + " start ....")
	cmd := exec.Command(t.ServiceInfo.Script.ExecPath, t.ServiceInfo.Script.StartCommand)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	if err = cmd.Start(); err != nil {
		return
	}
	time.Sleep(25 * time.Second)

	exit, err = t.isProcessExist()
	if err != nil {
		return false, err
	}
	if !exit {
		t.StartFaultNum++
		if FAULT_MAX_START_NUM != -1 {
			if t.StartFaultNum >= FAULT_MAX_START_NUM && t.State != 4 {
				t.Done <- struct{}{}
			}
		}
		return
	}

	t.StartFaultNum = 0
	return true, nil
}

// HealthCheck 发送请求
// 运行状况检查将在容器启动后首先运行interval秒，然后在每次之前的检查完成后再次运行interval秒
// 如果检查的单次运行时间超过timeout秒，则认为检查失败。
// start period为需要时间引导的容器提供初始化时间。在此期间的探测失败将不计入最大重试次数。
// 但是，如果在启动期间健康检查成功，则认为容器已启动，所有连续失败都将计入最大重试次数。
// retries 重试次数
func (t *task) healthCheck() (state bool, err error) {
	if t.ServiceInfo.HealthCheck.CMD == "" || t.ServiceInfo.HealthCheck.Interval == 0 {
		return true, nil
	}
	if !t.ServiceInfo.HealthCheck.IsRe {
		time.Sleep(time.Second * time.Duration(t.ServiceInfo.HealthCheck.StartPeriod))
		t.ServiceInfo.HealthCheck.IsRe = true
	}
	num := t.ServiceInfo.HealthCheck.Retries
	for num > 0 {
		cmd := exec.Command("sh", "-c", t.ServiceInfo.HealthCheck.CMD)
		err = cmd.Run()
		if err == nil {
			return true, nil
		}
		logger.Log().Warning(t.ServiceInfo.HealthCheck.CMD + " ERR : " + err.Error())
		time.Sleep(time.Second * time.Duration(t.ServiceInfo.HealthCheck.Interval))
		num--
	}
	return false, nil
}

// IsProcessExist 检查进程是否存在
func (t *task) isProcessExist() (exist bool, err error) {
	// 等待5S
	time.Sleep(5 * time.Second)
	psCmd := exec.Command("ps", "-C", t.ServiceInfo.Name)
	psData, _ := psCmd.Output()
	t.Pids = []int{}
	psFields := strings.Fields(string(psData))
	for i, v := range psFields {
		if i <= 3 {
			continue
		}
		if v == t.ServiceInfo.Name {
			// 取出执行文件全路径
			programCmd := exec.Command("ls", "-l", "/proc/"+psFields[i-3]+"/exe")
			programData, _ := programCmd.Output()
			programFields := strings.Fields(string(programData))
			if len(programFields) < 10 {
				continue
			}
			if strings.Contains(programFields[10], EXEC_BASE_PATH+t.ServiceInfo.Script.ProgramFilePath) {
				pid, _ := strconv.Atoi(psFields[i-3])
				t.Pids = append(t.Pids, pid)
			}
		}
	}
	if len(t.Pids) != 0 {
		return true, nil
	}
	logger.Log().Warning(t.ServiceInfo.Name + " process not exits !")
	return
}

// Stop  停止
func (t *task) stop() (state bool, err error) {

	logger.Log().Warning(t.ServiceInfo.Name + " stop ....")

	exit, err := t.isProcessExist()
	if err != nil {
		return false, err
	}
	if !exit {
		return true, err
	}

	// 使用脚本杀死进程
	cmd := exec.Command(t.ServiceInfo.Script.ExecPath, t.ServiceInfo.Script.StopCommand)
	err = cmd.Run()
	if err != nil {
		return
	}
	exit, err = t.isProcessExist()
	if err != nil {
		return
	}
	if exit {
		// 使用kill
		if err = t.kill(); err != nil {
			return
		}
		exit, err = t.isProcessExist()
		if err != nil {
			return
		}
		if exit {
			// 使用kill -9
			if err = t.kill9(); err != nil {
				return
			}
			exit, err = t.isProcessExist()
			if err != nil {
				return
			}
			if exit {
				t.StopFaultNum++
				if FAULT_MAX_STOP_NUM != -1 && t.State != 4 {
					if t.StopFaultNum >= FAULT_MAX_STOP_NUM {
						t.Done <- struct{}{}
					}
				}

				return
			}
		}
	}
	t.StopFaultNum = 0
	return true, nil
}

// Kill 杀死进程
func (t *task) kill() (err error) {
	exit, err := t.isProcessExist()
	if err != nil {
		return err
	}
	if !exit {
		return
	}
	tmpPids := []int{}
	for _, v := range t.Pids {
		cmd := exec.Command("kill", fmt.Sprint(v))
		err = cmd.Run()
		if err != nil {
			logger.Log().Error(` exec.Command("kill", fmt.Sprint(v)) err :` + err.Error())
		} else {
			continue
		}
		tmpPids = append(tmpPids, v)
	}
	t.Pids = tmpPids
	return
}

// Kill9 强制杀死进程
func (t *task) kill9() (err error) {
	exit, err := t.isProcessExist()
	if err != nil {
		return err
	}
	if !exit {
		return
	}

	tmpPids := []int{}
	for _, v := range t.Pids {
		cmd := exec.Command("kill", "-9", fmt.Sprint(v))
		err = cmd.Run()
		if err != nil {
			logger.Log().Error(` exec.Command("kill", fmt.Sprint(v)) err :` + err.Error())
		} else {
			continue
		}
		tmpPids = append(tmpPids, v)
	}
	t.Pids = tmpPids
	return
}

// Restart 重新启动
func (t *task) restart() (state bool, err error) {
	state, err = t.stop()
	if err != nil || !state {
		return
	}
	state, err = t.start()
	return
}

// Daemon 守护进程
func (t *task) daemon() {

	for {

		processExist, err := t.isProcessExist()
		if err != nil {
			logger.Log().Error(t.ServiceInfo.Name+" func isProcessExist  ERR: %s", err.Error())
			time.Sleep(DAEMON_DETECTION)
			continue
		}

		// 如果已经启动
		if processExist {
			t.State = 1
			if t.ServiceInfo.HealthCheck.IsHealthCheck {
				apiState, err := t.healthCheck()
				if err != nil {
					logger.Log().Error(t.ServiceInfo.Name+" func  healthCheck  ERR: %s", err.Error())
					t.stopHandle()
					time.Sleep(DAEMON_DETECTION)
					continue
				}

				if apiState {
					t.State = 2
				} else {
					logger.Log().Warning(t.ServiceInfo.Name + " healthCheck 验证失败 ！")
					t.stopHandle()
				}

			}
		} else {
			t.startHandle()
		}

		// 接收异常退出信号
		if len(t.Done) != 0 {
			<-t.Done
			t.State = 4
			close(t.Done)
			break
		}
		time.Sleep(DAEMON_DETECTION)
	}
}

// StopHandle
func (t *task) stopHandle() {
	t.State = 3
	stopState, err := t.stop()
	if err != nil {
		logger.Log().Error(t.ServiceInfo.Name+" func stop  ERR: %s", err.Error())
		return
	}
	if stopState {
		t.State = 0
		return
	}
	logger.Log().Warning(t.ServiceInfo.Name + "   process stop fault  ")
}

// StartHandle
func (t *task) startHandle() {
	t.State = 1
	// 未启动
	startState, err := t.start()
	if err != nil {
		logger.Log().Error(t.ServiceInfo.Name+" func start  ERR: %s", err.Error())
		return
	}
	if startState {
		return
	}
	logger.Log().Warning(t.ServiceInfo.Name + " process start fault ")
}
