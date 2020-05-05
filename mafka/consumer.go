package mafka

import (
	"github.com/tsthght/syncer/message"
	"s3common/s3castleclient"
	"s3common/s3mafkaclient"
	"encoding/json"
	"fmt"
)

func newConsumer(castleManager *s3castleclient.CastleClientManager, cfg *ProducerConfig) (consumer *s3mafkaclient.MafkaConsumer, err error) {
	if castleManager == nil {
		return nil, s3mafkaclient.ErrProducerInvalidConfig
	}

	if consumer, err = s3mafkaclient.NewMafkaConsumerWithCastleManager(cfg.Topic, cfg.Group, castleManager); err != nil {
		return nil, err
	} else {
		return consumer, nil
	}
}

type BasicHandler struct {
	ID int
}

func (handler *BasicHandler) RecvMessage(msg *s3mafkaclient.ConsumerMessage) int {
	fmt.Printf("#### handler: cluster=[%s] id=[%d] recvMessage messages topic=[%s] partition=[%d] offset=[%d] key=[%v], value=[%s] \n", msg.Cluster, handler.ID, msg.Topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))

	var mm message.Message
	if tmpErr := json.Unmarshal(msg.Value, &mm); tmpErr != nil {
		fmt.Printf("recvMsgsByConusmer:  json unmarshal fail, err=%s\n", tmpErr)
	} else {
		fmt.Printf("===== %v\n", mm)
	}
	/*if flag {
		var ob objectInfo
		if tmpErr := json.Unmarshal(msg.Value, &ob); tmpErr != nil {
			fmt.Printf("recvMsgsByConusmer:  json unmarshal fail, err=%s\n", tmpErr)
		} else {
			testRecvMsgsLock.Lock()
			testRecvMsgs = append(testRecvMsgs, ob)
			testRecvMsgsLock.Unlock()
		}
	}*/
	return s3mafkaclient.ConsumeStatusFailure
}