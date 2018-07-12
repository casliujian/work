package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var testmap = map[int]TestStruct{}
type TestStruct struct {
	Item1 int
	Item2 string
}

func main() {
	//testmap[1] = TestStruct{1, "2"}
	//ts := (testmap[1])
	//ts.Item1 = 3
	//println(testmap[1].Item1)
	//println(testmap[1].Item2)
	//ia := []int32{1,2,3,4}
	//ifl := []float64{1.1,2.2,3.3}
	//var i int32
	//i = 20
	////iab := *((*[]byte)(unsafe.Pointer(&ia)))
	////iab := make([]byte, 12)
	//buf := new(bytes.Buffer)
	//iab := *((*[]byte)(unsafe.Pointer(&ia)))
	//bufW := bufio.NewWriter(buf)
	//binary.Write(bufW, binary.LittleEndian, i)
	//
	//buf2 := new(bytes.Buffer)
	//enc := gob.NewEncoder(buf2)
	//dec := gob.NewDecoder(buf2)
	//enc.Encode(ia)
	////enc.Encode(ia)
	//var ia2 []int32
	//dec.Decode(&ia2)
	//
	//buf3 := new(bytes.Buffer)
	//encf := gob.NewEncoder(buf3)
	////decf := gob.NewDecoder(buf3)
	//encf.Encode(ifl)
	//
	////enc.Encode(ia)
	//
	//buf4 := new(bytes.Buffer)
	//buf4.Write(buf2.Bytes())
	//buf4.Write(buf3.Bytes())
	//decall := gob.NewDecoder(buf4)
	//var all []float64
	//decall.Decode(&all)
	//
	//fmt.Printf("array %d\n", i)
	//fmt.Printf("bytes %v\n", iab)
	//fmt.Printf("buf %x\n", buf.Bytes())
	//fmt.Printf("buf2 %x, and length is %d\n", buf2.Bytes(), len(buf2.Bytes()))
	//fmt.Printf("dec of buf2 %v\n", ia2)
	//fmt.Printf("all data %v\n", all)
	//var foo []int32
	foo := []int32{1,2,3,4}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, foo)
	fmt.Println(buffer.Len())
	fmt.Printf("buffer %x\n", buffer.Bytes())
}
