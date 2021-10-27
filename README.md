- [简介](#简介)
- [配置](#配置)
  - [config.yaml](#configyaml)
  - [services.json](#servicesjson)
- [API](#api)
  - [服务配置](#服务配置)
    - [获取服务配置列表](#获取服务配置列表)
    - [获取单个服务配置](#获取单个服务配置)
    - [更新或者增加服务配置](#更新或者增加服务配置)
    - [删除单个服务配置](#删除单个服务配置)
  - [任务配置](#任务配置)
    - [获取任务列表](#获取任务列表)
    - [获取单个任务](#获取单个任务)
    - [更新或者增加任务](#更新或者增加任务)
    - [删除任务](#删除任务)
- [运行](#运行)

### 简介
微服务架构前期 ,服务高可用解决方案。

- 服务名称：

![点击图片，进入可视化编辑](./doc/static/daemon.png)

### 配置

#### config.yaml
位置： /conf/config.yaml
- SERVER 服务名称
- PUBLIC_CONF 公共配置文件地址
- SERVER_PORT 服务端口
#### services.json
位置： /conf/services.json
- name 服务名称
- priority 优先级
- isInit 是否属于初始化服务 （如果为true 优先级高于所有 为false的服务）
- isDaemon 是否守护进程（如果 false 仅air_nanny 服务启动时拉起一次）
- script 脚本
	- execPath 脚本执行路径
	- programFilePath 执行程序路径 程序自动拼接 公共工作路径（公共配置中的rootpath）
	- startCommand 启动指令
	- stopCommand 停止指令
- healthCheck 保活验证
	-	isHealthCheck 是否启用此验证
	-	interval  运行状况检查将在容器启动后首先运行interval秒，然后在每次之前的检查完成后再次运行interval秒
	-	start_period 为需要时间引导的容器提供初始化时间。在此期间的探测失败将不计入最大重试次数。但是，如果在启动期间健康检查成功，则认为容器已启动，所有连续失败都将计入最大重试次数。
	-	retries 重试次数
### API

#### 服务配置

##### 获取服务配置列表
```http
GET /serviceInfo/get HTTP/1.1
{
  "data": [
    {
      "name": "godbproxy",
      "priority": 3,
      "isInit": false,
      "isDaemon": true,
      "script": {
        "execPath": "../godbproxy/run.sh",
        "programFilePath": "/godbproxy/godbproxy",
        "startCommand": "start",
        "stopCommand": "stop"
      },
      "healthCheck": {
        "isHealthCheck": true,
		....
    }
  ],
  "msg": "SUCCESS"
}
```
##### 获取单个服务配置
```http
GET /serviceInfo/find?name=godbproxy HTTP/1.1
{
  "data": {
    "name": "godbproxy",
    "priority": 3,
    "isInit": false,
    "isDaemon": true,
    "script": {
      "execPath": "../godbproxy/run.sh",
      "programFilePath": "/godbproxy/godbproxy",
      "startCommand": "start",
      "stopCommand": "stop"
    },
    "healthCheck": {
      "isHealthCheck": true,
     ....
    }
  },
  "msg": "SUCCESS"
}
```
##### 更新或者增加服务配置
```http
POST /serviceInfo/set HTTP/1.1
content-type: application/json

{
        "name":"godbproxy" ,
        "priority":3 ,
        "isInit":false,
        "isDaemon":true,
        "script":{
            "execPath":"../../机载端数据处理/godbproxy/run.sh" ,
            "programFilePath":"/mnt/e/project/机载端数据处理/godbproxy/godbproxy" , 
            "startCommand":"start" , 
            "stopCommand":"stop"
        },
        "healthCheck":{
            "isHealthCheck":true ,
        ...
        }
}
```
##### 删除单个服务配置
```http
DELETE /serviceInfo/del?name=godbproxy HTTP/1.1
```

#### 任务配置
##### 获取任务列表
```http
GET /task/get HTTP/1.1
{
  "data": [
    {
      "ServiceInfo": {
        "name": "godbproxy",
        "priority": 3,
        "isInit": false,
        "isDaemon": true,
        "script": {
          "execPath": "../godbproxy/run.sh",
          "programFilePath": "/godbproxy/godbproxy",
          "startCommand": "start",
          "stopCommand": "stop"
        },
        "healthCheck": {
          "isHealthCheck": true,
         ...
          }
        }
      },
      "State": 2,
      "StartFaultNum": 0,
      "StopFaultNum": 0,
      "Pids": [
        771
      ]
    }
  ],
  "msg": "SUCCESS"
}
```

##### 获取单个任务
```http
GET /task/find?name=godbproxy HTTP/1.1
{
  "data": {
    "ServiceInfo": {
      "name": "godbproxy",
      "priority": 3,
      "isInit": false,
      "isDaemon": true,
      "script": {
        "execPath": "../godbproxy/run.sh",
        "programFilePath": "/godbproxy/godbproxy",
        "startCommand": "start",
        "stopCommand": "stop"
      },
      "healthCheck": {
        "isHealthCheck": true,
        ...
      }
    },
    "State": 2,
    "StartFaultNum": 0,
    "StopFaultNum": 0,
    "Pids": [
      771
    ]
  },
  "msg": "SUCCESS"
}
```
##### 更新或者增加任务
```
POST /task/set HTTP/1.1
content-type: application/json

{
        "name":"godbproxy" ,
        "priority":3 ,
        "isInit":false,
        "isDaemon":true,
        "script":{
            "execPath":"../../机载端数据处理/godbproxy/run.sh" ,
            "programFilePath":"/mnt/e/project/机载端数据处理/godbproxy/godbproxy" , 
            "startCommand":"start" , 
            "stopCommand":"stop"
        },
        "healthCheck":{
            "isHealthCheck":true ,
          ...
        }
}
```
##### 删除任务
```http
DELETE /task/del?name=godbproxy HTTP/1.1
```

### 运行

  1. git clone https://github.com/qingolang/godaemon.git
  2. cd godaemon
  3. 修改 conf/config.yaml
  4. 完善 services.json 要守护的服务
  5. 项目根目录下运行以下指令
  
  ```shell
    go mod tidy
    go build 
    ./run.sh start 
  ```