package mafka

import (
	"log"
	"sync"
	"time"

	"github.com/tsthght/syncer/config"
	"github.com/tsthght/syncer/message"
	"s3common/s3mafkaclient"
)

type AsyProducer struct {
	asyProducer *s3mafkaclient.MafkaAsynProducer

	toBeAckTotalSize       int
	toBeAckCommitTSMu      sync.Mutex
	resumeProduce          chan struct{}
	resumeProduceCloseOnce sync.Once

	LastApplyTimestamp  int64
	LastSuccessTime int64

	CallBack MafkaCallBack

	cfg *config.ProducerConfig

	shutdown chan struct{}
}

func NewAsyProducer(cfg *config.ProducerConfig) (*AsyProducer, error) {
	producer := &AsyProducer{}
	producer.cfg = cfg
	castle, err := newCastleClient(producer.cfg)
	if err != nil {
		return nil, err
	}

	producer.CallBack = MafkaCallBack{make(chan []interface{}), make(chan interface{})}

	producer.asyProducer, err = newAsynProducer(castle, producer.cfg, producer.CallBack)
	if err != nil {
		return nil, err
	}

	producer.shutdown = make(chan struct{})
	producer.LastSuccessTime = time.Now().Unix()

	return producer, nil
}

func (p *AsyProducer) Close() {
	close(p.shutdown)
}

func (p *AsyProducer) GetWaitThreshold() int64 {
	return p.cfg.WaitThreshold
}

func (p *AsyProducer) Async(msg interface{}) string {
	m := msg.(message.Message)
	if err := p.asyProducer.SendMessageToChan(m); err != nil {
		return err.Error()
	}

	waitResume := false

	p.toBeAckCommitTSMu.Lock()
	p.toBeAckTotalSize += len(m.Msg)
	if p.toBeAckTotalSize > p.cfg.StallThreshold {
		p.resumeProduce = make(chan struct{})
		p.resumeProduceCloseOnce = sync.Once{}
		waitResume = true
	}
	p.toBeAckCommitTSMu.Unlock()

	if waitResume {
		select {
		case <-p.resumeProduce:
			return ""
		}
	}

	return ""
}

func (p *AsyProducer) Run () {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case msgs := <- p.CallBack.SuccessChan:
				for _, msg := range msgs {
					m := msg.(message.Message)
					p.LastApplyTimestamp = m.ApplyTime

					p.toBeAckCommitTSMu.Lock()
					p.LastSuccessTime = time.Now().Unix()
					p.toBeAckTotalSize -= len(m.Msg)
					if p.toBeAckTotalSize < p.cfg.StallThreshold && p.resumeProduce != nil {
						p.resumeProduceCloseOnce.Do(func() {
							close(p.resumeProduce)
						})
					}
					p.toBeAckCommitTSMu.Unlock()
				}
			case <-p.shutdown:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <- p.CallBack.FailureChan:
				e := err.(error)
				log.Fatal("fail to produce message to Mafka, please check the state of kafka server", e.Error())
			case <-p.shutdown:
				return
			}
		}
	}()

	for {
		select {
		case <-p.shutdown:
			wg.Wait()
			p.asyProducer.Close()
			return
		}
	}
}