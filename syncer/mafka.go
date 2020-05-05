package main

//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	"errors"

	"github.com/pingcap/tidb-binlog/drainer/loopbacksync"
	"github.com/pingcap/tidb-binlog/drainer/relay"
	dsync "github.com/pingcap/tidb-binlog/drainer/sync"
	"github.com/tsthght/syncer/orderlist"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/pingcap/tidb-binlog/drainer/translator"
)

type MafkaSyncer struct {
	toBeAckCommitTS orderlist.MapList
	configFile      string
	*BaseSyncer
}

func NewMafkaSyncer(
	cfg *dsync.DBConfig,
	cfgFile string,
	tableInfoGetter translator.TableInfoGetter,
	worker int,
	batchSize int,
	queryHistogramVec *prometheus.HistogramVec,
	sqlMode *string,
	destDBType string,
	relayer relay.Relayer,
	info *loopbacksync.LoopBackSync,
	enableDispatch bool,
	enableCausility bool) (dsyncer dsync.Syncer, err error) {
	if cfgFile == "" {
		return nil, errors.New("config file name is empty")
	}

	executor := &MafkaSyncer{}

	return executor, nil
}

func (ms *MafkaSyncer) Sync(item *dsync.Item) error {
	return nil
}

func (ms *MafkaSyncer) Close() error {
	return nil
}

func (ms *MafkaSyncer) SetSafeMode(mode bool) bool {
	return false
}

// go mod edit -replace=github.com/tsthght/tidb-binlog=github.com/tsthght/tidb-binlog@release3.0-plugin