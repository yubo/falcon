/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/huangaz/tsdb/lib/bucketLogWriter"
	"github.com/huangaz/tsdb/lib/bucketMap"
	"github.com/huangaz/tsdb/lib/dataTypes"
	"github.com/huangaz/tsdb/lib/keyListWriter"
	"github.com/yubo/falcon"

	"golang.org/x/net/context"
)

const (
	dataDirectory = "/tmp/tsdb"
)

// shard -> bucket -> item
type TsdbModule struct {
	sync.RWMutex
	service *Service
	buckets map[int]*bucketMap.BucketMap
	ids     []int

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *PutRequest) PrintForDebug() string {
	var res string
	res += fmt.Sprintf("--PutRequest PrintForDebug--\n")
	for _, item := range p.Items {
		res += fmt.Sprintf("%4s%20s%10s%3d%8s%5d%8s%15f\n", "key:", string(item.Key),
			"shardId:", item.ShardId, "time:", item.Timestamp, "value:", item.Value)
	}
	return res
}

func (g *GetResponse) PrintForDebug() string {
	var res string
	res += fmt.Sprintf("--GetResponse PrintForDebug--\n")
	res += fmt.Sprintf("key: %s\n", string(g.Key))
	for _, dp := range g.Dps {
		res += fmt.Sprintf("%10s%5d%8s%10f\n", "timestamp:", dp.Timestamp, "value:", dp.Value)
	}
	return res
}

func (t *TsdbModule) prestart(s *Service) error {
	s.tsdb = t
	return nil
}

func (t *TsdbModule) start(s *Service) (err error) {
	t.ctx, t.cancel = context.WithCancel(context.Background())
	ids_ := strings.Split(s.Conf.Configer.Str(C_SHARD_IDS), ",")
	t.ids = make([]int, len(ids_))
	for i := 0; i < len(t.ids); i++ {
		t.ids[i], err = strconv.Atoi(ids_[i])
		if err != nil {
			return err
		}
	}
	t.buckets = make(map[int]*bucketMap.BucketMap)

	k := keyListWriter.NewKeyListWriter(dataDirectory, 100)
	b := bucketLogWriter.NewBucketLogWriter(4*3600, dataDirectory, 100, 0)

	for _, shardId := range t.ids {
		// todo
		if err := createShardPath(shardId); err != nil {
			return err
		}
		t.buckets[shardId] = bucketMap.NewBucketMap(6, 4*3600, int64(shardId), dataDirectory,
			k, b, bucketMap.UNOWNED)
	}

	// check
	go func() {
		for _, m := range t.buckets {
			if ok := m.SetState(bucketMap.PRE_OWNED); !ok {
				glog.Fatal("set state failed")
			}
			if err := m.ReadKeyList(); err != nil {
				glog.Fatal(err)
			}
			if err := m.ReadData(); err != nil {
				glog.Fatal(err)
			}
		}
		for _, m := range t.buckets {
			more, err := m.ReadBlockFiles()
			if err != nil {
				glog.Fatal(err)
			}

			for more {
				more, err = m.ReadBlockFiles()
				if err != nil {
					glog.Fatal(err)
				}
			}
		}

	}()

	return nil
}

func (t *TsdbModule) stop(s *Service) error {
	return nil
}

func (t *TsdbModule) reload(s *Service) error {
	// TODO
	/*
		ids_ := strings.Split(s.Conf.Configer.Str(C_SHARD_IDS), ",")
		new_ids := make([]int, len(ids_))
		for i := 0; i < len(ids); i++ {
			ids[i], err = strconv.Atoi(ids_[i])
			if err != nil {
				return err
			}
		}

		// diff new_ids t.ids

	*/

	return nil
}

// TODO
func (t *TsdbModule) put(req *PutRequest) (*PutResponse, error) {
	res := &PutResponse{}
	for _, item := range req.Items {
		m := t.buckets[int(item.ShardId)]
		if m == nil {
			return res, falcon.EEXIST
		}

		newRows, dataPoints, err := m.Put(string(item.Key), dataTypes.DataPoint{Value: item.Value,
			Timestamp: item.Timestamp}, 0, false)
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

// TODO
func (t *TsdbModule) get(req *GetRequest) (*GetResponse, error) {
	res := &GetResponse{
		Key: make([]byte, len(req.Key)),
	}
	if len(req.Key) == 0 {
		return nil, fmt.Errorf("null key!")
	}

	m := t.buckets[int(req.ShardId)]
	if m == nil {
		return nil, falcon.EEXIST
	}

	copy(res.Key, req.Key)
	state := m.GetState()
	switch state {
	case bucketMap.UNOWNED:
		return nil, fmt.Errorf("Don't own shard %d", req.ShardId)
	case bucketMap.PRE_OWNED, bucketMap.READING_KEYS, bucketMap.READING_KEYS_DONE, bucketMap.READING_LOGS, bucketMap.PROCESSING_QUEUED_DATA_POINTS:
		return nil, fmt.Errorf("Shard %d in progress", req.ShardId)
	default:
		datas, err := m.Get(string(req.Key), req.Start, req.End)
		if err != nil {
			return res, err
		}

		for _, dp := range datas {
			res.Dps = append(res.Dps, &DataPoint{Value: dp.Value, Timestamp: dp.Timestamp})
		}

		if state == bucketMap.READING_BLOCK_DATA {
			return res, fmt.Errorf("Shard %d in progress", req.ShardId)
		} else if req.Start < m.GetReliableDataStartTime() {
			return res, fmt.Errorf("missing too much data")
		}

		return res, nil

	}
}

func createShardPath(shardId int) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", dataDirectory, shardId), 0755)
}
