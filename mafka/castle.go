package mafka

import (
	"github.com/tsthght/syncer/config"
	"s3common/s3castleclient"
)

func newCastleClient(cfg *config.ProducerConfig) (castleManager *s3castleclient.CastleClientManager, err error) {
	var ip, host string
	err, ip = GetLocalIP()
	if err != nil {
		return nil, err
	}
	err, host = GetLocalHost(ip)
	if err != nil {
		return nil, err
	}
	if castleManager, err = s3castleclient.NewCastleClientManager(host, ip,
		cfg.CastleAppkey, cfg.CastleName, cfg.BusinessAppkey); err != nil {
		return nil, err
	}
	return
}
