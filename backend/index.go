/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

/*
#include "cache.h"
*/
import "C"
import (
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/gotool/list"
)

func indexUpdate(e *cacheEntry, b *Backend) {
	var (
		err           error
		hid, cid, tid int64
		dstype        string
		step          int
	)

	statInc(ST_INDEX_UPDATE, 1)
	hid = -1
	err = b.indexDb.QueryRow("SELECT id FROM host WHERE host = ?", e.host).Scan(&hid)
	if err != nil {
		statInc(ST_INDEX_HOST_MISS, 1)
		if err == sql.ErrNoRows || hid < 0 {
			statInc(ST_INDEX_HOST_INSERT, 1)
			ret, err := b.indexDb.Exec("INSERT INTO host(host, ts, t_create) "+
				"VALUES (?, ?, now()) ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
				e.host(), int64(e.e.lastTs))
			if err != nil {
				statInc(ST_INDEX_HOST_INSERT_ERR, 1)
				glog.Warning(MODULE_NAME, err)
				return
			}

			hid, err = ret.LastInsertId()
			if err != nil {
				glog.Warning(MODULE_NAME, err)
				return
			}
		} else {
			glog.Warning(MODULE_NAME+e.host(), err)
			return
		}
	}

	// tag
	tags := strings.Split(e.tags(), ",")
	for _, tag := range tags {

		tid = -1
		err := b.indexDb.QueryRow("SELECT id FROM tag WHERE tag = ? and "+
			"host_id = ?", tag, hid).Scan(&tid)
		if err != nil {
			statInc(ST_INDEX_TAG_MISS, 1)
			if err == sql.ErrNoRows || tid < 0 {
				statInc(ST_INDEX_TAG_INSERT, 1)
				ret, err := b.indexDb.Exec("INSERT INTO tag(tag, host_id, "+
					"ts, t_create) "+
					"VALUES (?, ?, ?, now()) "+
					"ON DUPLICATE KEY "+
					"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
					tag, hid, int64(e.e.lastTs))
				if err != nil {
					statInc(ST_INDEX_TAG_INSERT_ERR, 1)
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

	// host_id
	counter := e.id()

	cid = -1
	step = 0
	dstype = "nil"

	err = b.indexDb.QueryRow("SELECT id,step,type FROM counter WHERE "+
		"host_id = ? and counter = ?",
		hid, counter).Scan(&cid, &step, &dstype)
	if err != nil {
		statInc(ST_INDEX_COUNTER_MISS, 1)
		if err == sql.ErrNoRows || cid < 0 {
			statInc(ST_INDEX_COUNTER_INSERT, 1)
			ret, err := b.indexDb.Exec("INSERT INTO counter(host_id,counter,"+
				"step,type,ts,t_create) "+
				"VALUES (?,?,?,?,?,now()) "+
				"ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts),"+
				"step=VALUES(step),type=VALUES(type)",
				hid, counter, int(e.e.step), e.typ(),
				int64(e.e.lastTs))
			if err != nil {
				statInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
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
			_, err := b.indexDb.Exec("UPDATE counter SET step = ?, "+
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

func (this *Backend) indexTrashWorker() {
	var (
		e     *cacheEntry
		p, _p *list.ListHead
	)
	ticker := falconTicker(time.Second*INDEX_TRASH_LOOPTIME,
		this.Conf.Params.Debug)
	q0 := &this.cache.idx0q
	q2 := &this.cache.idx2q
	for {
		select {
		case <-ticker:
			for p = q2.head.Next; p != &q2.head; p = p.Next {
				_p = p.Next
				e = list_idx_entry(p)
				if this.timeNow()-int64(e.e.lastTs) < INDEX_TIMEOUT {
					q2.Lock()
					p.Del()
					//q2.size--
					q2.Unlock()

					e.e.idxTs = 0
					statInc(ST_INDEX_TRASH_PICKUP, 1)
					q0.enqueue(p)
				}
				p = _p
			}
		}
	}
}

func (p *Backend) indexWorker() {
	var (
		e   *cacheEntry
		l   *list.ListHead
		now int64
	)
	ticker := falconTicker(time.Second/INDEX_QPS, p.Conf.Params.Debug)
	q0 := &p.cache.idx0q
	q1 := &p.cache.idx1q
	q2 := &p.cache.idx2q

	for {
		select {
		case <-ticker:
			statInc(ST_INDEX_TICK, 1)

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
			now = p.timeNow()

			if now-int64(e.e.idxTs) < INDEX_UPDATE_CYCLE_TIME {
				q1.addHead(l)
				continue
			}

			if now-int64(e.e.lastTs) > INDEX_TIMEOUT {
				// timeout entry move to q2
				q2.enqueue(l)
				continue
			}
		out:
			e.e.idxTs = C.int64_t(now)
			q1.enqueue(l)
			p.indexUpdateCh <- e
		}
	}
}

func (p *Backend) updateWorker() {
	for {
		select {
		case e := <-p.indexUpdateCh:
			go indexUpdate(e, p)
		}
	}
}

func (p *Backend) indexStart() {
	var err error

	p.indexDb, err = sql.Open("mysql", p.Conf.Dsn)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.indexDb.SetMaxIdleConns(p.Conf.DbMaxIdle)
	p.indexDb.SetMaxOpenConns(0)

	err = p.indexDb.Ping()
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.indexUpdateCh = make(chan *cacheEntry, INDEX_MAX_OPEN_CONNS)

	go p.indexWorker()
	go p.indexTrashWorker()
	go p.updateWorker()

	glog.Info(MODULE_NAME, "indexStart ok")
}

func (p *Backend) indexStop() {
}
