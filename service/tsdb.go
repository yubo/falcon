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
	"github.com/huangaz/tsdb/lib/keyListWriter"

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

func (p *TsdbModule) prestart(s *Service) error {
	s.tsdb = p
	return nil
}

func (p *TsdbModule) start(s *Service) (err error) {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	ids_ := strings.Split(s.Conf.Configer.Str(C_SHARD_IDS), ",")
	p.ids = make([]int, len(ids_))
	for i := 0; i < len(p.ids); i++ {
		p.ids[i], err = strconv.Atoi(ids_[i])
		if err != nil {
			return err
		}
	}
	p.buckets = make(map[int]*bucketMap.BucketMap)

	k := keyListWriter.NewKeyListWriter(dataDirectory, 100)
	b := bucketLogWriter.NewBucketLogWriter(4*3600, dataDirectory, 100, 0)

	for _, shardId := range p.ids {
		// todo
		if err := createShardPath(shardId); err != nil {
			return err
		}
		p.buckets[shardId] = bucketMap.NewBucketMap(6, 4*3600, int64(shardId), dataDirectory, k, b, bucketMap.UNOWNED)
	}

	// check
	go func() {
		for _, m := range p.buckets {
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
		for _, m := range p.buckets {
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

func (p *TsdbModule) stop(s *Service) error {
	return nil
}

func (p *TsdbModule) reload(s *Service) error {
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

		// diff new_ids p.ids

	*/

	return nil
}

// TODO
func (p *TsdbModule) put(req *PutRequest) (res *PutResponse, err error) {
	/*
		m := p.buckets[item.ShardId]
		if m == nil {
			return falcon.EEXIST
		}

		_, _, err := m.Put(item.Key, dataTypes.DataPoint{Value: item.Value,
			Timestamp: item.Timestamp}, 0, false)
		return err

	*/
	return
}

// TODO
func (p *TsdbModule) get(req *GetRequest) (res *GetResponse, err error) {
	return
}

func createShardPath(shardId int) error {
	return os.MkdirAll(fmt.Sprintf("%s/%d", dataDirectory, shardId), 0755)
}
