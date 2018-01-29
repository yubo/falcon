/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/huangaz/tsdb/lib/bucketLogWriter"
	"github.com/huangaz/tsdb/lib/bucketMap"
	"github.com/huangaz/tsdb/lib/dataTypes"
	"github.com/huangaz/tsdb/lib/keyListWriter"
	"github.com/yubo/falcon"

	"golang.org/x/net/context"
)

const (
	dataDirectory = "/tmp/tsdb"
	shardNum      = dataTypes.GORILLA_SHARDS
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
	allowedTimestampBehind uint32 = 15 * 60
	cleanInterval                 = time.Hour
	finalizeInterval              = 10 * time.Minute
)

// shard -> bucket -> item
type TsdbModule struct {
	sync.RWMutex
	service *Service
	buckets map[int]*bucketMap.BucketMap

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
	t.service = s
	return nil
}

func (t *TsdbModule) start(s *Service) (err error) {
	t.ctx, t.cancel = context.WithCancel(context.Background())

	keyWriters := make([]*keyListWriter.KeyListWriter, keyWriterNum)
	for i := 0; i < len(keyWriters); i++ {
		keyWriters[i] = keyListWriter.NewKeyListWriter(dataDirectory, keyWriterQueueSize)
	}

	bucketLogWriters := make([]*bucketLogWriter.BucketLogWriter, logWiterNum)
	for i := 0; i < len(bucketLogWriters); i++ {
		bucketLogWriters[i] = bucketLogWriter.NewBucketLogWriter(bucketSize, dataDirectory,
			logWriterQueueSize, allowedTimestampBehind)
	}

	t.buckets = make(map[int]*bucketMap.BucketMap)
	for shardId := 0; shardId < shardNum; shardId++ {
		k := keyWriters[rand.Intn(len(keyWriters))]
		b := bucketLogWriters[rand.Intn(len(bucketLogWriters))]
		if err := createShardPath(shardId); err != nil {
			return err
		}
		t.buckets[shardId] = bucketMap.NewBucketMap(bucketNum, bucketSize, int64(shardId),
			dataDirectory, k, b, bucketMap.UNOWNED)
	}

	t.reload(t.service)

	go t.cleanWorker()
	go t.finalizeBucketWorker()

	return nil
}

func (t *TsdbModule) stop(s *Service) error {
	t.cancel()
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
			if state == bucketMap.UNOWNED {
				shardToBeAdded = append(shardToBeAdded, k)
			} else if state == bucketMap.PRE_UNOWNED {
				v.CancelUnowning()
			}
		} else if state == bucketMap.OWNED {
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
			return res, falcon.EEXIST
		}

		newRows, dataPoints, err := m.Put(string(dp.Key.Key), dataTypes.DataPoint{Value: dp.Value.Value,
			Timestamp: dp.Value.Timestamp}, 0, false)
		if err != nil {
			return res, err
		}

		if newRows == bucketMap.NOT_OWNED && dataPoints == bucketMap.NOT_OWNED {
			return res, fmt.Errorf("key not own!")
		}

		res.N++
	}

	return res, nil
}

func (t *TsdbModule) _get(key *Key, start, end int64) (res []*TimeValuePair, err error) {
	t.RLock()
	defer t.RUnlock()

	m := t.buckets[int(key.ShardId)]
	if m == nil {
		return nil, falcon.EEXIST
	}

	state := m.GetState()
	switch state {
	case bucketMap.UNOWNED:
		return nil, fmt.Errorf("Don't own shard %d", key.ShardId)

	case bucketMap.PRE_OWNED, bucketMap.READING_KEYS, bucketMap.READING_KEYS_DONE,
		bucketMap.READING_LOGS, bucketMap.PROCESSING_QUEUED_DATA_POINTS:
		return nil, fmt.Errorf("Shard %d in progress", key.ShardId)
	default:
		datas, err := m.Get(string(key.Key), start, end)
		if err != nil {
			return res, err
		}

		for _, dp := range datas {
			res = append(res, &TimeValuePair{Value: dp.Value, Timestamp: dp.Timestamp})
		}

		if state == bucketMap.READING_BLOCK_DATA {
			log.Printf("Shard %d in process", key.ShardId)
		}

		if start < m.GetReliableDataStartTime() {
			log.Printf("missing too much data")
		}

		return res, nil

	}
}

func (t *TsdbModule) get(req *GetRequest) (res *GetResponse, err error) {

	res = &GetResponse{Data: make([]*DataPoints, len(req.Keys))}
	for i, key := range req.Keys {
		if res.Data[i].Values, err = t._get(key, req.Start, req.End); err != nil {
			return
		}
	}
	return
}

func createShardPath(shardId int) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", dataDirectory, shardId), 0755)
}

func (t *TsdbModule) cleanWorker() {
	ticker := time.NewTicker(cleanInterval).C

	for {
		select {
		case <-ticker:
			t.RLock()
			defer t.RUnlock()

			for _, v := range t.buckets {
				state := v.GetState()
				if state == bucketMap.OWNED {
					v.CompactKeyList()
					err := v.DeleteOldBlockFiles()
					if err != nil {
						log.Println(err)
						continue
					}
				} else if state == bucketMap.PRE_UNOWNED {
					err := v.SetState(bucketMap.UNOWNED)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *TsdbModule) finalizeBucketWorker() {
	ticker := time.NewTicker(finalizeInterval).C

	for {
		select {
		case <-ticker:
			finalizeTimeSeries()
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *TsdbModule) addShard(shards []int) {
	t.RLock()
	t.RUnlock()

	for _, shardId := range shards {
		m := t.buckets[shardId]
		if m == nil {
			log.Printf("Invalid shardId :%d", shardId)
			continue
		}

		if m.GetState() >= bucketMap.PRE_OWNED {
			continue
		}

		if err := m.SetState(bucketMap.PRE_OWNED); err != nil {
			log.Println(err)
			continue
		}

		if err := m.ReadKeyList(); err != nil {
			log.Println(err)
			continue
		}

		if err := m.ReadData(); err != nil {
			log.Println(err)
			continue
		}
	}

	go func() {
		for _, shardId := range shards {
			m := t.buckets[shardId]
			if m == nil {
				log.Printf("Invalid shardId :%d", shardId)
				continue
			}

			if m.GetState() != bucketMap.READING_BLOCK_DATA {
				continue
			}

			for more, _ := m.ReadBlockFiles(); more; more, _ = m.ReadBlockFiles() {
			}

			/*
				more, err := m.ReadBlockFiles()
				if err != nil {
					log.Println(err)
					continue
				}

				for more {
					more, err = m.ReadBlockFiles()
					if err != nil {
						log.Println(err)
						break
					}
				}
			*/

		}
	}()
}

func (t *TsdbModule) dropShard(shards []int) error {
	t.RLock()
	defer t.RUnlock()

	for _, shardId := range shards {
		m := t.buckets[shardId]
		if m == nil {
			log.Printf("Invalid shardId :%d", shardId)
			continue
		}

		if m.GetState() != bucketMap.OWNED {
			continue
		}

		if err := m.SetState(bucketMap.PRE_UNOWNED); err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}

func finalizeTimeSeries() {
}
