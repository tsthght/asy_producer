package main

import "C"

import (
	"github.com/tsthght/syncer/config"
	"github.com/tsthght/syncer/mafka"
	"github.com/tsthght/syncer/message"
)

var p *mafka.AsyProducer = nil
var abcdef int64 = 1024

//export InitProducerOnce
func InitProducerOnce(fn *C.char) *C.char {
	if p != nil {
		return C.CString("")
	}
	cfg := config.NewProducerConfig()
	err := cfg.Parse(C.GoString(fn))
	if err != nil {
		return C.CString(err.Error())
	}

	p, err = mafka.NewAsyProducer(cfg)
	if err != nil {
		return C.CString(err.Error())
	}
	go p.Run()
	return C.CString("")
}

//export AsyncMessage
func AsyncMessage (msg *C.char, t C.long) *C.char {
	if p == nil {
		return C.CString("")
	}

	m := message.Message{C.GoString(msg), int64(t)}
	return C.CString(p.Async(m))
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(p.LastApplyTimestamp)
}

//export CloseProducer
func CloseProducer() {
	p.Close()
}

func main() {}