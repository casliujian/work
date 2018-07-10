package main

import (
	// "strings"
	"encoding/json"
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
	pipeNum string `json:"pipeNum"` 
	dataItemNum string `json:"dataItemNum"`
}
func serveWs(ws *websocket.Conn) {
	jsonData := &RequestJSON{"3","4"}
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
	// ws.ReadJSON(&jsonData)
	_, msgb, _ := ws.ReadMessage()
	println(msgb)
	msg := `{"pipeNum": "1", "dataItemNum": "2000"}`
	// dec := json.NewDecoder(strings.NewReader(msg))
	// dec.Decode(&jsonData)
	json.Unmarshal([]byte(msg), jsonData)
	println(jsonData.dataItemNum)
	fmt.Printf("client request pipe %s for %s data items\n", jsonData.pipeNum, jsonData.dataItemNum)
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
	// mt, msg, _ := ws.ReadMessage()
	// fmt.Printf("receved %s", msg)
	// ws.WriteMessage(mt, msg)

}

func listen(ip string, port int) {
	wssAddr := ip+":"+strconv.Itoa(port)
	http.HandleFunc("/", serveHttp)
	http.ListenAndServe(wssAddr, nil)
}

func main() {
	listen("localhost", 1999)
}
