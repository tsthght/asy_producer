package mafka

import (
	"errors"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tsthght/syncer/message"
	"go.uber.org/zap"
	"s3common/s3mafkaclient"
)

type AsyProducer struct {
	asyProducer *s3mafkaclient.MafkaAsynProducer

	lastSuccessTime time.Time
	toBeAckTotalSize       int
	toBeAckCommitTSMu      sync.Mutex
	resumeProduce          chan struct{}
	resumeProduceCloseOnce sync.Once

	LastApplyTimestamp  int64

	CallBack MafkaCallBack

	cfg *ProducerConfig

	Consumer *s3mafkaclient.MafkaConsumer//for test
}

func NewAsyProducer(fn string) (*AsyProducer, error) {
	producer := &AsyProducer{}
	producer.cfg = NewProducerConfig()
	if err := producer.cfg.Parse(fn); err != nil {
		return nil, err
	}
	castle, err := newCastleClient(producer.cfg)
	if err != nil {
		return nil, err
	}

	producer.CallBack = MafkaCallBack{make(chan []interface{}), make(chan interface{})}

	producer.asyProducer, err = newAsynProducer(castle, producer.cfg, producer.CallBack)
	if err != nil {
		return nil, err
	}

	producer.Consumer, err = newConsumer(castle, producer.cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
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
		}
	}

	return ""
}

func (p *AsyProducer) GetProducerConfig () *ProducerConfig {
	return p.cfg
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
					p.lastSuccessTime = time.Now()
					p.toBeAckTotalSize -= len(m.Msg)
					if p.toBeAckTotalSize < p.cfg.StallThreshold && p.resumeProduce != nil {
						p.resumeProduceCloseOnce.Do(func() {
							close(p.resumeProduce)
						})
					}
					p.toBeAckCommitTSMu.Unlock()
				}
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
				log.Fatal("fail to produce message to Mafka, please check the state of kafka server", zap.Error(e))
			}
		}
	}()
}

func (p *AsyProducer) Consuem() {
	go p.Consumer.ConsumeMessage(&BasicHandler{1})
}

type ProducerConfig struct {
	fs                         *flag.FlagSet
	Topic                      string      `toml:"topic" json:"topic"`
	Group                      string      `toml:"group" json:"group"`
	BusinessAppkey             string      `toml:"business-appkey" json:"business-appkey"`
	CastleAppkey               string      `toml:"castle-appkey" json:"castle-appkey"`
	CastleName                 string      `toml:"castle-name" json:"castle-name"`
	StallThreshold             int         `toml:"stall-threshold" json:"stall-threshold"`
	WaitThreshold              int64       `toml:"wait-threshold" json:"wait-threshold"`
	MaxRetryTimes              int         `toml:"max-retry-times" json:"max-retry-times"`
	MaxAsyncBufferChanSize     int64       `toml:"max-chan-size" json:"max-chan-size"`

	LocalHost                  string
}

func NewProducerConfig() *ProducerConfig {
	cfg := &ProducerConfig{}
	fs := flag.NewFlagSet("syncer", flag.ContinueOnError)
	fs.StringVar(&cfg.Topic, "topic", "", "mafka 's topic")
	fs.StringVar(&cfg.Group, "group", "", "mafka 's group")
	fs.StringVar(&cfg.BusinessAppkey, "business-appkey", "", "business appkey")
	fs.StringVar(&cfg.CastleAppkey, "castle-appkey", "", "castle appkey")
	fs.StringVar(&cfg.CastleName, "castle-name", "", "castle name")
	fs.IntVar(&cfg.StallThreshold, "stall-threshold", 90 * 1024 * 1024, "stall threshold")
	fs.Int64Var(&cfg.WaitThreshold, "wait-threshold", 30000, "wait threshold (ms)")
	fs.IntVar(&cfg.MaxRetryTimes, "max-retry-times", 10000, "max retry times")
	fs.Int64Var(&cfg.MaxAsyncBufferChanSize, "max-chan-size", 1 << 30, "max async buffer chan size")

	cfg.fs = fs
	return cfg
}

func (c *ProducerConfig) StrictDecodeFile(path string) error {
	metaData, err := toml.DecodeFile(path, c)
	if err != nil {
		return err
	}

	if undecoded := metaData.Undecoded(); len(undecoded) > 0 {
		var undecodedItems []string
		for _, item := range undecoded {
			undecodedItems = append(undecodedItems, item.String())
		}
		err = errors.New("contained unknown configuration options")
	}

	return err
}

func (c *ProducerConfig) Parse(filename string) error {
	// load config file if specified
	if filename == "" {
		return errors.New("config file name is nil")
	}
	if err := c.StrictDecodeFile(filename); err != nil {
		return err
	}
	return nil
}