# README

### Environments (As far as I know):
* `Go` >= 1.9.2
* `Python` >= 3.6.5

### Package Required:
* Go packages:
    - `nanomsg.org/go-mangos`
    - `githum.com/gorilla/websocket`
    - `github.com/json-iterator/go`
* Python packages:
    - `asyncio`
    - `websockets`

### How to Install Packages:
* Go: 
    
    `go get -u <packagename>`
* Python: 
    
    `pip3 install <packagename>`

### Folders:
* `go`: Go source code
    - `datasource.go`: NanoMSG server.
    - `wsserver.go`: The main program, including a NanoMSG client a websocket server.
* `py`: Python source code
    - `wsclient.py`: The websocket client, for testing purpose only.
    - `test.py`: The main testing program.
    - `datasource.py`: A TCP server, which serves as one of the data source of the main program. This file is also for testing purpose only.


### How to RUN:

1. `go run datasource.go -url=<URL>`
2. `go run wsserver.go -nano=<true/false> -url=<URL>`
3. `python3 test.py <IP> <PORT> <CLIENTNUMBERS>`

### Where is the Result:

The result is output to the terminal which runs `wsserver.go`.

### Data Format Sent to the Websocket Client:

Each reply from the websocket server to the client has the following format.

```json
{
    "pipe": <INT>,
    "count": <INT>,
    "data": <FLOATARRAY>
}
```
In the above format, `pipe` is the number of pipe; `count` is the number of data sent to the client; `data` is the actural data sent to the client. For instance, the following json item shows a typical reply.
```json
{
    "pipe": 1,
    "count": 2000,
    "data": [f1, f2, ..., f2000]
}
```