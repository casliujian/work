package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	msgType := 4
	var b = bytes.NewBuffer([]byte{})
	binary.Read(b, binary.LittleEndian, &msgType)
	//b := (*[]byte)(unsafe.Pointer(&msgType))
	println(b.Bytes())
	fmt.Printf("b.bytes() is %x\n", b.Bytes())
	b2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(b2, uint32(msgType))
	fmt.Printf("b2 is %x\n", b2)

	var highFreq float64 = 4
	//var highFreqb = make([]byte, 8)
	var highFreqBuf bytes.Buffer
	binary.Write(&highFreqBuf, binary.LittleEndian, highFreq)
	fmt.Printf("highFreqb %x\n", highFreqBuf.Bytes())
}
