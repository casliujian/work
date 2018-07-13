# 软件说明

### 运行环境 (目前):
* `Go` >= 1.9.2
* `Python` >= 3.6.5

### 用到的程序包:
* Go 程序包:
    - `nanomsg.org/go-mangos`
    - `githum.com/gorilla/websocket`
    - `github.com/json-iterator/go`
* Python 程序包:
    - `asyncio`
    - `websockets`

### 如何安装程序包:
* Go: 
    
    `go get -u <程序包名>`
* Python: 
    
    `pip3 install <程序包名>`

### 文件夹:
* `go`: Go 代码
    - `datasource.go`: NanoMSG 服务端。
    - `wsserver.go`: 主程序，包括一个NanoMSG客户端和一个websocket服务端。
* `py`: Python 代码
    - `wsclient.py`: 一个websocket客户端，仅供测试。
    - `test.py`: 测试程序入口。
    - `datasource.py`: 一个TCP服务端，可用来作为主程序的数据来源。
    


### 如何运行程序:

1. `go run datasource.go -url <URL>`
2. `go run wsserver.go -nURL <URL>`
3. `python3 test.py <IP> <PORT> <CLIENTNUMBERS>`

### 怎样查看程序运行结果:

程序的运行结果在运行 `wsserver.go` 文件的终端中输出。

### 发到 websocket 的客户端的数据格式:

Websocket 客户端接收的数据具有如下格式。

```json
{
    "pipe": <INT>,
    "count": <INT>,
    "data": <FLOATARRAY>
}
```
在以上形式中，`pipe` 指的是通道编号；`count` 指的是从通道接收的数字个数；`data` 指的是从通道接收到的实际数字的数组。比如，一个典型的接收数据如下。
```json
{
    "pipe": 1,
    "count": 2000,
    "data": [f1, f2, ..., f2000]
}
```