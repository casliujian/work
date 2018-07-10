package main

import (
	"net"
	"strconv"
	"io"
	"time"
	"fmt"
	"unsafe"
	"github.com/gorilla/websocket"
	"net/http"
)

var intsize = 4
var floatsize = 8
var intnumber = 4
var floatnumber = 65536
var bigendian = false
var keepReceiving = true

func readData(ip string, port int) {
	addr := ip+":"+strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err!= nil {
		println("connection to data source failed")
		return
	}
	//defer conn.Close()
	println("connected to "+addr)
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
		//println("received pipe number")
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
		//println("received real data")

		// TODO: Send the corresponding data to each client
		// Safe way
		//pipeNums := make([]int32, intnumber)
		//pipeNumReader := bytes.NewBuffer(pipebuf)
		//binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)
		//
		//datas := make([]float64, floatnumber)
		//dataReader := bytes.NewBuffer(databuf)
		//binary.Read(dataReader, binary.LittleEndian, &datas)
		//
		//println(pipeNums[0])
		//println(datas[0])

		// Unsafe way
		pipeNumsU := *((*[]int32) (unsafe.Pointer(&pipebuf)))
		datasU := *((*[]float64)(unsafe.Pointer(&databuf)))
		//println(pipeNumsU[0])
		//println(datasU[0])
		sendPipeData(pipeNumsU, datasU)


		//switch {
		//case intsize == 4 && floatsize == 8 && bigendian == false:
		//	pipeNums := make([]int32, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)
		//
		//	datas := make([]float64, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.LittleEndian, &datas)
		//	//handlePipeData(pipeNums, datas)
		//case intsize == 4 && floatsize == 8 && bigendian == true:
		//	pipeNums := make([]int32, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.BigEndian, &pipeNums)
		//
		//	datas := make([]float64, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.BigEndian, &datas)
		//case intsize == 8 && floatsize == 8 && bigendian == false:
		//	pipeNums := make([]int64, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)
		//
		//	datas := make([]float64, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.LittleEndian, &datas)
		//case intsize == 8 && floatsize == 8 && bigendian == true:
		//	pipeNums := make([]int64, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.BigEndian, &pipeNums)
		//
		//	datas := make([]float64, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.BigEndian, &datas)
		//case intsize == 4 && floatsize == 4 && bigendian == false:
		//	pipeNums := make([]int32, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)
		//
		//	datas := make([]float32, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.LittleEndian, &datas)
		//case intsize == 4 && floatsize == 4 && bigendian == true:
		//	pipeNums := make([]int32, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.BigEndian, &pipeNums)
		//
		//	datas := make([]float32, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.BigEndian, &datas)
		//case intsize == 8 && floatsize == 4 && bigendian == false:
		//	pipeNums := make([]int64, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.LittleEndian, &pipeNums)
		//
		//	datas := make([]float32, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.LittleEndian, &datas)
		//case intsize == 8 && floatsize == 4 && bigendian == true:
		//	pipeNums := make([]int64, intnumber)
		//	pipeNumReader := bytes.NewReader(pipebuf)
		//	binary.Read(pipeNumReader, binary.BigEndian, &pipeNums)
		//
		//	datas := make([]float32, floatnumber)
		//	dataReader := bytes.NewReader(databuf)
		//	binary.Read(dataReader, binary.BigEndian, &datas)
		//}

		round = round + 1
	}
	println("end receiving loop")
	finishNano := time.Now().UnixNano()
	fmt.Printf("%f round(s) per second\n", float32(round)/float32((finishNano-startNano)/1000000000))
	conn.Close()
}

func sendPipeData(pipeNums []int32, datas []float64) {

}


func stopReceiving() {
	keepReceiving = false
}

var upgrader = websocket.Upgrader{}
func serveWs(ws *websocket.Conn) {
	jsonData := struct {}{}
	ws.ReadJSON()
}

func serveHttp(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("error when handling http request")
		return
	}
	defer ws.Close()
	serveWs(ws)
}

func listenWsClient(ip string, port int) {
	wssAddr := ip+":"+strconv.Itoa(port)

}



func main() {
	go readData("localhost", 2000)
	time.Sleep(10*time.Second)
	stopReceiving()
	time.Sleep(1*time.Second)
}