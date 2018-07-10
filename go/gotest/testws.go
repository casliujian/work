package main

import (
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
	pipeNum int32
	dataItemNum int32
}
func serveWs(ws *websocket.Conn) {
	jsonData := RequestJSON{}
	//ws.ReadJSON(&jsonData)
	//println(jsonData.pipeNum)

	//_, jd, _ := ws.NextReader()
	//buf := make([]byte, 200)
	//jd.Read(buf)
	//println("received :"+string(buf))
	//
	//
	//mt, p, err := ws.ReadMessage()
	//if err != nil && err != io.EOF {
	//	println("error when serving websocket")
	//	return
	//}
	//println(p)
	//ws.WriteMessage(mt, p)
	ws.ReadJSON(&jsonData)
	fmt.Printf("client request pipe %d for %d data items\n", jsonData.pipeNum, jsonData.dataItemNum)
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
