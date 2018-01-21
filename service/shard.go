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
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
)

// shard -> bucket -> item
type ShardModule struct {
	sync.RWMutex
	service    *Service
	bucketMap  map[int32]*bucketEntry
	shardTotal int

	trashQueue queue //trash ie queue
	idxQueue   queue //index update queue
	putQueue   queue //put api queue

	ctx    context.Context
	cancel context.CancelFunc

	db orm.Ormer
}

func (p *ShardModule) prestart(s *Service) error {
	p.bucketMap = make(map[int32]*bucketEntry)
	s.shard = p
	return nil
}

// TODO
func (p *ShardModule) start(s *Service) (err error) {
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
			dpEntryMap: make(map[string]*dpEntry),
		}
	}
	p.service = s
	p.idxQueue.init()
	p.putQueue.init()
	p.trashQueue.init()

	dbmaxidle, _ := s.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := s.Conf.Configer.Int(C_DB_MAX_CONN)
	db, _, err := falcon.NewOrm("service_index",
		s.Conf.Configer.Str(C_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go p.indexWorker(db)
	go p.cleanWorker()

	return nil
}

func (p *ShardModule) stop(s *Service) error {
	return nil
}

func (p *ShardModule) reload(s *Service) error {
	return nil
}

func (p *ShardModule) getTrashItem() *dpEntry {
	l := p.trashQueue.dequeue()
	if l == nil {
		return nil
	}
	return list_entry(l)
}

func (p *ShardModule) putTrashItem(ie *dpEntry) {
	p.trashQueue.enqueue(&ie.list)
}

func (p *ShardModule) put(dp *DataPoint) (*dpEntry, error) {
	bucket, err := p.getBucket(dp.Key.ShardId)
	if err != nil {
		return nil, err
	}

	ie, err := bucket.getDpEntry(string(dp.Key.Key))
	if err != nil {
		ie, err = bucket.createDpEntry(dp)
		if err != nil {
			return nil, err
		}
		p.idxQueue.addHead(&ie.list)
		p.putQueue.enqueue(&ie.list_p)
		return ie, ie.put(dp)
	}

	if err := ie.put(dp); err != nil {
		return ie, err
	}

	p.putQueue.moveTail(&ie.list_p)
	return ie, nil
}

func (p *ShardModule) get(key *Key, start, end int64) (*DataPoints, error) {
	bucket, err := p.getBucket(key.ShardId)
	if err != nil {
		return nil, err
	}

	ie, err := bucket.getDpEntry(string(key.Key))
	if err != nil {
		return nil, err
	}

	return &DataPoints{
		Key:    key,
		Values: ie.getValues(start, end),
	}, nil
}

func (p *ShardModule) getBucket(shardId int32) (*bucketEntry, error) {
	p.RLock()
	defer p.RUnlock()

	if bucket, ok := p.bucketMap[shardId]; ok {
		return bucket, nil
	}
	return nil, falcon.ErrNoExits
}

func (p *ShardModule) cleanWorker() {
	var now int64

	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			now = timer.now()
			for l := p.putQueue.dequeue(); l != nil; l = p.putQueue.dequeue() {
				e := list_p_entry(l)

				if now-e.lastTs > INDEX_EXPIRE_TIME {
					// DEL_HOOK
					p.delEntryHandle(e)

				} else {
					p.putQueue.addHead(l)
					break
				}
			}
		}
	}
}

func (p *ShardModule) indexWorker(db orm.Ormer) {
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

// TODO fix me
func (p *ShardModule) addEntryHandle(e *dpEntry) {
}

// already del from p.putQueue
func (p *ShardModule) delEntryHandle(e *dpEntry) {
	bucket, _ := p.getBucket(e.key.ShardId)
	bucket.unlink(string(e.key.Key))

	p.idxQueue.Lock()
	defer p.idxQueue.Unlock()
	e.list.Del()
}
