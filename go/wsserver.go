package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"unsafe"
)

var intsize = 4
var floatsize = 8
var intnumber = 4
var floatnumber = 65536
var bigendian = false
var keepReceiving = true

// Read data from a tcp socket
func readData(ip string, port int) {
	addr := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		println("connection to data source failed")
		return
	}
	println("connected to " + addr)
	round := 0
	startNano := time.Now().UnixNano()
	for keepReceiving {
		// Receiving pipe number
		pipebufLen := intsize * intnumber
		pipebuf := make([]byte, pipebufLen)
		receivedPipebufLen := 0
		for receivedPipebufLen < pipebufLen {
			n, err := conn.Read(pipebuf[receivedPipebufLen:])
			if err != nil && err != io.EOF {
				println("read pipe number error")
				conn.Close()
				return
			}
			receivedPipebufLen = receivedPipebufLen + n
		}
		// Receiving real data
		databufLen := floatsize * floatnumber
		databuf := make([]byte, databufLen)
		receivedDatabufLen := 0
		for receivedDatabufLen < databufLen {
			n, err := conn.Read(databuf[receivedDatabufLen:])
			if err != nil && err != io.EOF {
				println("read data error")
				conn.Close()
				return
			}
			receivedDatabufLen = receivedDatabufLen + n
		}

		// Send the corresponding data to each client
		// The safe way, may be not very efficient, and may be further improved
		pipeNums := make([]int32, intnumber)
		pipeNumReader := bytes.NewBuffer(pipebuf)
		binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)

		datas := make([]float64, floatnumber)
		dataReader := bytes.NewBuffer(databuf)
		binary.Read(dataReader, binary.LittleEndian, &datas)
		sendPipeData(pipeNums, datas)

		round = round + 1
	}
	println("end receiving loop")
	finishNano := time.Now().UnixNano()
	fmt.Printf("%f round(s) per second\n", float32(round)/float32((finishNano-startNano)/1000000000))
	conn.Close()
}

func sendPipeData(pipeNums []int32, rawData []float64) {
	for _, ritem := range registry {
		if ritem.PipeNum == pipeNums[0] {
			data := make([]float64, ritem.DataNum)
			interval := len(rawData) / int(ritem.DataNum)
			for i, _ := range data {
				data[i] = rawData[i*interval]
			}
			ritem.DataChan <- data
		}
	}
}

func stopReceiving() {
	keepReceiving = false
}

func shuttingDown(timelimit int) {
	time.Sleep(time.Duration(timelimit) * time.Second)
	stopReceiving()
	for i, cr := range record {
		finishTime := time.Now().UnixNano()
		fmt.Printf("Client %d received %d rounds in %d seconds, speed %f round(s) per second\n", i, cr.Rounds, timelimit, (float32(cr.Rounds) / (float32(finishTime-cr.StartTime) / 1000000000)))
	}
	os.Exit(1)
}

var record = map[int]*ClientRecord{}

type ClientRecord struct {
	Rounds    int
	StartTime int64
}

var registry = map[int]registryItem{}

type registryItem struct {
	PipeNum  int32
	DataNum  int32
	DataChan chan []float64
}

var rid = 0

func newRid() int {
	rid = rid + 1
	return rid
}

func registerDataChan(registerId int, pipeNum int32, dataNum int32, dataChan chan []float64) {
	registry[registerId] = registryItem{PipeNum: pipeNum, DataNum: dataNum, DataChan: dataChan}
	newClientRecord := new(ClientRecord)
	newClientRecord.Rounds = 0
	newClientRecord.StartTime = time.Now().UnixNano()
	record[registerId] = newClientRecord
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
	PipeNum     int32 `json:"pipeNum"`
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
		data := <-dataChan
		datab := *((*[]byte)(unsafe.Pointer(&data)))
		err := ws.WriteMessage(websocket.BinaryMessage, datab)
		currentRecord := record[registerId]
		currentRecord.Rounds = currentRecord.Rounds + 1
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
	wssAddr := ip + ":" + strconv.Itoa(port)
	http.HandleFunc("/", serveHttp)
	http.ListenAndServe(wssAddr, nil)
}

func main() {
	wsport := flag.Int("port", 1999, "Websocket server listening port")
	timeout := flag.Int("timeout", 30, "Timeout for testing purpose (seconds)")
	flag.Parse()

	go readData("localhost", 2000)
	// Shut down the system in `timeout` seconds
	go shuttingDown(*timeout)
	listen("localhost", *wsport)
	//time.Sleep(10*time.Second)
	//stopReceiving()
	//time.Sleep(1*time.Second)
}
