package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
	"time"
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

	var bar int32 = 1
	println(unsafe.Alignof(bar))
	println(unsafe.Sizeof(bar))

	fa := make([]float64, 65536)
	for i,_ := range fa {
		fa[i] = float64(i)+0.00001
	}
	//fa := []float64{1.1,2.2}

	//fab := make([]byte, 2*8)
	var start1 = time.Now().UnixNano()
	var fab bytes.Buffer
	//fabb := bytes.NewBuffer(fab)
	binary.Write(&fab, binary.LittleEndian, fa)
	fmt.Printf("fab:\t %x\n", fab.Bytes())
	var stop1 = time.Now().UnixNano()




	var start2 = time.Now().UnixNano()
	var fabig = make([]float64, len(fa)*8)
	//fabig := []float64{3.3,4.4}
	//fabig[0:2] = fa[:]
	//fa = append(fa[:], fabig[:]...)
	copy(fabig, fa)
	//var float64align = unsafe.Sizeof(fa[0])
	//var fa0off = unsafe.Offsetof(fa[0])
	//var fap = (unsafe.Pointer(&fa))
	fab2 := (*[]byte)(unsafe.Pointer(&fabig))
	fmt.Printf("fab2:\t %x\n", (*fab2)[0:len(fa)*8])
	var stop2 = time.Now().UnixNano()

	println("two conversion are the same:", len(fab.Bytes()) ==len(*fab2))
	println("using binary.write use time\t", stop1-start1)
	println("using unsafe.pointer use time\t", stop2-start2)

}
