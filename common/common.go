package main

import "C"

import (
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
	return C.CString("")
}

//export AsyncMessage
func AsyncMessage (msg *C.char, t C.long) {
	abcdef ++
	fmt.Printf("%s\n", C.GoString(msg))
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(abcdef)
}

func main() {}