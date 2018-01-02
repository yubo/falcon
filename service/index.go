/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"database/sql"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
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

func indexUpdateEndpoint(e *itemEntry, db orm.Ormer) (id int64, err error) {
	if err = db.Raw("SELECT id FROM endpoint WHERE endpoint = ?",
		e.endpoint).QueryRow(&id); err == nil || err != sql.ErrNoRows {
		return
	}

	statsInc(ST_INDEX_HOST_INSERT, 1)
	res, err := db.Raw("INSERT INTO endpoint(endpoint, ts, t_create) VALUES (?, ?, now()) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)", e.endpoint, e.lastTs).Exec()
	if err != nil {
		statsInc(ST_INDEX_HOST_INSERT_ERR, 1)
		return
	}

	return res.LastInsertId()
}

func indexUpdateTagEndpoint(e *itemEntry, db orm.Ormer, hid int64) (err error) {
	var tid int64

	tags := strings.Split(string(e.tags), ",")
	for _, tag := range tags {

		if err = db.Raw("SELECT id FROM tag_endpoint WHERE tag = ? and endpoint_id = ?",
			tag, hid).QueryRow(&tid); err == nil {
			continue
		}
		if err != sql.ErrNoRows {
			return
		}

		statsInc(ST_INDEX_TAG_INSERT, 1)
		_, err = db.Raw("INSERT INTO tag_endpoint(tag, endpoint_id, ts, t_create) VALUES (?, ?, ?, now()) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)", tag, hid, e.lastTs).Exec()
		if err != nil {
			statsInc(ST_INDEX_TAG_INSERT_ERR, 1)
			return
		}
	}
	return
}

func indexUpdateEndpointCounter(e *itemEntry, db orm.Ormer, hid int64) (err error) {
	var id int64

	counter := e.key()
	err = db.Raw("SELECT id FROM counter WHERE endpoint_id = ? and counter = ?",
		hid, counter).QueryRow(&id)
	if err == nil || err != sql.ErrNoRows {
		return
	}

	statsInc(ST_INDEX_COUNTER_INSERT, 1)
	_, err = db.Raw("INSERT INTO endpoint_counter(endpoint_id,counter,type,ts,t_create) VALUES (?,?,?,?,now()) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts),type=VALUES(type)", hid, counter, e.typ, e.lastTs).Exec()
	if err != nil {
		statsInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
		return
	}
	return
}

func indexUpdate(e *itemEntry, db orm.Ormer) {

	statsInc(ST_INDEX_UPDATE, 1)

	hid, err := indexUpdateEndpoint(e, db)
	if err != nil {
		return
	}

	if err = indexUpdateTagEndpoint(e, db, hid); err != nil {
		return
	}

	if err = indexUpdateEndpointCounter(e, db, hid); err != nil {
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
				continue
			}

			e = list_entry(l)
			if now-e.idxTs < INDEX_UPDATE_CYCLE_TIME {
				lruQueue.addHead(l)
				continue
			}

			// remove timeout entry
			if now-e.lastTs > INDEX_TIMEOUT {
				// DEL_HOOK
				delEntryHandle(e)
				bucket, _ := b.shard.getBucket(e.shardId)
				bucket.unlink(e.key())
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

	glog.Info(MODULE_NAME, "indexStart ok")
	return nil
}

func (p *IndexModule) stop(b *Service) error {
	p.cancel()
	return nil
}

func (p *IndexModule) reload(b *Service) error {
	return nil
}
