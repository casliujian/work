package main

import "time"

var clientRequests = ClientRequests{}

type ClientRequests = map[int]ClientRequest

type ClientRequest struct {
	//ID int
	ServerIP string
	Subs []SubscribeInfo
	DataChannel DataChannel
}

type SubscribeInfo struct {
	Pipe int
	Count int
	StartFreq float64
	StopFreq float64
}

// 添加客户端数据请求，并返回新添加的客户端请求的记录
func addClientRequest(rid int, serverIP string, subs *[]SubscribeInfo, dataChannel *DataChannel) *ClientRequest {
	//println("adding client request for client", rid)
	cr := ClientRequest{}
	cr.ServerIP = serverIP
	cr.Subs = *subs
	cr.DataChannel = *dataChannel
	clientRequests[rid] = cr
	println("now client", rid, "has", len(cr.Subs), "subscribes")
	return &cr
}

// 添加数据订阅信息
func addSubscribe(rid int, pipe int, count int, startFreq float64, stopFreq float64) {
	println("client", rid, "adding subscribe for pipe", pipe)
	cr := clientRequests[rid]
	//subs := cr.Subs
	cr.Subs = append(cr.Subs, SubscribeInfo{Pipe:pipe, Count:count, StartFreq:startFreq, StopFreq:stopFreq})
	println("now client", rid, "has", len(cr.Subs), "subscribes")
}

// 删除数据订阅信息
func removeSubscribe(rid int, pipe int) {
	cr := clientRequests[rid]
	var index = -1
	for i, v := range cr.Subs {
		if v.Pipe == pipe {
			index = i
			break
		}
	}
	if index != -1 {
		subs := cr.Subs
		cr.Subs = append(subs[:index], subs[index+1:]...)
	} else {
		println("cannot remove subscribe info, because client", rid, "didn't require pipe", pipe)
		return
	}
}

// 修改数据订阅信息
func modifySubscribe(rid int, pipe int, count int, startFreq float64, stopFreq float64) {
	cr := clientRequests[rid]
	var index = -1
	for i, v := range cr.Subs {
		if v.Pipe == pipe {
			index = i
			break
		}
	}
	if index != -1 {
		sub := &cr.Subs[index]
		sub.Count = count
		sub.StartFreq = startFreq
		sub.StopFreq = stopFreq
	} else {
		println("cannot modify subscribe info of client", rid, "for pipe", pipe)
	}
}


// 对每个Pipe，记录该Pipe收到的数据的最大值和最小值
var minmaxData = map[int]MinmaxDataRecord{}
var playbackMinmaxData = map[int]MinmaxDataRecord{}

type MinmaxDataRecord struct {
	Count int
	//DataX []float32
	MinDataY []float32
	MaxDataY []float32
}

// 更新实时数据最大和最小值，返回更新后的最大和最小值，如果没有更新，则返回nil
func updateMinmaxDataRecord(minmaxData map[int]MinmaxDataRecord, pipe int, count int, dataY []float32) ([]float32, []float32) {
	var minUpdated = false
	var maxUpdated = false

	currentRecord, exists := minmaxData[pipe]
	if exists == false {
		minmaxData[pipe] = MinmaxDataRecord{Count:count, MinDataY:dataY, MaxDataY:dataY}
		return dataY, dataY
	} else {
		currentMinDataY := currentRecord.MinDataY
		currentMaxDataY := currentRecord.MaxDataY
		var i int
		for i=0;i<currentRecord.Count;i++ {
			if dataY[i] < currentMinDataY[i] {
				currentMinDataY[i] = dataY[i]
				minUpdated = true
			}
			if dataY[i] > currentMaxDataY[i] {
				currentMaxDataY[i] = dataY[i]
				maxUpdated = true
			}
		}
		var returnedMin, returnedMax []float32 = nil, nil
		if minUpdated ==true {
			returnedMin = currentMinDataY
		}
		if maxUpdated == true {
			returnedMax = currentMaxDataY
		}
		return returnedMin, returnedMax
	}
}

// 发给客户端的实时数据包，json格式
type ResponseData struct {
	MsgType int `json:"msgType"`
	Data DataInfo `json:"data"`
	Warning string `json:"warning"`
	PlaybackIndex int `json:"playbackIndex"`
}

type DataInfo struct {
	TimeStamp string `json:"timeStamp"`
	DataX []float32 `json:"dataX"`
	DataY []float32 `json:"dataY"`
}

