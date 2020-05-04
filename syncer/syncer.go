package syncer

import (
	"sync"
	"time"
	"github.com/tsthght/asy_producer/mafka"
	"s3common/s3mafkaclient"

)

//DemoSyncer is a syncer demo
type MafkaSyncer struct {
	cfg *mafka.ProducerConfig
	asynProducer *s3mafkaclient.MafkaAsynProducer

	shutdown chan struct{}

	toBeAckCommitTSMu      sync.Mutex
	toBeAckCommitTS        map[int64]*MetaInfo
	toBeAckTotalSize       int
	resumeProduce          chan struct{}
	resumeProduceCloseOnce sync.Once
	lastSuccessTime time.Time
}
