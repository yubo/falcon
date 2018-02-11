/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/tsdb"
	"golang.org/x/net/context"
)

var (
	// 13*2h
	bucketNum uint8 = 13
	// 2h
	bucketSize         uint64 = 2 * 3600
	keyWriterNum              = 2
	keyWriterQueueSize uint32 = 100
	logWiterNum               = 2
	logWriterQueueSize uint32 = 100

	// 15min
	allowTimestampBehind uint32 = 15 * 60

	// 1 min
	allowTimestampAhead uint32 = 60

	// 1 hour
	cleanInterval = time.Hour

	// 1o min
	finalizeInterval = 10 * time.Minute

	adjustTimestamp = false
)

// shard -> bucket -> item
type TsdbModule struct {
	sync.RWMutex
	buckets map[int]*tsdb.BucketMap

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *PutRequest) PrintForDebug() string {
	var res string
	res += fmt.Sprintf("--PutRequest PrintForDebug--\n")
	for _, dp := range p.Data {
		res += fmt.Sprintf("%4s%20s%10s%3d%8s%5d%8s%15f\n", "key:", string(dp.Key.Key),
			"shardId:", dp.Key.ShardId, "time:", dp.Value.Timestamp, "value:", dp.Value.Value)
	}
	return res
}

func (g *GetResponse) PrintForDebug() string {
	var res string
	res += fmt.Sprintf("--GetResponse PrintForDebug--\n")
	for _, dp := range g.Data {
		res += fmt.Sprintf("key: %s\n", string(dp.Key.Key))
		for _, v := range dp.Values {
			res += fmt.Sprintf("%10s%5d%8s%10f\n", "timestamp:", v.Timestamp, "value:", v.Value)
		}
	}
	return res
}

func (t *TsdbModule) prestart(s *Service) error {
	s.tsdb = t
	return nil
}

func (t *TsdbModule) start(s *Service) (err error) {
	dataDirectory := s.Conf.Configer.Str(C_TSDB_DIR)
	t.ctx, t.cancel = context.WithCancel(context.Background())

	keyWriters := make([]*tsdb.KeyListWriter, keyWriterNum)
	for i := 0; i < len(keyWriters); i++ {
		keyWriters[i] = tsdb.NewKeyListWriter(dataDirectory, keyWriterQueueSize)
	}

	bucketLogWriters := make([]*tsdb.BucketLogWriter, logWiterNum)
	for i := 0; i < len(bucketLogWriters); i++ {
		bucketLogWriters[i] = tsdb.NewBucketLogWriter(bucketSize, dataDirectory,
			logWriterQueueSize, allowTimestampBehind)
	}

	t.buckets = make(map[int]*tsdb.BucketMap)
	for i := 0; i < falcon.SHARD_NUM; i++ {
		k := keyWriters[rand.Intn(len(keyWriters))]
		b := bucketLogWriters[rand.Intn(len(bucketLogWriters))]
		if err := createShardPath(i, dataDirectory); err != nil {
			return err
		}
		newMap := tsdb.NewBucketMap(bucketNum, bucketSize, int32(i),
			dataDirectory, k, b, tsdb.UNOWNED)

		t.buckets[i] = newMap
	}

	t.reload(s)

	go t.cleanWorker()
	go t.finalizeBucketWorker()

	return nil
}

func (t *TsdbModule) stop(s *Service) error {
	t.cancel()

	var shardToBeDropped []int
	for k, v := range t.buckets {
		if v.GetState() != tsdb.UNOWNED {
			shardToBeDropped = append(shardToBeDropped, k)
		}
	}

	t.dropShard(shardToBeDropped)

	return nil
}

func (t *TsdbModule) reload(s *Service) error {
	t.RLock()
	defer t.RUnlock()

	ids := strings.Split(s.Conf.Configer.Str(C_SHARD_IDS), ",")
	newMap := make(map[int]bool, len(ids))
	for i := 0; i < len(ids); i++ {
		id, err := strconv.Atoi(ids[i])
		if err != nil {
			return err
		}
		newMap[id] = true
	}

	var shardToBeAdded []int
	var shardToBeDropped []int

	for k, v := range t.buckets {
		state := v.GetState()
		if _, ok := newMap[k]; ok {
			if state == tsdb.UNOWNED {
				shardToBeAdded = append(shardToBeAdded, k)
			} else if state == tsdb.PRE_UNOWNED {
				v.CancelUnowning()
			}
		} else if state == tsdb.OWNED {
			shardToBeDropped = append(shardToBeDropped, k)
		}
	}

	sort.Ints(shardToBeAdded)
	sort.Ints(shardToBeDropped)

	go t.addShard(shardToBeAdded)
	go t.dropShard(shardToBeDropped)

	return nil
}

func (t *TsdbModule) put(req *PutRequest) (*PutResponse, error) {
	t.RLock()
	defer t.RUnlock()

	res := &PutResponse{}
	for _, dp := range req.Data {
		m := t.buckets[int(dp.Key.ShardId)]
		if m == nil {
			return res, falcon.ErrNoExits
		}

		// Adjust 0, late, or early timestamps to now. Disable only for testing.
		now := timer.now()
		if adjustTimestamp {
			if dp.Value.Timestamp == 0 || dp.Value.Timestamp < now-int64(allowTimestampBehind) ||
				dp.Value.Timestamp > now+int64(allowTimestampAhead) {
				dp.Value.Timestamp = now
			}
		}

		newRows, dataPoints, err := m.Put(string(dp.Key.Key), &tsdb.TimeValuePair{Value: dp.Value.Value,
			Timestamp: dp.Value.Timestamp}, 0, false)
		if err != nil {
			return res, err
		}

		if newRows == tsdb.NOT_OWNED && dataPoints == tsdb.NOT_OWNED {
			return res, fmt.Errorf("key not own!")
		}

		res.N++
	}

	return res, nil
}

