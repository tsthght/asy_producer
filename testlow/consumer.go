package main

import (
	"sync"
	"fmt"
	"os"
	"time"

	"github.com/tsthght/syncer/config"
	"github.com/tsthght/syncer/mafka"
)

func main() {
	cfg := config.NewProducerConfig()
	err := cfg.Parse("../mafka/mafka.toml")
	if err != nil{
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	p, err := mafka.NewAsyConsumer(cfg)
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(1)
	}
	defer p.Close()
	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go p.PrintConsuem()
	wg.Wait()
}
