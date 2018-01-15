/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yubo/falcon"
	"github.com/yubo/gotool/list"
	"golang.org/x/net/context"
)

type IndexModule struct {
	ctx    context.Context
	cancel context.CancelFunc

	indexUpdateCh chan *itemEntry
	db            orm.Ormer
}

func indexUpdateEndpoint(endpoint string, lastTs int64, db orm.Ormer) (id int64, err error) {
	_, err = db.Raw("INSERT INTO endpoint(endpoint, ts, t_create) VALUES (?, ?, now()) ON DUPLICATE KEY UPDATE ts=?, t_modify=now()", endpoint, lastTs, lastTs).Exec()
	if err != nil {
		statsInc(ST_INDEX_HOST_INSERT_ERR, 1)
		return
	}

	err = db.Raw("SELECT id FROM endpoint WHERE endpoint = ?", endpoint).QueryRow(&id)
	return
}

func indexUpdateTagEndpoint(tags_ string, lastTs int64, db orm.Ormer, hid int64) (err error) {
	tags := strings.Split(tags_, ",")
	for _, tag := range tags {
		_, err = db.Raw("INSERT INTO tag_endpoint(tag, endpoint_id, ts, t_create) VALUES (?, ?, ?, now()) ON DUPLICATE KEY UPDATE ts=?, t_modify=now()", tag, hid, lastTs, lastTs).Exec()
		if err != nil {
			statsInc(ST_INDEX_TAG_INSERT_ERR, 1)
			return
		}
	}
	return
}

func indexUpdateEndpointCounter(counter, tags, typ string, ts int64, db orm.Ormer, hid int64) (err error) {
	if len(tags) > 0 {
		counter += "/" + tags
	}

	_, err = db.Raw("INSERT INTO endpoint_counter(endpoint_id,counter,type,ts,t_create) VALUES (?,?,?,?,now()) ON DUPLICATE KEY UPDATE ts=?, type=?,t_modify=now()", hid, counter, typ, ts, ts, typ).Exec()
	if err != nil {
		statsInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
	}
	return
}

func indexUpdate(e *itemEntry, db orm.Ormer) {

	statsInc(ST_INDEX_UPDATE, 1)

	hid, err := indexUpdateEndpoint(e.endpoint, e.lastTs, db)
	if err != nil {
		return
	}

	if err = indexUpdateTagEndpoint(e.tags, e.lastTs, db, hid); err != nil {
		return
	}

	if err = indexUpdateEndpointCounter(e.metric, e.tags, e.typ, e.lastTs, db, hid); err != nil {
		return
	}

	return
}

/*
 * lru　队列是以索引更新时间先进先出
 */
func (p *IndexModule) indexWorker(b *Service) {
	var (
		e   *itemEntry
		l   *list.ListHead
		now int64
	)
	//ticker := falconTicker(time.Second/INDEX_QPS, b.Conf.Debug)
	ticker := time.NewTicker(time.Second / INDEX_QPS).C
	newQueue := &b.shard.newQueue
	lruQueue := &b.shard.lruQueue

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			statsInc(ST_INDEX_TICK, 1)
			now = timer.now()

			if l = newQueue.dequeue(); l != nil {
				lruQueue.enqueue(l)

				// immediate update
				e = list_entry(l)
				e.idxTs = now
				p.indexUpdateCh <- e

				// ADD_HOOK
				addEntryHandle(e)
				continue
			}

			if l = lruQueue.dequeue(); l == nil {
				time.Sleep(time.Second)
				continue
			}

			e = list_entry(l)
			if now-e.idxTs < INDEX_UPDATE_CYCLE_TIME {
				lruQueue.addHead(l)
				time.Sleep(time.Second)
				continue
			}

			// remove timeout entry
			if now-e.lastTs > INDEX_TIMEOUT {
				// DEL_HOOK
				delEntryHandle(e)
				bucket, _ := b.shard.getBucket(e.shardId)
				bucket.unlink(e.key)
				continue
			}

			e.idxTs = now
			lruQueue.enqueue(l)
			p.indexUpdateCh <- e
		}
	}
}

func (p *IndexModule) updateWorker() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case e := <-p.indexUpdateCh:
			go indexUpdate(e, p.db)
		}
	}
}

// indexModule depend on cacheModule
func (p *IndexModule) prestart(b *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.indexUpdateCh = make(chan *itemEntry, INDEX_MAX_OPEN_CONNS)
	return nil
}

func (p *IndexModule) start(b *Service) error {
	var err error

	dbmaxidle, _ := b.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := b.Conf.Configer.Int(C_DB_MAX_CONN)

	p.db, err = falcon.NewOrm("service_index",
		b.Conf.Configer.Str(C_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	go p.indexWorker(b)
	go p.updateWorker()

	return nil
}

func (p *IndexModule) stop(b *Service) error {
	p.cancel()
	return nil
}

func (p *IndexModule) reload(b *Service) error {
	return nil
}
