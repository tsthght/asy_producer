package main

import (
   "fmt"
   "os"
   "sync"
   "time"

   "github.com/tsthght/syncer/mafka"
)


func main() {
   p, err := mafka.NewAsyProducer("./config/mafka.toml")
   if err != nil {
      fmt.Printf("%v", err.Error())
      os.Exit(1)
   }
   cfg := p.GetProducerConfig()
   fmt.Printf("%v\n", cfg.Topic)

   var wg sync.WaitGroup
   wg.Add(1)
   go func() {
      defer wg.Done()
      for {
         select {
         case msgs := <- p.CallBack.SuccessChan:
            for _, msg := range msgs {
               fmt.Printf("##msg = %s\n", msg)
            }
         }
      }
   }()

   //just for test
   go p.Consumer.ConsumeMessage(&mafka.BasicHandler{1})

   for i := 1000; i > 0 ; i-- {
      m := mafka.MafkaMessage{i}
      p.Async(m)
      time.Sleep(5 * time.Second)
   }
}
