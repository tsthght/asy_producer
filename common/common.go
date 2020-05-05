package main

import "C"

import (
	"sync"
	"fmt"

	"github.com/tsthght/syncer/mafka"
)

var p *mafka.AsyProducer = nil
var abcdef int64 = 1024

//export InitProducerOnce
func InitProducerOnce(fn *C.char) *C.char {
	if p != nil {
		return C.CString("")
	}
	var err error
	p, err = mafka.NewAsyProducer(C.GoString(fn))
	if err != nil {
		return C.CString(err.Error())
	}
	go p.Run()
	return C.CString("")
}

//export AsyncMessage
func AsyncMessage (msg *C.char, t C.long) {
	if p == nil {
		return
	}

	m := Message{C.GoString(msg), int64(t)}
	p.Async(m)
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(abcdef)
}

func main() {}

type Message struct {
	Msg string `json:"message"`
	ApplyTime int64 `json:"timestamp"`
}