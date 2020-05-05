package main

//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	"fmt"
	"os"
	"time"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func main() {
	cfg := &log.Config{Level: "debug", File: log.FileLogConfig{Filename:"./testmedium.log"}, DisableTimestamp: true}
	lg, _globalP, _ := log.InitLogger(cfg)
	lg = lg.WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	log.ReplaceGlobals(lg, _globalP)

	ret := C.InitProducerOnce(C.CString("../mafka/mafka.toml"))
	if len(C.GoString(ret)) > 0 {
		fmt.Printf("ret: %s\n", C.GoString(ret))
		os.Exit(1)
	}

	defer C.CloseProducer()

	time.Sleep(5 * time.Second)

	go func() {
		for i:=1;i > 0; i++ {
			fmt.Printf("##### %d \n", int64(C.GetLatestApplyTime()))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for i:=0; i< 10; i++ {
		ret := C.AsyncMessage(C.CString("hello golang"), C.long(i))
		if len(C.GoString(ret)) > 0 {
			fmt.Printf("async error: %s\n", C.GoString(ret))
		}
		time.Sleep(1 * time.Second)
	}
}


