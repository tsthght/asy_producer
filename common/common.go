package main

import "C"

import (
	"fmt"
	"os"

	"github.com/tsthght/syncer/mafka"
)

var p *mafka.AsyProducer = nil

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
	p.LastApplyTime ++
	fmt.Printf("%s\n", string(msg))
}

//export GetLatestApplyTime
func GetLatestApplyTime() C.long {
	return C.long(p.LastApplyTime)
}

func main() {}