func (t *TsdbModule) _get(key *tsdb.Key, start, end int64) (res []*tsdb.TimeValuePair, err error) {
	t.RLock()
	defer t.RUnlock()

	m := t.buckets[int(key.ShardId)]
	if m == nil {
		return nil, falcon.EEXIST
	}

	state := m.GetState()
	switch state {
	case tsdb.UNOWNED:
		return nil, fmt.Errorf("Don't own shard %d", key.ShardId)

	case tsdb.PRE_OWNED, tsdb.READING_KEYS, tsdb.READING_KEYS_DONE,
		tsdb.READING_LOGS, tsdb.PROCESSING_QUEUED_DATA_POINTS:
		return nil, fmt.Errorf("Shard %d in progress", key.ShardId)
	default:
		datas, err := m.Get(string(key.Key), start, end)
		if err != nil {
			return res, err
		}

		for _, dp := range datas {
			res = append(res, &tsdb.TimeValuePair{Value: dp.Value, Timestamp: dp.Timestamp})
		}

		if state == tsdb.READING_BLOCK_DATA {
			glog.V(1).Infof("Shard %d in process", key.ShardId)
		}

		if start < m.GetReliableDataStartTime() {
			glog.V(1).Infof("missing too much data")
		}

		return res, nil

	}
}

func (t *TsdbModule) get(req *GetRequest) (res *GetResponse, err error) {

	res = &GetResponse{Data: make([]*tsdb.DataPoints, len(req.Keys))}
	for i, key := range req.Keys {
		dps := &tsdb.DataPoints{Key: key}
		if dps.Values, err = t._get(key, req.Start, req.End); err != nil {
			return
		}
		res.Data[i] = dps
	}
	return
}

func createShardPath(shardId int, dataDirectory string) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", dataDirectory, shardId), 0755)
}

func (t *TsdbModule) cleanWorker() {
	ticker := time.NewTicker(cleanInterval).C

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker:
			t.RLock()
			defer t.RUnlock()

			for _, v := range t.buckets {
				state := v.GetState()
				if state == tsdb.OWNED {
					v.CompactKeyList()
					err := v.DeleteOldBlockFiles()
					if err != nil {
						glog.Infof("%s %v", MODULE_NAME, err)
						continue
					}
				} else if state == tsdb.PRE_UNOWNED {
					err := v.SetState(tsdb.UNOWNED)
					if err != nil {
						glog.Infof("%s %v", MODULE_NAME, err)
						continue
					}
				}
			}
		}
	}
}

func (t *TsdbModule) finalizeBucketWorker() {
	ticker := time.NewTicker(finalizeInterval).C

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker:
			timestamp := timer.now() - int64(allowTimestampAhead+allowTimestampBehind) -
				int64(tsdb.Duration(1, bucketSize))
			t.finalizeBucket(timestamp)
		}
	}
}

func (t *TsdbModule) addShard(shards []int) {
	t.RLock()
	t.RUnlock()

	for _, shardId := range shards {
		m := t.buckets[shardId]
		if m == nil {
			glog.Infof("%s Invalid shardId :%d", MODULE_NAME, shardId)
			continue
		}

		if m.GetState() >= tsdb.PRE_OWNED {
			continue
		}

		if err := m.SetState(tsdb.PRE_OWNED); err != nil {
			glog.Infof("%s %v", MODULE_NAME, err)
			continue
		}

		if err := m.ReadKeyList(); err != nil {
			glog.Infof("%s %v", MODULE_NAME, err)
			continue
		}

		if err := m.ReadData(); err != nil {
			glog.Infof("%s %v", MODULE_NAME, err)
			continue
		}
	}

	go func() {
		for _, shardId := range shards {
			m := t.buckets[shardId]
			if m == nil {
				glog.Infof("%s Invalid shardId :%d", MODULE_NAME, shardId)
				continue
			}

			if m.GetState() != tsdb.READING_BLOCK_DATA {
				continue
			}

			more, err := m.ReadBlockFiles()
			if err != nil {
				glog.Infof("%s %v", MODULE_NAME, err)
				continue
			}

			for more {
				more, err = m.ReadBlockFiles()
				if err != nil {
					glog.Infof("%s %v", MODULE_NAME, err)
					break
				}
			}

		}
	}()
}

func (t *TsdbModule) dropShard(shards []int) error {
	t.RLock()
	defer t.RUnlock()

	for _, shardId := range shards {
		m := t.buckets[shardId]
		if m == nil {
			glog.Infof("%s Invalid shardId :%d", MODULE_NAME, shardId)
			continue
		}

		if m.GetState() != tsdb.OWNED {
			continue
		}

		if err := m.SetState(tsdb.PRE_UNOWNED); err != nil {
			glog.Infof("%s %v", MODULE_NAME, err)
			continue
		}
	}

	return nil
}

func (t *TsdbModule) finalizeBucket(timestamp int64) {
	go func() {
		t.RLock()
		defer t.RUnlock()

		for _, bucket := range t.buckets {
			bucketToFinalize := bucket.Bucket(timestamp)
			_, err := bucket.FinalizeBuckets(bucketToFinalize)
			if err != nil {
				glog.Infof("%s %v", MODULE_NAME, err)
				continue
			}
		}
	}()
}