// 发给客户端的数据回放信息，json格式
type Playback struct {
	MsgType int	`json:"msgType"`
	Success int 	`json:"success"`
	ServerIP string	`json:"serverIP"`
	Pipe int 	`json:"pipe"`
	StartDate string	`json:"startDate"`
	StopDate string 	`json:"stopDate"`
	Data int			`json:"data"`
}

// 记录数据回放的信息
var playbackRecords = map[int][]PlaybackRecord{}
type PlaybackRecord struct {
	Pipe int
	Count int
	PlayPeriod int
	StartIndex int
	StartFreq float64
	StopFreq float64
	Sending bool
}

func preparePlayback(cid int, pipe int, count int, playPeriod int, startIndex int, startFreq float64, stopFreq float64) {
	prExists := false
	prs, exists := playbackRecords[cid]
	if exists == false {
		pr := PlaybackRecord{Pipe:pipe, Count: count, PlayPeriod:playPeriod, StartIndex: startIndex, StartFreq:startFreq, StopFreq:stopFreq, Sending:false}
		playbackRecords[cid] = []PlaybackRecord{pr}
		prExists = true
	} else {
		for _, pr := range prs {
			if pr.Pipe == pipe {
				prExists = true
				pr.Count = count
				pr.PlayPeriod = playPeriod
				pr.StartIndex = startIndex
				pr.StartFreq = startFreq
				pr.StopFreq = stopFreq
			}
		}
		if prExists == false {
			pr := PlaybackRecord{Pipe:pipe, Count:count, PlayPeriod:playPeriod, StartIndex: startIndex, StartFreq: startFreq, StopFreq: stopFreq, Sending:false}
			playbackRecords[cid] = append(prs, pr)
		}
	}
}

func sendPlayback(msgType uint32, cid int, pipe int) {
	var removedPrIndex = -1
	prs := playbackRecords[cid]
	for i, pr := range prs {
		if pr.Pipe == pipe && pr.Sending == true {
			//pr.PlayPeriod = playPeriod
			//pr.StartIndex = startIndex
			//pr.Sending = true

			//startNano := time.Now().UnixNano()
			seekSuccess := spuSeek(pr.StartIndex)
			if seekSuccess == 0 {
				println("cannot seek index", pr.StartIndex, "in playback")
				break
			}
			var hasNext = spuNext()
			for hasNext == 1 {
				pr.StartIndex = pr.StartIndex + 1
				datab := spuFetch()
				if msgTypeRecieved := msgTypeOfMsg(&datab); msgTypeRecieved == msgType {
					sendPlaybackData(msgTypeRecieved, pipeOfMsg(&datab), timeOfMsg(&datab), lowFreqOfMsg(&datab), highFreqOfMsg(&datab), dataOfMsg(&datab), pr.StartIndex)
				}
				//stopNano := time.Now().UnixNano()
				time.Sleep(time.Millisecond*time.Duration(pr.PlayPeriod))
			}
			// 如果一个回放推送结束，则删除对应的记录
			removedPrIndex = i
		}
	}
	if removedPrIndex != -1 {
		playbackRecords[cid] = append(prs[0:removedPrIndex], prs[removedPrIndex+1:]...)
	}
}

func startPlayback(cid int, pipe int) {
	prs, _ := playbackRecords[cid]
	for _, pr := range prs {
		if pr.Pipe == pipe {
			pr.Sending = true
		}
	}

}

func pausePlayback(cid int, pipe int) {
	prs, _ := playbackRecords[cid]
	for _, pr := range prs {
		if pr.Pipe == pipe {
			pr.Sending = false
		}
	}
}

func setPlaybackStartIndex(cid int, pipe int, startIndex int) {
	prs, _ := playbackRecords[cid]
	for _, pr := range prs {
		if pr.Pipe == pipe {
			pr.StartIndex = startIndex
		}
	}
}

func setPlaybackPeriod(cid int, pipe int, playPeriod int) {
	prs, _ := playbackRecords[cid]
	for _, pr := range prs {
		if pr.Pipe == pipe {
			pr.PlayPeriod = playPeriod
		}
	}
}

func setPlaybackFreq(cid int, pipe int, startFreq float64, stopFreq float64) {
	prs, _ := playbackRecords[cid]
	for _, pr := range prs {
		if pr.Pipe == pipe {
			pr.StartFreq = startFreq
			pr.StopFreq = stopFreq
		}
	}
}


// NanoMSG客户端收到数据后放到通道内，等待websocket服务端发送
type DataChannel = chan []byte

func closeClientConnection(rid int) {
	dc := clientRequests[rid].DataChannel
	delete(clientRequests, rid)
	delete(playbackRecords, rid)
	close(dc)
	//delete(dataChannels, rid)
}