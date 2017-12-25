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

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/gotool/list"
	"golang.org/x/net/context"
)

type IndexModule struct {
	indexUpdateCh chan *cacheEntry
	indexDb       *sql.DB
	ctx           context.Context
	cancel        context.CancelFunc
}

func indexUpdate(e *cacheEntry, db *sql.DB) {
	var (
		err           error
		hid, cid, tid int64
		dstype        string
	)

	statsInc(ST_INDEX_UPDATE, 1)
	hid = -1
	err = db.QueryRow("SELECT id FROM endpoint WHERE endpoint = ?", e.endpoint).Scan(&hid)
	if err != nil {
		statsInc(ST_INDEX_HOST_MISS, 1)
		if err == sql.ErrNoRows || hid < 0 {
			statsInc(ST_INDEX_HOST_INSERT, 1)
			ret, err := db.Exec("INSERT INTO endpoint(endpoint, ts, t_create) "+
				"VALUES (?, ?, now()) ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
				e.endpoint, e.lastTs)
			if err != nil {
				statsInc(ST_INDEX_HOST_INSERT_ERR, 1)
				glog.Warning(MODULE_NAME, err)
				return
			}

			hid, err = ret.LastInsertId()
			if err != nil {
				glog.Warning(MODULE_NAME, err)
				return
			}
		} else {
			glog.Warning(MODULE_NAME+string(e.endpoint), err)
			return
		}
	}

	// tag
	tags := strings.Split(string(e.tags), ",")
	for _, tag := range tags {

		tid = -1
		err := db.QueryRow("SELECT id FROM tag WHERE tag = ? and "+
			"endpoint_id = ?", tag, hid).Scan(&tid)
		if err != nil {
			statsInc(ST_INDEX_TAG_MISS, 1)
			if err == sql.ErrNoRows || tid < 0 {
				statsInc(ST_INDEX_TAG_INSERT, 1)
				ret, err := db.Exec("INSERT INTO tag_endpoint(tag, endpoint_id, "+
					"ts, t_create) "+
					"VALUES (?, ?, ?, now()) "+
					"ON DUPLICATE KEY "+
					"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
					tag, hid, e.lastTs)
				if err != nil {
					statsInc(ST_INDEX_TAG_INSERT_ERR, 1)
					glog.Warning(MODULE_NAME, err)
					return
				}

				tid, err = ret.LastInsertId()
				if err != nil {
					glog.Warning(MODULE_NAME, err)
					return
				}
			} else {
				glog.Warning(MODULE_NAME+tag, hid, err)
				return
			}
		}

	}

	// endpoint_id
	counter := e.key()

	cid = -1
	dstype = "nil"

	err = db.QueryRow("SELECT id, type FROM counter WHERE "+
		"endpoint_id = ? and counter = ?",
		hid, counter).Scan(&cid, &dstype)
	if err != nil {
		statsInc(ST_INDEX_COUNTER_MISS, 1)
		if err == sql.ErrNoRows || cid < 0 {
			statsInc(ST_INDEX_COUNTER_INSERT, 1)
			ret, err := db.Exec("INSERT INTO endpoint_counter(endpoint_id,counter,"+
				"type,ts,t_create) "+
				"VALUES (?,?,?,?,now()) "+
				"ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts),"+
				"type=VALUES(type)",
				hid, counter, e.typ,
				e.lastTs)
			if err != nil {
				statsInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
				glog.Warning(MODULE_NAME, err)
				return
			}

			cid, err = ret.LastInsertId()
			if err != nil {
				glog.Warning(MODULE_NAME, err)
				return
			}
		} else {
			glog.Warning(MODULE_NAME+counter, hid, err)
			return
		}
	} /* else {
		if !(e.step == step && e.typ == dstype) {
			_, err := db.Exec("UPDATE counter SET step = ?, "+
				"type = ? where id = ?",
				e.step, e.typ, cid)
			if err != nil {
				glog.Warning(err)
				return
			}
		}
	}*/
	return
}

func (this *IndexModule) indexTrashWorker(b *Service) {
	var (
		e     *cacheEntry
		p, _p *list.ListHead
	)
	ticker := falconTicker(time.Second*INDEX_TRASH_LOOPTIME,
		b.Conf.Debug)
	q0 := &b.cache.idx0q
	q2 := &b.cache.idx2q
	for {
		select {
		case <-this.ctx.Done():
			return
		case <-ticker:
			for p = q2.head.Next; p != &q2.head; p = p.Next {
				_p = p.Next
				e = list_idx_entry(p)
				if b.timeNow()-e.lastTs < INDEX_TIMEOUT {
					q2.Lock()
					p.Del()
					//q2.size--
					q2.Unlock()

					e.idxTs = 0
					statsInc(ST_INDEX_TRASH_PICKUP, 1)
					q0.enqueue(p)
				}
				p = _p
			}
		}
	}
}

func (p *IndexModule) indexWorker(b *Service) {
	var (
		e   *cacheEntry
		l   *list.ListHead
		now int64
	)
	ticker := falconTicker(time.Second/INDEX_QPS, b.Conf.Debug)
	q0 := &b.cache.idx0q
	q1 := &b.cache.idx1q
	q2 := &b.cache.idx2q

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			statsInc(ST_INDEX_TICK, 1)

			l = q0.dequeue()
			if l != nil {
				// immediate update , enqueue q1
				e = list_idx_entry(l)
				goto out
			}

			if l = q1.dequeue(); l == nil {
				continue
			}

			e = list_idx_entry(l)
			now = b.timeNow()

			if now-e.idxTs < INDEX_UPDATE_CYCLE_TIME {
				q1.addHead(l)
				continue
			}

			if now-e.lastTs > INDEX_TIMEOUT {
				// timeout entry move to q2
				q2.enqueue(l)
				continue
			}
		out:
			e.idxTs = now
			q1.enqueue(l)
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
			go indexUpdate(e, p.indexDb)
		}
	}
}

// indexModule depend on cacheModule
func (p *IndexModule) prestart(b *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.indexUpdateCh = make(chan *cacheEntry, INDEX_MAX_OPEN_CONNS)
	return nil
}

func (p *IndexModule) start(b *Service) error {
	var err error

	dbmaxidle, _ := b.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := b.Conf.Configer.Int(C_DB_MAX_CONN)

	p.indexDb, err = sql.Open("mysql", b.Conf.Configer.Str(C_DSN))
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.indexDb.SetMaxIdleConns(dbmaxidle)
	p.indexDb.SetMaxOpenConns(dbmaxconn)

	err = p.indexDb.Ping()
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	go p.indexWorker(b)
	go p.indexTrashWorker(b)
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