package main

import (
	"nanomsg.org/go-mangos/protocol/sub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"encoding/binary"
	"nanomsg.org/go-mangos"
	"bytes"
	"time"
	"fmt"
	"os"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"flag"
	"github.com/json-iterator/go"
)

var intsize = 4
var floatsize = 8
var intnumber = 4
var floatnumber = 65536
var bigendian = false
var keepReceiving = true



// Read data from NanoMSG
/*
消息格式：
	消息类型	uint		4字节
	消息类型	uint		4字节
	通道号		uint		4字节
	时间		long		8字节
	最低频率	double		8字节
	最高频率	double		8字节
	数据区		[]floag32	4*字节

msgType为115时检索宽带频谱数据
msgType为116时检索高分辨频谱数据
*/
func readDataNanoMSG(url string, msgType uint32) {
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
	msgTypeb := make([]byte, 4)
	binary.LittleEndian.PutUint32(msgTypeb, msgType)
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

		//println("ws server received from NanoMSG length:", len(msg))
		var msgTypeReceived = msgTypeOfMsg(&msg)
		var pipe = pipeOfMsg(&msg)
		var time = timeOfMsg(&msg)
		var lowFreq = lowFreqOfMsg(&msg)
		var highFreq = highFreqOfMsg(&msg)
		var data = dataOfMsg(&msg)

		fmt.Printf("Data received from NanoMSG: msgType: %d, pipe %d, time %d, lowFreq %f, highFreq %f, data %d\n", msgTypeReceived, pipe, time, lowFreq, highFreq, len(*data))

		if msgTypeReceived := msgTypeOfMsg(&msg); msgTypeReceived == msgType {
			sendPipeData(msgTypeReceived, pipeOfMsg(&msg), timeOfMsg(&msg), lowFreqOfMsg(&msg), highFreqOfMsg(&msg), dataOfMsg(&msg))
		}

		//var msgTypeReceived uint32 = 0
		//var pipe uint32 = 0
		//var time uint64 = 0
		//var lowFreq float64 = 0
		//var highFreq float64 = 0
		//var data = make([]float32, len(msg[36:])/4)
		//
		//binary.Read(bytes.NewBuffer(msg[0:4]), binary.LittleEndian, &msgTypeReceived)
		//if msgTypeReceived == msgType {
		//	binary.Read(bytes.NewBuffer(msg[8:12]), binary.LittleEndian, &pipe)
		//	binary.Read(bytes.NewBuffer(msg[12:20]), binary.LittleEndian, &time)
		//	binary.Read(bytes.NewBuffer(msg[20:28]), binary.LittleEndian, &lowFreq)
		//	binary.Read(bytes.NewBuffer(msg[28:36]), binary.LittleEndian, &highFreq)
		//	binary.Read(bytes.NewBuffer(msg[36:]), binary.LittleEndian, &data)
		//
		//	//fmt.Printf("Data received from NanoMSG: msgType: %d, pipe %d, time %d, lowFreq %f, highFreq %f, data %d\n", msgTypeReceived, pipe, time, lowFreq, highFreq, len(data))
		//
		//	sendPipeData(msgTypeReceived, pipe, time, lowFreq, highFreq, &data)
		//}
	}

}

func msgTypeOfMsg(msgp *[]byte) uint32 {
	msg := *msgp
	var msgTypeReceived uint32 = 0
	binary.Read(bytes.NewBuffer(msg[0:4]), binary.LittleEndian, &msgTypeReceived)
	return msgTypeReceived
}

func pipeOfMsg(msgp *[]byte) uint32 {
	msg := *msgp
	var pipe uint32 = 0
	binary.Read(bytes.NewBuffer(msg[8:12]), binary.LittleEndian, &pipe)
	return pipe
}

func timeOfMsg(msgp *[]byte) uint64 {
	msg := *msgp
	var time uint64 = 0
	binary.Read(bytes.NewBuffer(msg[12:20]), binary.LittleEndian, &time)
	return time
}

func lowFreqOfMsg(msgp *[]byte) float64 {
	msg := *msgp
	var lowFreq float64 = 0
	binary.Read(bytes.NewBuffer(msg[20:28]), binary.LittleEndian, &lowFreq)
	return lowFreq
}

func highFreqOfMsg(msgp *[]byte) float64 {
	msg := *msgp
	var highFreq float64 = 0
	binary.Read(bytes.NewBuffer(msg[28:36]), binary.LittleEndian, &highFreq)
	return highFreq
}

