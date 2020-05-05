package main

//#cgo CFLAGS: -I.
//#cgo LDFLAGS: -L. -lcommon
//
//#include "common.h"
import "C"
import "fmt"

func main() {
	for i:=0; i< 10; i++ {
		C.AsyncMessage("", 21)
		p := int64(C.GetLatestApplyTime())
		fmt.Printf("## %d\n", p)
	}
}


