package mafka

import (
	"s3common/s3castleclient"
)

func newCastleClient(cfg *ProducerConfig) (castleManager *s3castleclient.CastleClientManager, err error) {
	var ip, host string
	err, ip = GetLocalIP()
	if err != nil {
		return nil, err
	}
	err, host = GetLocalHost()
	if err != nil {
		return nil, err
	}
	if castleManager, err = s3castleclient.NewCastleClientManager(host, ip,
		cfg.CastleAppkey, cfg.CastleName, cfg.BusinessAppkey); err != nil {
		return nil, err
	}
	return
}
