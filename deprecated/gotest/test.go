package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

var testmap = map[int]TestStruct{}
type TestStruct struct {
	Item1 int
	Item2 string
}


func main() {
	foo := []int32{1,2,3,4}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, foo)
	fmt.Println(buffer.Len())
	fmt.Printf("buffer %x\n", buffer.Bytes())

	da := []float64{0.1,0.2,0.3,0.4}
	dab1 := new(bytes.Buffer)
	binary.Write(dab1, binary.LittleEndian, da)
	fmt.Printf("dab1 %x\n", dab1.Bytes())

	var dab2 =make([]byte,1000,1000)

	dab2 = *((*[]byte)(unsafe.Pointer(&(da))))
	dab3 := *((*[]byte)(unsafe.Pointer(&(da[1:]))))
	dab4 := *((*[]byte)(unsafe.Pointer(&(da[2:]))))
	dab5 := *((*[]byte)(unsafe.Pointer(&(da[3:]))))
	fmt.Printf("dab2 %x\n", dab2)
	fmt.Printf("dab3 %x\n", dab3)
	fmt.Printf("dab4 %x\n", dab4)
	fmt.Printf("dab5 %x\n", dab5)

	//fmt.Printf("%x\n", [8]byte(3.2))

}
