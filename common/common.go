package main

import "C"

import (
	"fmt"
	"os"

	"github.com/tsthght/syncer/mafka"
)

var p *mafka.AsyProducer = nil

//export InitProducerOnce
func InitProducerOnce(fn C.string) {
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
func AsyncMessage (msg C.string, t C.int64) {

}

//export GetLatestApplyTime
func GetLatestApplyTime() C.int64 {
	return p.LastApplyTime
}

func main() {}