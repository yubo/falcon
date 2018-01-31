/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/tsdb"
	"golang.org/x/net/context"
)

// shard -> bucket -> item
type CacheModule struct {
	//sync.RWMutex
	service *Service
	buckets []*cacheBucket

	idxQueue queue //index update queue

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *CacheModule) prestart(s *Service) error {
	p.service = s
	s.cache = p
	return nil
}

// TODO
func (p *CacheModule) start(s *Service) (err error) {
	conf := &s.Conf.Configer

	p.buckets = make([]*cacheBucket, falcon.SHARD_NUM)
	for i := 0; i < falcon.SHARD_NUM; i++ {
		p.buckets[i] = &cacheBucket{entries: make(map[string]*cacheEntry)}
	}

	p.reload(s)

	p.idxQueue.init()

	dbmaxidle, _ := conf.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := conf.Int(C_DB_MAX_CONN)
	db, _, err := falcon.NewOrm("service_index", conf.Str(C_DSN),
		dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go p.indexWorker(db)
	go p.cleanWorker()

	return nil
}

func (p *CacheModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *CacheModule) reload(s *Service) error {
	conf := &s.Conf.Configer

	ids := strings.Split(conf.Str(C_SHARD_IDS), ",")
	idmap := make(map[int]bool, len(ids))

	for i := 0; i < len(ids); i++ {
		id, err := strconv.Atoi(ids[i])
		if err != nil {
			return err
		}
		if id < 0 || id >= falcon.SHARD_NUM {
			return falcon.EINVAL
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

func (p *CacheModule) put(dp *tsdb.DataPoint) (*cacheEntry, error) {
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
		p.idxQueue.addHead(&e.list)
	}
	return e, e.put(dp)
}

func (p *CacheModule) get(key *tsdb.Key, start, end int64) (*tsdb.DataPoints, error) {
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

func (p *CacheModule) getBucket(shardId int32) (*cacheBucket, error) {
	if shardId < 0 || shardId >= falcon.SHARD_NUM {
		return nil, falcon.ErrNoExits
	}
	bucket := p.buckets[shardId]

	if bucket.getState() != CACHE_BUCKET_ENABLE {
		return nil, falcon.ErrNoExits
	}
	return bucket, nil
}

func (p *CacheModule) cleanWorker() {
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

func (p *CacheModule) indexWorker(db orm.Ormer) {
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
			now := timer.now()

			if now-e.idxTs > INDEX_UPDATE_INTERVAL {
				// update
				e.idxTs = now
				p.idxQueue.enqueue(l)
				indexUpdate(e, db)
			} else {
				p.idxQueue.addHead(l)
				time.Sleep(time.Second)
			}
		}
	}
}
