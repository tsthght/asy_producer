package syncer

import dsync "github.com/pingcap/tidb-binlog/drainer/sync"

type MetaInfo struct {
	size int
	item *dsync.Item
}
