package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/json-iterator/go"
)

var intsize = 4
var floatsize = 8
var intnumber = 4
var floatnumber = 65536
var bigendian = false
var keepReceiving = true

// Read data from NanoMSG
func readDataNanoMSG(url string) {
	sock, err := sub.NewSocket()
	if err != nil {
		println("NanoMSG client cannot make a socket")
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		println("NanoMSG client cannot connect to server via socket")
		return
	}
	if err = sock.SetOption(mangos.OptionSubscribe, []byte("")); err != nil {
		println("NanoMSG client cannot subscribe topics")
		return
	}

	for {
		msg, err := sock.Recv()
		if err != nil {
			println("error when receiving pipe number from NanoMSG")
			return
		}
		PipeNums1 := make([]int32, intnumber)
		pipeReader1 := bytes.NewBuffer(msg[0:16])
		binary.Read(pipeReader1, binary.LittleEndian, &PipeNums1)
		//
		Datas1 := make([]float64, floatnumber)
		dataReader1 := bytes.NewBuffer(msg[16:])
		binary.Read(dataReader1, binary.LittleEndian, &Datas1)

		sendPipeData(&(PipeNums1), &(Datas1))

	}

}

// Read data from raw tcp socket
func readDataRAW(url string) {
	addr := url
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		println("connection to data source failed")
		return
	}
	println("connected to " + addr)
	sockround := 0
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
		sendPipeData(&pipeNums, &datas)

		sockround = sockround + 1
	}
	println("end receiving loop")
	finishNano := time.Now().UnixNano()
	fmt.Printf("%f round(s) per second\n", float32(sockround)/float32((finishNano-startNano)/1000000000))
	conn.Close()
}

func sendPipeData(pipeNumsP *[]int32, rawDataP *[]float64) {
	pipeNums := *pipeNumsP
	rawData := *rawDataP
	//println("trying sending from pipe", pipeNums[0])
	for _, ritem := range registry {
		dataRequests, dataChan := ritem.DataRequests, ritem.DataChan
		for _, dr := range dataRequests {
			if dr.PipeNum == pipeNums[0] {
				data := make([]float64, dr.DataNum)
				interval := len(rawData) / int(dr.DataNum)
				var i int
				for i = 0; i < int(dr.DataNum); i++ {
					data[i] = rawData[(i)*interval]
				}
				dataChan <- DataItemInfo{dr.PipeNum, dr.DataNum, &data}
			}
		}
	}
}

func stopReceiving() {
	keepReceiving = false
}

func shuttingDown(timelimit int) {
	if timelimit == 0 {
		return
	}
	time.Sleep(time.Duration(timelimit) * time.Second)
	stopReceiving()
	for i, cr := range record {
		finishTime := time.Now().UnixNano()
		print("Sending to client ", i, " has speed [")
		timeInterval := (float32(finishTime-cr.StartTime) / 1000000000)
		for k,v := range cr.Rounds {
			fmt.Printf("pipe %d : %f rounds/s ",k, float32(v)/timeInterval)
		}
		println("]")

		//fmt.Printf("Client %d received %d rounds in %f seconds, speed %f round(s) per second\n", i, cr.Rounds, (float32(finishTime-cr.StartTime) / 1000000000), (float32(cr.Rounds) / (float32(finishTime-cr.StartTime) / 1000000000)))
	}
	os.Exit(1)
}

var record = map[int]*ClientRecord{}

type ClientRecord struct {
	Rounds    map[int32]int
	StartTime int64
}

func addClientRecordRound(cr *ClientRecord, pipe int32) {
	_, exists := cr.Rounds[pipe]
	if exists == true {
		cr.Rounds[pipe] += 1
	} else {
		cr.Rounds[pipe] = 1
	}
}

var registry = map[int]RegistryItem{}

type DataRequest struct {
	PipeNum int32
	DataNum int32
}

type DataItemInfo struct {
	PipeInfo int32
	CountInfo int32
	Data *[]float64
}

type RegistryItem struct {
	DataRequests []DataRequest
	DataChan     chan DataItemInfo
}

var rid = 0

func newRid() int {
	rid = rid + 1
	return rid
}

func registerDataChan(registerId int, dataRequests *[]DataRequest, dataChan chan DataItemInfo) {
	registry[registerId] = RegistryItem{DataRequests: *dataRequests, DataChan: dataChan}
	newClientRecord := new(ClientRecord)
	//newClientRecord.Rounds = 0
	newClientRecord.Rounds = map[int32]int{}
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

type Subscribe struct {
	Pipe      int32 `json:"pipe"`
	Count     int32 `json:"count"`
	StartFreq int32 `json:"startFreq"`
	StopFreq  int32 `json:"stopFreq"`
}

type RequestJSON struct {
	ServerIP  string      `json:"serverIP"`
	Subscribe []Subscribe `json:"subscribe"`
}

type ResponseJSON struct {
	PipeNum int32     `json:"pipe"`
	Count   int32     `json:"count"`
	Data    []float64 `json:"data"`
}

func serveWs(ws *websocket.Conn) {
	request := &RequestJSON{}
	ws.ReadJSON(request)
	//println(jsonData.DataItemNum)
	requestItems := make([]DataRequest, len(request.Subscribe))
	for i, subs := range request.Subscribe {
		requestItems[i] = DataRequest{PipeNum: subs.Pipe, DataNum: subs.Count}
		fmt.Printf("client request pipe %d for %d data items\n", subs.Pipe, subs.Count)
	}
	// Make a new channel to fetch data to be sent for each client
	dataChan := make(chan DataItemInfo, 10)
	registerId := newRid()
	registerDataChan(registerId, &requestItems, dataChan)

	for {
		data := <-dataChan
		var responseData = ResponseJSON{}
		responseData.PipeNum = data.PipeInfo
		responseData.Count = data.CountInfo
		responseData.Data = *(data.Data)

		json := jsoniter.ConfigCompatibleWithStandardLibrary
		datab, err := json.Marshal(&responseData)
		ws.WriteMessage(websocket.BinaryMessage, datab)

		currentRecord := record[registerId]
		addClientRecordRound(currentRecord, data.PipeInfo)
		if err != nil {
			unregisterDataChan(registerId)
			ws.Close()
			break
		}

	}
}

func serveHttp(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("error when handling http request")
		return
	}
	//defer ws.Close()
	go serveWs(ws)
}

func listen(ip string, port int) {
	wssAddr := ip + ":" + strconv.Itoa(port)
	http.HandleFunc("/", serveHttp)
	http.ListenAndServe(wssAddr, nil)
}

func main() {
	usingNano := flag.Bool("nano", true, "Whether using NanoMSG or not")
	url := flag.String("url", "tcp://127.0.0.1:8080", "URL for connecting NanoMSG, in the form of either \"tcp://\" or \"ipc://\"")
	wsport := flag.Int("port", 1999, "Websocket server listening port")
	timeout := flag.Int("timeout", 30, "Timeout for testing purpose (seconds)")
	flag.Parse()

	// Fetch data from raw tcp socket
	//go readDataRAW("192.168.9.72", 2000)

	// Fetch data from NanoMSG
	if *usingNano == true {
		go readDataNanoMSG(*url)
	} else {
		go readDataRAW(*url)
	}

	//go readDataNanoMSG(*nanoUrl)

	// Shut down the system in `timeout` seconds
	go shuttingDown(*timeout)
	listen("", *wsport)

	// Don't remove this, otherwise the program fails
	time.Sleep(10 * time.Second)

}
