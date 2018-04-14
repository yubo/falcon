/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/lib/tsdb"
	"golang.org/x/net/context"
)

// shard -> bucket -> item
type cacheModule struct {
	//sync.RWMutex
	//service *Service
	buckets []*cacheBucket

	idxQueue queue //index update queue

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *cacheModule) prestart(s *Service) error {
	//p.service = s
	p.idxQueue.init()

	s.cache = p
	return nil
}

// TODO
func (p *cacheModule) start(s *Service) (err error) {
	p.buckets = make([]*cacheBucket, SHARD_NUM)
	for i := 0; i < SHARD_NUM; i++ {
		p.buckets[i] = &cacheBucket{entries: make(map[string]*cacheEntry)}
	}

	p.reload(s)

	db, _, err := core.NewOrm("service_index", s.Conf.IdxDsn,
		s.Conf.DbMaxIdle, s.Conf.DbMaxConn)
	if err != nil {
		return err
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go p.indexWorker(db)
	go p.cleanWorker()

	return nil
}

func (p *cacheModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *cacheModule) reload(s *Service) error {
	ids := s.Conf.ShardIds
	idmap := make(map[int]bool, len(ids))

	for i := 0; i < len(ids); i++ {
		id := ids[i]
		if id < 0 || id >= SHARD_NUM {
			return core.EINVAL
		}
		idmap[id] = true
	}

	for id, bucket := range p.buckets {
		if idmap[id] {
			bucket.setState(CACHE_BUCKET_ENABLE)
		} else {
			bucket.setState(CACHE_BUCKET_DISABLE)
		}
	}

	return nil
}

func (p *cacheModule) put(dp *tsdb.DataPoint) (*cacheEntry, error) {
	bucket, err := p.getBucket(dp.Key.ShardId)
	if err != nil {
		return nil, err
	}

	e, err := bucket.getCacheEntry(string(dp.Key.Key))
	if err != nil {
		e, err = bucket.createCacheEntry(dp)
		if err != nil {
			return nil, err
		}
		p.idxQueue.enqueue(&e.list)
	}
	return e, e.put(dp)
}

func (p *cacheModule) get(key *tsdb.Key, start, end int64) (*tsdb.DataPoints, error) {
	bucket, err := p.getBucket(key.ShardId)
	if err != nil {
		return nil, err
	}

	e, err := bucket.getCacheEntry(string(key.Key))
	if err != nil {
		return nil, err
	}

	return &tsdb.DataPoints{
		Key:    key,
		Values: e.getValues(start, end),
	}, nil
}

func (p *cacheModule) getBucket(shardId int32) (*cacheBucket, error) {
	if shardId < 0 || shardId >= SHARD_NUM {
		return nil, core.ErrNoExits
	}
	bucket := p.buckets[shardId]

	if bucket.getState() != CACHE_BUCKET_ENABLE {
		return nil, core.ErrNoExits
	}
	return bucket, nil
}

func (p *cacheModule) cleanWorker() {
	ticker := time.NewTicker(time.Second * CACHE_CLEAN_INTERVAL).C

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			for _, bucket := range p.buckets {
				bucket.clean(CACHE_EXPIRE_TIME)
			}
		}
	}
}

func (p *cacheModule) indexWorker(db orm.Ormer) {
	//ticker := falconTicker(time.Second/INDEX_QPS, b.Conf.Debug)
	ticker := time.NewTicker(time.Second / INDEX_QPS).C

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			statsInc(ST_INDEX_TICK, 1)
			l := p.idxQueue.dequeue()
			if l == nil {
				time.Sleep(time.Second)
				continue
			}

			e := list_entry(l)
			indexUpdate(e, db)
		}
	}
}
