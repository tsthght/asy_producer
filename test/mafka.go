package main
//#cgo CFLAGS: -I ../common
//#cgo LDFLAGS: -L ../common  -Wl,-rpath=/usr/local/lib -lcommon
//
//#include "libcommon.h"
import "C"
import (
	dsync "github.com/pingcap/tidb-binlog/drainer/sync"
	"github.com/tsthght/syncer/orderlist"
)

type MafkaSyncer struct {
	toBeAckCommitTS orderlist.MapList
	configFile      string
}

func (ms *MafkaSyncer) Sync(item *dsync.Item) error {
	return nil
}
