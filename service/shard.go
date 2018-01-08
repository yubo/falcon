/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

// shard -> bucket -> item

type ShardModule struct {
	sync.RWMutex
	service    *Service
	newQueue   cacheq //immediate queue
	lruQueue   cacheq //lru queue
	bucketMap  map[int32]*bucketEntry
	shardTotal int
}

func (p *ShardModule) put(item *Item) (*itemEntry, error) {
	bucket, err := p.getBucket(item.ShardId)
	if err != nil {
		return nil, err
	}

	ie, err := bucket.getItem(string(item.Key))
	if err != nil {
		ie, err = bucket.addItem(item)
		if err != nil {
			return nil, err
		}
		p.newQueue.enqueue(&ie.list)
	}

	return ie, ie.put(item)
}

func (p *ShardModule) get(req *GetRequest) ([]*DataPoint, error) {
	bucket, err := p.getBucket(req.ShardId)
	if err != nil {
		return nil, err
	}

	ie, err := bucket.getItem(string(req.Key))
	if err != nil {
		return nil, err
	}

	return ie.getDps(req.Start, req.End)
}

func (p *ShardModule) getBucket(shardId int32) (*bucketEntry, error) {
	p.RLock()
	defer p.RUnlock()

	if bucket, ok := p.bucketMap[shardId]; ok {
		return bucket, nil
	}
	return nil, falcon.ErrNoExits
}

func (p *ShardModule) prestart(s *Service) (err error) {
	glog.V(3).Infof(MODULE_NAME + " cache prestart \n")
	p.bucketMap = make(map[int32]*bucketEntry)

	ids_ := strings.Split(s.Conf.Configer.Str(C_SHARD_IDS), ",")
	ids := make([]int, len(ids_))
	for i := 0; i < len(ids); i++ {
		ids[i], err = strconv.Atoi(ids_[i])
		if err != nil {
			return err
		}
	}

	for _, shardId := range ids {
		p.bucketMap[int32(shardId)] = &bucketEntry{
			itemMap: make(map[string]*itemEntry),
		}
	}
	p.newQueue.init()
	p.lruQueue.init()
	s.shard = p
	p.service = s
	return nil
}

func (p *ShardModule) start(s *Service) error {
	glog.V(3).Infof(MODULE_NAME + " cache start \n")
	return nil
}

func (p *ShardModule) stop(s *Service) error {
	//p.cache.close()
	return nil
}

func (p *ShardModule) reload(s *Service) error {
	return nil
}