func dataOfMsg(msgp *[]byte) *[]float32 {
	msg := *msgp
	var data = make([]float32, len(msg[36:])/4)
	binary.Read(bytes.NewBuffer(msg[36:]), binary.LittleEndian, &data)
	return &data
}

// 发送三份数据：实时数据(msgType: 11s)、历史最低数据(msgType: 12)、历史最高数据(msgType: 13)，每份数据格式一样
func sendPipeData(msgType uint32, pipe uint32, time uint64, lowFreq float64, highFreq float64, data *[]float32) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for rid, cr := range clientRequests {
		for _, sub := range cr.Subs {
			//println(sub.Pipe, uint32(sub.Pipe))
			if pipe == uint32(sub.Pipe) {
				// println("prepare to add a data to channel of client", rid)
				// 根据请求的startFreq和stopFreq以及通道的startFreq和stopFreq来计算返回的数据的长度

				dataLen := len(*data)
				dataIndexStart := int((sub.StartFreq - lowFreq)/(highFreq-lowFreq)*float64(dataLen))
				dataIndexStop := int((sub.StopFreq - lowFreq)/(highFreq-lowFreq)*float64(dataLen))

				var dataCount = dataIndexStop - dataIndexStart

				var dataX, dataY []float32
				var actualDataCount int
				if dataCount < sub.Count {
					// 请求的数据量大于实际的数据量
					actualDataCount = dataCount
					dataX = make([]float32, dataCount)
					dataY = make([]float32, dataCount)
					var i int
					for i=0; i<dataCount; i++ {
						dataX[i] = float32(sub.StartFreq) + float32(float32(i)/float32(dataLen)) * float32(highFreq - lowFreq)
						//dataY[i] =
					}
					dataY = (*data)[dataIndexStart:dataIndexStop]
				} else {
					// 请求的数据量小于实际的数据量，此时需要取样
					actualDataCount = sub.Count
					dataX = make([]float32, sub.Count)
					dataY = make([]float32, sub.Count)
					interval := dataCount/sub.Count
					var i int
					for i=0; i<sub.Count; i++ {
						dataX[i] = float32(sub.StartFreq) + float32(float32(i*interval)/float32(dataLen)) * float32(highFreq - lowFreq)
						dataY[i] = (*data)[dataIndexStart+i*interval]
						//dataY[i] = float32(lowFreq) + float32(i)*(float32(highFreq - lowFreq))/float32(dataLen)
					}
				}

				dataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:dataY}
				responseData := ResponseData{MsgType:11, Data:dataXY, Warning:"", PlaybackIndex:0}
				responseDatab, _ := json.Marshal(&responseData)
				select {
				case clientRequests[rid].DataChannel <- responseDatab:
				default:
				}

				minDataY, maxDataY := updateMinmaxDataRecord(minmaxData, int(pipe), actualDataCount, dataY)
				if minDataY != nil {
					minDataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:minDataY}
					minResponseData := ResponseData{MsgType:12, Data:minDataXY, Warning:""}
					minResponseDatab,_ := json.Marshal(&minResponseData)
					select {
					case clientRequests[rid].DataChannel <- minResponseDatab:
					default:
					}

				}
				if maxDataY != nil {
					maxDataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:maxDataY}
					maxResponseData := ResponseData{MsgType:13, Data:maxDataXY, Warning:"", PlaybackIndex:0}
					maxResponseDatab,_ := json.Marshal(&maxResponseData)
					select {
					case clientRequests[rid].DataChannel <- maxResponseDatab:
					default:
					}

				}

				// println("added a data to channel of client", rid)
			}
		}
	}
}

