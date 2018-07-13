package main

import (
	"nanomsg.org/go-mangos/protocol/pub"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"
	"flag"
	"bytes"
	"encoding/binary"
)

func server(url string) {
	sock, err := pub.NewSocket()
	if err != nil {
		println("server cannot make a socket")
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	err = sock.Listen(url)
	if err != nil {
		println("server cannot listen on the socket")
		return
	}
	pipe1 := make([]int32, 4)
	data1 := make([]float64, 65536)
	pipe2 := make([]int32, 4)
	data2 := make([]float64, 65536)
	pipe3 := make([]int32, 4)
	data3 := make([]float64, 65536)
	pipe4 := make([]int32, 4)
	data4 := make([]float64, 65536)

	//for i, _ := range pipe1 {
	//	pipe1[i] = 1
	//	pipe2[i] = 2
	//	pipe3[i] = 3
	//	pipe4[i] = 4
	//}

	pipe1 = []int32{1,1,1,1}
	pipe2 = []int32{2,2,2,2}
	pipe3 = []int32{3,3,3,3}
	pipe4 = []int32{4,4,4,4}

	for i, _ := range data1 {
		data1[i] = float64(i)+0.000001
		data2[i] = float64(i)+0.000002
		data3[i] = float64(i)+0.000003
		data4[i] = float64(i)+0.000004
	}



	for {
		//println("data source sending data")
		var tmpPipeBuf1, tmpPipeBuf2, tmpPipeBuf3, tmpPipeBuf4 bytes.Buffer
		var tmpDataBuf1, tmpDataBuf2, tmpDataBuf3, tmpDataBuf4 bytes.Buffer
		//encPipe1 := gob.NewEncoder(&tmpPipeBuf1)
		//encPipe2 := gob.NewEncoder(&tmpPipeBuf2)
		//encPipe3 := gob.NewEncoder(&tmpPipeBuf3)
		//encPipe4 := gob.NewEncoder(&tmpPipeBuf4)
		//
		//encData1 := gob.NewEncoder(&tmpDataBuf1)
		//encData2 := gob.NewEncoder(&tmpDataBuf2)
		//encData3 := gob.NewEncoder(&tmpDataBuf3)
		//encData4 := gob.NewEncoder(&tmpDataBuf4)
		//
		//encPipe1.Encode(pipe1)
		//encPipe2.Encode(pipe2)
		//encPipe3.Encode(pipe3)
		//encPipe4.Encode(pipe4)
		//
		//encData1.Encode(data1)
		//encData2.Encode(data2)
		//encData3.Encode(data3)
		//encData4.Encode(data4)

		binary.Write(&tmpPipeBuf1, binary.LittleEndian, pipe1)
		binary.Write(&tmpPipeBuf2, binary.LittleEndian, pipe2)
		binary.Write(&tmpPipeBuf3, binary.LittleEndian, pipe3)
		binary.Write(&tmpPipeBuf4, binary.LittleEndian, pipe4)

		binary.Write(&tmpDataBuf1, binary.LittleEndian, data1)
		binary.Write(&tmpDataBuf2, binary.LittleEndian, data2)
		binary.Write(&tmpDataBuf3, binary.LittleEndian, data3)
		binary.Write(&tmpDataBuf4, binary.LittleEndian, data4)




		//pipeBytes1 := *((*[]byte)(unsafe.Pointer(&pipe1)))
		//println("length of pipe1", len(pipe1), ", and length of pipeBytes1", len(pipeBytes1))
		//pipeBytes2 := *((*[]byte)(unsafe.Pointer(&pipe1)))
		//pipeBytes3 := *((*[]byte)(unsafe.Pointer(&pipe1)))
		//pipeBytes4 := *((*[]byte)(unsafe.Pointer(&pipe1)))
		//
		//dataBytes1 := *((*[]byte)(unsafe.Pointer(&data1)))
		//dataBytes2 := *((*[]byte)(unsafe.Pointer(&data2)))
		//dataBytes3 := *((*[]byte)(unsafe.Pointer(&data3)))
		//dataBytes4 := *((*[]byte)(unsafe.Pointer(&data4)))


		var buf1, buf2, buf3, buf4 bytes.Buffer
		buf1.Write(tmpPipeBuf1.Bytes())
		buf1.Write(tmpDataBuf1.Bytes())
		buf2.Write(tmpPipeBuf2.Bytes())
		buf2.Write(tmpDataBuf2.Bytes())
		buf3.Write(tmpPipeBuf3.Bytes())
		buf3.Write(tmpDataBuf3.Bytes())
		buf4.Write(tmpPipeBuf4.Bytes())
		buf4.Write(tmpDataBuf4.Bytes())

		//fmt.Printf("pipebuf1 %x\n", (tmpPipeBuf1.Bytes()))
		//p1 := make([]int32, 4)
		//binary.Read(bytes.NewBuffer(tmpPipeBuf1.Bytes()), binary.LittleEndian, p1)
		//fmt.Printf("pipe1 %v\n", p1)
		//println("size of databuf1", len(tmpPipeBuf1.Bytes()))
		//println("length of buf1", len(buf1.Bytes()))

		// Sending every pipe data
		sock.Send(buf1.Bytes())
		//fmt.Printf("buf1.bytes() %x\n", buf1.Bytes())
		sock.Send(buf2.Bytes())
		sock.Send(buf3.Bytes())
		sock.Send(buf4.Bytes())


	}

}

func main() {
	url := flag.String("url", "tcp://127.0.0.1:8000", "URL of the server")
	flag.Parse()
	println("now starting server")
	server(*url)
}