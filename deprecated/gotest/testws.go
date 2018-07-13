package main

import (
	// "strings"
	// "encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"fmt"
	"unsafe"
)

var registry = map[int]registryItem{}
//type registryItems = []registryItem
type registryItem struct {
	PipeNum int32
	DataNum int32
	DataChan chan[]float64
}

var rid = 0
func newRid() int {
	rid = rid + 1
	return rid
}

func registerDataChan(registerId int, pipeNum int32, dataNum int32, dataChan chan []float64) {
	registry[registerId] = registryItem{PipeNum:pipeNum, DataNum:dataNum, DataChan:dataChan}
}

func unregisterDataChan(registerId int) {
	rItem := registry[registerId]
	close(rItem.DataChan)
	delete(registry, registerId)
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type RequestJSON struct {
	PipeNum int32 `json:"pipeNum"`
	DataItemNum int32 `json:"dataItemNum"`
}
func serveWs(ws *websocket.Conn) {
	jsonData := &RequestJSON{}
	ws.ReadJSON(jsonData)
	println(jsonData.DataItemNum)
	fmt.Printf("client request pipe %d for %d data items\n", jsonData.PipeNum, jsonData.DataItemNum)

	// Make a new channel to fetch data to be sent for each client
	dataChan := make(chan []float64)
	registerId := newRid()
	registerDataChan(registerId, jsonData.PipeNum, jsonData.DataItemNum, dataChan)

	for {
		data := <- dataChan
		datab := *((*[]byte)(unsafe.Pointer(&data)))
		err := ws.WriteMessage(websocket.BinaryMessage, datab)
		if err != nil {
			unregisterDataChan(registerId)
			break
		}
	}
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
	println("Listening")
}
