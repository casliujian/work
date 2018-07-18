# 软件说明

### 运行环境 (目前):
* `Go` >= 1.9.2

### 用到的程序包:
* Git (安装 Go 的程序包需要用到)
* Go 程序包:
    - `nanomsg.org/go-mangos`
    - `githum.com/gorilla/websocket`
    - `github.com/json-iterator/go`
    - `github.com/toolkits/net`


### 如何安装程序包:
* Go: 
    
    `go get -u <程序包名>`


### 文件夹:
* `go`: Go 代码
    - `wsserver.go`: 主程序，包括一个NanoMSG客户端和一个websocket服务端。
    - `client.go`, `localip.go`, `message.go`, `influxDB.go`: 服务程序，被主程序调用

    


### 如何运行程序:


1. `go run wsserver.go -url=<URL> -port=<PORT> -timeout=<TIME>`
    - `-url=<URL>`: `URL` 为连接到 NanoMSG 的服务端的地址，可为 `tcp://***` 形式或 `ipc://***` 形式
    - `-port=<PORT>`: `PORT` 为 websocket 服务端绑定的端口号, 默认为 `1999`
    - `-timeout=<TIME>`: `TIME` 为指定的系统运行时间，一般用来测试。单位为秒，比如 `-timeout=60` 代表系统运行 60 秒后停止，默认值为 30。特殊地，当指定 `TIME` 为 0 时，系统一直运行直到人为强制退出。

