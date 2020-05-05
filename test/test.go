package main

//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	"fmt"
	"github.com/pingcap/log"
)

func main() {
	for i:=0; i< 10; i++ {
		C.AsyncMessage(C.CString("hello golang"), 21)
		p := int64(C.GetLatestApplyTime())
		log.Info("test")
		fmt.Printf("## %d\n", p)
	}
}


