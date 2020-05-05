package main

import (
   "fmt"
   "os"
   "sync"
   "time"

   "github.com/tsthght/syncer/config"
   "github.com/tsthght/syncer/mafka"
   "github.com/tsthght/syncer/message"
)


func main() {
   cfg := config.NewProducerConfig()
   err := cfg.Parse("./mafka/mafka.toml")
   if err != nil{
      fmt.Printf("%s\n", err.Error())
      os.Exit(1)
   }
   p, err := mafka.NewAsyProducer(cfg)
   if err != nil {
      fmt.Printf("%v", err.Error())
      os.Exit(1)
   }

   var wg sync.WaitGroup
   wg.Add(1)
   go func() {
      defer wg.Done()
      for {
         select {
         case msgs := <- p.CallBack.SuccessChan:
            for _, msg := range msgs {
               fmt.Printf("## msg = %s\n", msg)
            }
         }
      }
   }()

   //just for test
   //go p.Consumer.ConsumeMessage(&mafka.BasicHandler{1})

   for i := 10; i > 0 ; i-- {
      m := message.Message{"demo", int64(i)}
      p.Async(m)
      time.Sleep(5 * time.Second)
   }
}
