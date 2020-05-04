package mafka

import (
	"s3common/s3castleclient"
	"s3common/s3mafkaclient"
)

func newAsynProducer(castleManager *s3castleclient.CastleClientManager, cfg *ProducerConfig,
	callback s3mafkaclient.SenderCallback) (asynProducer *s3mafkaclient.MafkaAsynProducer, err error) {
	if castleManager == nil {
		return nil, s3mafkaclient.ErrProducerInvalidConfig
	}

	if asynProducer, err = s3mafkaclient.NewMafkaAsynProducer(cfg.Topic, castleManager, nil, callback); err != nil {
		return nil, err
	} else {
		return asynProducer, nil
	}
}
