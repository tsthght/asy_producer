package main

//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	"fmt"
	"time"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	dsync "github.com/pingcap/tidb-binlog/drainer/sync"
)

func main() {
	cfg := &log.Config{Level: "debug", File: log.FileLogConfig{Filename:"./test.log"}, DisableTimestamp: true}
	lg, _globalP, _ := log.InitLogger(cfg)
	lg = lg.WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	log.ReplaceGlobals(lg, _globalP)

	ret := C.InitProducerOnce(C.CString("../mafka/mafka.toml"))
	fmt.Printf("ret: %s\n", C.GoString(ret))

	C.ConsumeMessage()

	for i:=0; i< 10; i++ {
		C.AsyncMessage(C.CString("hello golang"), C.long(i))
		time.Sleep(1 * time.Second)
	}

	it := dsync.Item{}
	fmt.Printf("%v\n", it)
}


