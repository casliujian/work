package main
/*
#cgo CFLAGS: -I/home/ofserver1/ppjs/include
//-L/home/ofserver1/ppjs/lib -Wl,-R/home/ofserver1/ppjs/lib -lspuapi -lcurl
#cgo LDFLAGS: -L/home/ofserver1/ppjs/lib -lspuapi -lcurl
// -Wl,-R/home/ofserver1/ppjs/lib -lspuapi -lcurl

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include "spuapi.h"

SpuHandle h = NULL;
char buf[1024*1024];

void initSpu() {
	h = spuNew();
}

void closeSpu() {
	spuFree(h);
}

int cgo_spuQueryAll(int pipe, long btime, long etime) {
	return spuQueryAll(h, pipe, btime, etime);
}

int cgo_spuQuerySam(int pipe, long btime, long etime) {
	return spuQuerySam(h, pipe, btime, etime);
}

//int cgo_spuQueryAlm(long btime, long etime) {
//	return spuQueryAlm(h, btime, etime);
//}

int cgo_spuCount() {
	return spuCount(h);
}

int cgo_spuFirst() {
	return spuFirst(h);
}

int cgo_spuLast() {
	return spuLast(h);
}

int cgo_spuNext() {
	return spuNext(h);
}

int cgo_spuPrev() {
	return spuPrev(h);
}

int cgo_spuSeek(int index) {
	return spuSeek(h, index);
}

int cgo_spuFetch() {
	return spuFetch(h, buf, sizeof(buf));
}

//int x = 1;
char y[4] = "abcd";
//char* y = (char*)malloc(4);
//y[0] = 'a';
//y[1] = 'b';
//y[2] = 'c';
//y[3] = 'd';
//void test() {
//	printf("helloworld %d\n", x);
//}
//
//void addx() {
//	x = x + 1;
//}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

func main() {

	gcy := C.GoString(C.y)
	yb := *(*[]byte)(unsafe.Pointer(&gcy))
	fmt.Printf("%x\n", yb)

	C.initSpu()

	querySamSuccess := C.cgo_spuQuerySam(0, 0, 2000)
	if querySamSuccess == 0 {
		println("Read from spuQuerySam fails")
		C.closeSpu()
		return
	}
	count := C.cgo_spuCount()
	println("spuCount:", count)

	var cnt = 0
	var hasNext = C.cgo_spuNext()
	for hasNext == 1 {
		ret := C.cgo_spuFetch()

		if ret == -1 {
			println("spuFetch error")
			//C.closeSpu()
			break
		}
		fmt.Printf("data read from spuFetch %d", ret)
		cnt += 1
		hasNext = C.cgo_spuNext()
	}

	//t := time.Date(2018, 7, 11, 11, 0,0,0,time.)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "1970-01-01 00:01:01", time.UTC)

	println(t.Unix(), time.Now().Day())
	//C.closeSpu()
}
