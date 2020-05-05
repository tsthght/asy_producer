package main

import "C"

import (
	"fmt"
	"os"

	"github.com/tsthght/syncer/mafka"
)

var p *mafka.AsyProducer = nil
var abcdef int64 = 1024

//export InitProducerOnce
func InitProducerOnce(fn *C.char) {
	if p != nil {
		return
	}
	var err error
	p, err = mafka.NewAsyProducer("./mafka/mafka.toml")
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(1)
	}
	cfg := p.GetProducerConfig()
	fmt.Printf("%v\n", cfg.Topic)
}

//export AsyncMessage
func AsyncMessage (msg *C.char, t C.long) {
	abcdef ++
	fmt.Printf("%s\n", msg)
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(abcdef)
}

func main() {}