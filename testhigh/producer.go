package main

//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	"container/list"
	"fmt"
	"os"
	"sync"
	"time"
	"errors"

	"github.com/pingcap/log"
	"github.com/tsthght/syncer/orderlist"
	"go.uber.org/zap"
)

//simulate baseSyncer
type baseSyncer struct {
	success chan *Item
}

func newBaseSyncer() *baseSyncer {
	return &baseSyncer{make(chan *Item, 8)}
}

func(b *baseSyncer) Success() chan *Item {
	return b.success
}

type MafkaSyncer struct {
	toBeAckCommitTSMu      sync.Mutex
	toBeAckCommitTS *orderlist.MapList
	shutdown chan struct{}
	*baseSyncer
}

type Item struct {
	data string
	ts int64
}

func (it *Item) GetKey() int64 {
	return it.ts
}

func NewMafkaSyncer(cfgFile string) (*MafkaSyncer, error){
	if cfgFile == "" {
		return nil, errors.New("config file name is empty")
	}
	ret := C.InitProducerOnce(C.CString(cfgFile))
	if len(C.GoString(ret)) > 0 {
		return nil, errors.New(C.GoString(ret))
	}

	executor := &MafkaSyncer{
		toBeAckCommitTS: orderlist.NewMapList(),
		shutdown: make(chan struct{}),
		baseSyncer: newBaseSyncer(),
	}

	return executor, nil
}

func (ms *MafkaSyncer) Sync(item *Item) error {
	//do sth translator
	C.AsyncMessage(C.CString(string(item.data)), C.long(item.ts))
	ms.toBeAckCommitTSMu.Lock()
	ms.toBeAckCommitTS.Push(item)
	ms.toBeAckCommitTSMu.Unlock()
	return nil
}

func (ms *MafkaSyncer) Close() {
	close(ms.shutdown)
}

func (ms *MafkaSyncer) SetSafeMode(mode bool) bool {
	return false
}

func (ms *MafkaSyncer) Run () {
	var wg sync.WaitGroup

	// handle successes from producer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for ; ; {
			ts := int64(C.GetLatestApplyTime())
			ms.toBeAckCommitTSMu.Lock()
			fmt.Printf("##### before : %d\n", ms.toBeAckCommitTS.Size())
			var next *list.Element
			for elem := ms.toBeAckCommitTS.GetDataList().Front(); elem != nil; elem = next {
				if elem.Value.(orderlist.Keyer).GetKey() <= ts {
					next = elem.Next()
					ms.success <- elem.Value.(*Item)
					ms.toBeAckCommitTS.Remove(elem.Value.(orderlist.Keyer))
				} else {
					break
				}
			}
			fmt.Printf("##### after : %d\n", ms.toBeAckCommitTS.Size())
			ms.toBeAckCommitTSMu.Unlock()

			time.Sleep(5 * time.Second)
		}
		fmt.Printf("################################\n")
		fmt.Printf("run exit\n")
		fmt.Printf("################################\n")
	}()

	for {
		select {
		case it := <- ms.success: {
			fmt.Printf("\n##success## %s, %d \n", it.data, it.ts)
		}
		case <-ms.shutdown:
			//ms.SetErr(nil)
			wg.Wait()
			C.CloseProducer()
			fmt.Printf("################################\n")
			fmt.Printf("main exit\n")
			fmt.Printf("################################\n")
			return
		}
	}
}

func main() {
	cfg := &log.Config{Level: "debug", File: log.FileLogConfig{Filename:"./testhigh.log"}, DisableTimestamp: true}
	lg, _globalP, _ := log.InitLogger(cfg)
	lg = lg.WithOptions(zap.AddStacktrace(zap.DPanicLevel))
	log.ReplaceGlobals(lg, _globalP)

	producer, err := NewMafkaSyncer("../mafka/mafka.toml")
	if err != nil {
		fmt.Printf("NewMafkaSyncer failed: %s\n", err.Error())
		os.Exit(1)
	}
	defer producer.Close()

	//generator data
	go func() {
		time.Sleep(5 * time.Second)
		for i:=1;i > 0; i++ {
			if i == 20 {
				producer.Close()
				break
			}
			it := Item{"hello mafka", int64(i)}
			producer.Sync(&it)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	producer.Run()

}
