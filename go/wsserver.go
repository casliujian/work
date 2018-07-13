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
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"nanomsg.org/go-mangos"
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
		//println("size of msg:", len(msg))

		PipeNums := make([]int32, intnumber)
		//pipeReader := bytes.NewBuffer(msg[0:15])
		//binary.Read(pipeReader, binary.LittleEndian, &PipeNums)
		//
		Datas := make([]float64, floatnumber)
		//dataReader := bytes.NewBuffer(msg[16:])
		//binary.Read(dataReader, binary.LittleEndian, &Datas)

		var pipeNum int32
		var data float64
		var i int
		for i = 0; i<intnumber; i++ {
			binary.Read(bytes.NewBuffer(msg[i*4:(i+1)*4]),binary.LittleEndian,&pipeNum)
			PipeNums[i] = pipeNum
		}
		for i=0; i<floatnumber; i++ {
			binary.Read(bytes.NewBuffer(msg[16+i*8:16+(i+1)*8]),binary.LittleEndian,&data)
			Datas[i] = data
		}


		//fmt.Printf("raw pipenums1 %x\n", msg[0:15])
		//fmt.Printf("pipeinfo %v\n", PipeNums)
		//fmt.Printf("PipeNums[0] %d\n", pipeNum)
		sendPipeData(PipeNums[:], Datas[:])

	}

}

// Read data from raw tcp socket
func readDataRAW(ip string, port int) {
	addr := ip + ":" + strconv.Itoa(port)
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
		sendPipeData(pipeNums, datas)

		sockround = sockround + 1
	}
	println("end receiving loop")
	finishNano := time.Now().UnixNano()
	fmt.Printf("%f round(s) per second\n", float32(sockround)/float32((finishNano-startNano)/1000000000))
	conn.Close()
}

func sendPipeData(pipeNums []int32, rawData []float64) {
	//println("trying sending from pipe", pipeNums[0])
	for _, ritem := range registry {
		dataRequests, dataChan := ritem.DataRequests, ritem.DataChan
		for _, dr := range dataRequests {
			if dr.PipeNum == pipeNums[0] {
				data := make([]float64, dr.DataNum+1)
				interval := len(rawData) / int(dr.DataNum)
				var i int
				for i=1; i<int(dr.DataNum)+1; i++ {
					data[i] = rawData[(i-1)*interval]
				}

				//for i,_ := range data {
				//	data[i+1] = rawData[i*interval]
				//}
				data[0] = float64(dr.PipeNum)
				//fmt.Printf("sending to client %d, pipe %d and data %d ...\n", c,  pipeNums, data[0])
				//dataChan := dr.DataChan
				dataChan <- data
			}
		}

		//
		//if ritem.PipeNum == pipeNums[0] {
		//	data := make([]float64, ritem.DataNum)
		//	interval := len(rawData) / int(ritem.DataNum)
		//	for i, _ := range data {
		//		data[i] = rawData[i*interval]
		//	}
		//	ritem.DataChan <- data
		//}
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
		fmt.Printf("Client %d received %d rounds in %f seconds, speed %f round(s) per second\n", i, cr.Rounds, (float32(finishTime-cr.StartTime) / 1000000000), (float32(cr.Rounds) / (float32(finishTime-cr.StartTime) / 1000000000)))
	}
	os.Exit(1)
}

var record = map[int]*ClientRecord{}

type ClientRecord struct {
	Rounds    int
	StartTime int64
}

var registry = map[int]RegistryItem{}

type DataRequest struct {
	PipeNum int32
	DataNum int32
	//DataChan chan []float64
}

type RegistryItem struct {
	DataRequests []DataRequest
	DataChan chan []float64
}

var rid = 0

func newRid() int {
	rid = rid + 1
	return rid
}

func registerDataChan(registerId int, dataRequests *[]DataRequest, dataChan chan []float64) {
	registry[registerId] = RegistryItem{DataRequests: *dataRequests, DataChan: dataChan}
	newClientRecord := new(ClientRecord)
	newClientRecord.Rounds = 0
	newClientRecord.StartTime = time.Now().UnixNano()
	record[registerId] = newClientRecord
}

func unregisterDataChan(registerId int) {
	rItem := registry[registerId]
	close(rItem.DataChan)
	//for _, v := range rItem.DataRequests {
	//	close(v.DataChan)
	//}
	//close(rItem.DataChan)
	delete(registry, registerId)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Subscribe struct {
	Pipe 		int32 `json:"pipe"`
	Count 		int32 `json:"count"`
	StartFreq 	int32 `json:"startFreq"`
	StopFreq	int32 `json:"stopFreq"`
}

type RequestJSON struct {
	ServerIP string `json:"serverIP"`
	Subscribe []Subscribe `json:"subscribe"`

	//PipeNum     int32 `json:"pipeNum"`
	//DataItemNum int32 `json:"dataItemNum"`
}

type ResponseJSON struct {
	PipeNum int32 `json:"pipe"`
	Count int32 `json:"count"`
	Data []float64 `json:"data"`
}

func serveWs(ws *websocket.Conn) {
	request := &RequestJSON{}
	ws.ReadJSON(request)
	//println(jsonData.DataItemNum)
	requestItems := make([]DataRequest, len(request.Subscribe))
	for i,subs := range request.Subscribe {
		requestItems[i] = DataRequest{PipeNum:subs.Pipe, DataNum:subs.Count}
		fmt.Printf("client request pipe %d for %d data items\n", subs.Pipe, subs.Count)
	}

	//fmt.Printf("client request pipe %d for %d data items\n", request.PipeNum, jsonData.DataItemNum)

	// Make a new channel to fetch data to be sent for each client
	dataChan := make(chan []float64)
	registerId := newRid()
	registerDataChan(registerId, &requestItems, dataChan)

	for {
		//for _, v := range requestItems {
			//dataChan := v.DataChan
			data := <-dataChan
			//datab := *((*[]byte)(unsafe.Pointer(&data)))
			var responseData = ResponseJSON{}
			responseData.PipeNum = int32(data[0])
			responseData.Count = int32(len(data[1:]))
			responseData.Data = data[1:]
			//response,_ := json.Marshal(responseData)


			//dataBytes := make([]byte, len(data)*8)
			//binary.Write(bytes.NewBuffer(dataBytes), binary.LittleEndian, &data)
			err := ws.WriteJSON(responseData)
			//err := ws.WriteMessage(websocket., response)
			//println("writing to websocket to client", responseData)
			currentRecord := record[registerId]
			currentRecord.Rounds = currentRecord.Rounds + 1
			if err != nil {
				unregisterDataChan(registerId)
				break
			}
		//}

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

	//go readDataRAW("localhost", 2000)

	go readDataNanoMSG("ipc://testipc")

	// Shut down the system in `timeout` seconds
	go shuttingDown(*timeout)
	listen("127.0.0.1", *wsport)
	time.Sleep(10*time.Second)
	//stopReceiving()
	//time.Sleep(1*time.Second)
}
