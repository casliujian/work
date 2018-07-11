package main

import (
	// "strings"
	// "encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"fmt"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type RequestJSON struct {
	PipeNum int `json:"pipeNum"` 
	DataItemNum int `json:"dataItemNum"`
}
func serveWs(ws *websocket.Conn) {
	jsonData := &RequestJSON{}
	ws.ReadJSON(jsonData)
	// _, msgb, _ := ws.ReadMessage()
	// println(msgb)
	// json.Unmarshal(msgb, jsonData)
	println(jsonData.DataItemNum)
	fmt.Printf("client request pipe %d for %d data items\n", jsonData.PipeNum, jsonData.DataItemNum)

	// Make a new channel to fetch data to be sent for each client
	dataChan := make(chan float32)


	time.Sleep(4*time.Second)
}

func serveHttp(w http.ResponseWriter, r *http.Request) {
	println("a http connection established")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("error when handling http request")
		return
	}
	defer ws.Close()
	println("a websocket connection established")
	serveWs(ws)
}

func listen(ip string, port int) {
	wssAddr := ip+":"+strconv.Itoa(port)
	http.HandleFunc("/", serveHttp)
	http.ListenAndServe(wssAddr, nil)
}

func main() {
	listen("localhost", 1999)
}