// 回放：发送三份数据：实时数据(msgType: 11s)、历史最低数据(msgType: 12)、历史最高数据(msgType: 13)，每份数据格式一样
func sendPlaybackData(msgType uint32, pipe uint32, time uint64, lowFreq float64, highFreq float64, data *[]float32, playbackIndex int) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for rid, prs := range playbackRecords {
		for _, pr := range prs {
			//println(sub.Pipe, uint32(sub.Pipe))
			if pipe == uint32(pr.Pipe) {
				// println("prepare to add a data to channel of client", rid)
				// 根据请求的startFreq和stopFreq以及通道的startFreq和stopFreq来计算返回的数据的长度

				dataLen := len(*data)
				dataIndexStart := int((pr.StartFreq - lowFreq)/(highFreq-lowFreq)*float64(dataLen))
				dataIndexStop := int((pr.StopFreq - lowFreq)/(highFreq-lowFreq)*float64(dataLen))

				var dataCount = dataIndexStop - dataIndexStart

				var dataX, dataY []float32
				var actualDataCount int
				if dataCount < pr.Count {
					// 请求的数据量大于实际的数据量
					actualDataCount = dataCount
					dataX = make([]float32, dataCount)
					dataY = make([]float32, dataCount)
					var i int
					for i=0; i<dataCount; i++ {
						dataX[i] = float32(pr.StartFreq) + float32(i/dataLen) * float32(highFreq - lowFreq)
						//dataY[i] =
					}
					dataY = (*data)[dataIndexStart:dataIndexStop]
				} else {
					// 请求的数据量小于实际的数据量，此时需要取样
					actualDataCount = pr.Count
					dataX = make([]float32, pr.Count)
					dataY = make([]float32, pr.Count)
					interval := dataCount/pr.Count
					var i int
					for i=0; i<pr.Count; i++ {
						dataX[i] = float32(pr.StartFreq) + float32(i*interval/dataLen) * float32(highFreq - lowFreq)
						dataY[i] = float32(lowFreq) + float32(i)*(float32(highFreq - lowFreq))/float32(dataLen)
					}
				}

				dataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:dataY}
				responseData := ResponseData{MsgType:11, Data:dataXY, Warning:"", PlaybackIndex:0}
				responseDatab, _ := json.Marshal(&responseData)
				select {
				case clientRequests[rid].DataChannel <- responseDatab:
				default:
				}

				minDataY, maxDataY := updateMinmaxDataRecord(playbackMinmaxData, int(pipe), actualDataCount, dataY)
				if minDataY != nil {
					minDataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:minDataY}
					minResponseData := ResponseData{MsgType:12, Data:minDataXY, Warning:""}
					minResponseDatab,_ := json.Marshal(&minResponseData)
					select {
					case clientRequests[rid].DataChannel <- minResponseDatab:
					default:
					}

				}
				if maxDataY != nil {
					maxDataXY := DataInfo{TimeStamp:"", DataX:dataX, DataY:maxDataY}
					maxResponseData := ResponseData{MsgType:13, Data:maxDataXY, Warning:"", PlaybackIndex:0}
					maxResponseDatab,_ := json.Marshal(&maxResponseData)
					select {
					case clientRequests[rid].DataChannel <- maxResponseDatab:
					default:
					}

				}

				// println("added a data to channel of client", rid)
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

// 在本地生成客户端的ID
var rid = 0
func newRid() int {
	rid = rid + 1
	return rid
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sending2client = false
func response2clients(ws *websocket.Conn, rid int) {

	for {
		//println("trying sending data to client", rid)
		//dataChannel :=
		cr, exists := clientRequests[rid]
		if exists == true {
			datab := <-cr.DataChannel
			//datab, _ := json.Marshal(&data)
			err := ws.WriteMessage(websocket.BinaryMessage, datab)
			println("sent data to client", rid)
			if err != nil {
				//println("cannot send data to client")
				closeClientConnection(rid)
			}
		} else {
			return
		}

	}
}

func serveWs(ws *websocket.Conn) {
	var currentRid int
	currentRid = newRid()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for {
		// 更改websocket接收形式，先接收[]byte，再解析成map
		_, buf, err := ws.ReadMessage()
		if err != nil {
			println("cannot read message from websocket client")
			ws.Close()
			return
		}
		m := bytes2map(buf)
		fmt.Printf("ws server received msg:%s\n", string(buf))
		switch interface2int(m["msgType"]) {
		// 收到订阅消息请求
		case 1:
			serverIP := m["serverIP"].(string)
			if serverIP == localIP {
				subs := m["subscribe"].([]interface{})
				println("received", len(subs), "requests")

				//_ = addClientRequest(currentRid, serverIP)
				//var dataChannel = make(DataChannel, 10)
				//addDataChannel(currentRid, dataChannel)
				var pipe, count int
				var startFreq, stopFreq float64
				var subsInfo = make([]SubscribeInfo, len(subs))
				for i, sub := range subs {
					subMap := sub.(map[string]interface{})

					pipe = int(subMap["pipe"].(float64))
					count = int(subMap["count"].(float64))
					startFreq = subMap["startFreq"].(float64)
					stopFreq = subMap["stopFreq"].(float64)
					//addSubscribe(currentRid, pipe, count, startFreq, stopFreq)
					subsInfo[i] = SubscribeInfo{Pipe:pipe, Count:count, StartFreq:startFreq, StopFreq:stopFreq}
				}
				var dataChannel = make(DataChannel, 10)
				_ = addClientRequest(currentRid, serverIP, &subsInfo, &dataChannel)
				//if sending2client == false {
				go response2clients(ws, currentRid)
					//sending2client = true
				//}
			}
		// 修改订阅请求
		case 2:
			serverIP := m["serverIP"].(string)
			if serverIP == localIP {
				pipe := int(m["pipe"].(float64))
				count := int(m["count"].(float64))
				startFreq := m["startFreq"].(float64)
				stopFreq := m["stopFreq"].(float64)
				modifySubscribe(currentRid, pipe, count, startFreq, stopFreq)
			}
			// 取消订阅
		case 3:
			serverIP := m["serverIP"].(string)
			if serverIP == localIP {
				pipe := int(m["pipe"].(float64))
				removeSubscribe(currentRid, pipe)
			}
			// 收到数据回放请求
		case 4:
			state := int(m["state"].(float64))
			pipe := int(m["pipe"].(float64))
			startDate := m["startDate"].(string)
			stopDate := m["stopDate"].(string)
			playPeriod := int(m["playPeriod"].(float64))
			startIndex := int(m["startIndex"].(float64))
			startFreq := m["startFreq"].(float64)
			stopFreq := m["stopFreq"].(float64)
			println("data playback request from client", currentRid, state, pipe, startDate, stopDate, playPeriod, startIndex)
			switch state {
			case 0:
				startTime,_ := time.ParseInLocation("2006-01-02 15:04:05", startDate, time.Local)
				stopTime,_ := time.ParseInLocation("2006-01-02 15:04:05", stopDate, time.Local)

				var querySuccess int
				if startTime.Day() == time.Now().Day() {
					querySuccess = spuQueryAll(pipe, startTime.Unix(), stopTime.Unix())
				} else {
					querySuccess = spuQuerySam(pipe, startTime.Unix(), stopTime.Unix())
				}
				var playBackResponse Playback
				if querySuccess == 1 {
					spuCount := spuCount()
					playBackResponse = Playback{MsgType:5, Success:1, ServerIP:localIP, Pipe:pipe, StartDate:startDate, StopDate:stopDate, Data:spuCount}
				} else {
					playBackResponse = Playback{MsgType:5, Success:0, ServerIP:localIP, Pipe:pipe, StartDate:startDate, StopDate:stopDate, Data:0}
				}
				playBackResponseb,_ := json.Marshal(&playBackResponse)
				select {
				case clientRequests[currentRid].DataChannel <- playBackResponseb:
					preparePlayback(currentRid, pipe, 600, playPeriod, startIndex, startFreq, stopFreq)
				default:
				}

			case 1:
				startPlayback(currentRid, pipe)
				go sendPlayback(115, currentRid, pipe)
			case 2:
				pausePlayback(currentRid, pipe)
			case 3:
				setPlaybackPeriod(currentRid, pipe, playPeriod)
			case 4:
				setPlaybackStartIndex(currentRid, pipe, startIndex)
			case 5:
				setPlaybackFreq(currentRid, pipe, startFreq, stopFreq)
			//default:

			}

			println("data playback, have not implemented")
			// 不明请求
		default:
			println("unknown request type")
		}
	}
}

func serveHttp(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("error when handling http request ", err.Error())

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
	//usingNano := flag.Bool("nano", true, "Whether using NanoMSG or not")
	url := flag.String("url", "tcp://127.0.0.1:8080", "URL for connecting NanoMSG, in the form of either \"tcp://\" or \"ipc://\"")
	wsport := flag.Int("port", 1999, "Websocket server listening port")
	timeout := flag.Int("timeout", 30, "Timeout for testing purpose (seconds)")
	flag.Parse()

	// Fetch data from raw tcp socket
	//go readDataRAW("192.168.9.72", 2000)

	// 得到本地IP地址，用来过滤掉非法请求：一个客户端的请求是合法的仅当请求中serverIP字段与本地IP一致
	localIP = getLocalIP()

	initSpu()

	// Fetch data from NanoMSG
	go readDataNanoMSG(*url, 115)


	//go readDataNanoMSG(*nanoUrl)

	// Shut down the system in `timeout` seconds
	go shuttingDown(*timeout)
	listen("", *wsport)
	closeSpu()
	// Don't remove this, otherwise the program fails
	time.Sleep(10 * time.Second)

}
