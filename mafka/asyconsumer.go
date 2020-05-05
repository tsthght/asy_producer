package mafka

import (
	"github.com/tsthght/syncer/config"
	"s3common/s3mafkaclient"
)

type AsyConsumer struct {
	consumer *s3mafkaclient.MafkaConsumer
	cfg *config.ProducerConfig
}

func NewAsyConsumer(cfg *config.ProducerConfig) (*AsyConsumer, error) {
	consumer := &AsyConsumer{}
	consumer.cfg = cfg
	castle, err := newCastleClient(consumer.cfg)
	if err != nil {
		return nil, err
	}
	consumer.consumer, err = newConsumer(castle, cfg)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (c *AsyConsumer) PrintConsuem() {
	c.consumer.ConsumeMessage(&BasicHandler{1})
}

func (c *AsyConsumer) Close() {
	c.consumer.Close()
}