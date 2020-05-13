package main

import "C"

import (
	"github.com/tsthght/syncer/config"
	"github.com/tsthght/syncer/mafka"
	"github.com/tsthght/syncer/message"
)

var p *mafka.AsyProducer = nil

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
func AsyncMessage (db, tb, sql *C.char, cts, ats, seq C.long) *C.char {
	if p == nil {
		return C.CString("")
	}

	m := message.Message{C.GoString(db), C.GoString(tb), C.GoString(sql), int64(cts), int64(ats), int64(seq)}
	return C.CString(p.Async(m))
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(p.LastApplyTimestamp)
}

//export CloseProducer
func CloseProducer() {
	if p != nil {
		p.Close()
	}
}

//export GetLatestSuccessTime
func GetLatestSuccessTime() C.long {
	return C.long(p.LastSuccessTime)
}

//export GetWaitThreshold
func GetWaitThreshold() C.long {
	return C.long(p.GetWaitThreshold())
}

func main() {}