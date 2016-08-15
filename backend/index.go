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
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/list"
)

var (
	indexConfig   BackendOpts
	indexDb       *sql.DB
	indexUpdateCh chan *cacheEntry
)

func indexUpdate(e *cacheEntry) {
	var (
		err           error
		hid, cid, tid int64
		dstype        string
		step          int
	)

	statInc(ST_INDEX_UPDATE, 1)
	hid = -1
	err = indexDb.QueryRow("SELECT id FROM host WHERE host = ?", e.host).Scan(&hid)
	if err != nil {
		statInc(ST_INDEX_HOST_MISS, 1)
		if err == sql.ErrNoRows || hid < 0 {
			statInc(ST_INDEX_HOST_INSERT, 1)
			ret, err := indexDb.Exec("INSERT INTO host(host, ts, t_create) "+
				"VALUES (?, ?, now()) ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
				e.host(), int64(e.e.lastTs))
			if err != nil {
				statInc(ST_INDEX_HOST_INSERT_ERR, 1)
				glog.Warning(err)
				return
			}

			hid, err = ret.LastInsertId()
			if err != nil {
				glog.Warning(err)
				return
			}
		} else {
			glog.Warning(e.host(), err)
			return
		}
	}

	// tag
	tags := strings.Split(e.tags(), ",")
	for _, tag := range tags {

		tid = -1
		err := indexDb.QueryRow("SELECT id FROM tag WHERE tag = ? and "+
			"host_id = ?", tag, hid).Scan(&tid)
		if err != nil {
			statInc(ST_INDEX_TAG_MISS, 1)
			if err == sql.ErrNoRows || tid < 0 {
				statInc(ST_INDEX_TAG_INSERT, 1)
				ret, err := indexDb.Exec("INSERT INTO tag(tag, host_id, "+
					"ts, t_create) "+
					"VALUES (?, ?, ?, now()) "+
					"ON DUPLICATE KEY "+
					"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
					tag, hid, int64(e.e.lastTs))
				if err != nil {
					statInc(ST_INDEX_TAG_INSERT_ERR, 1)
					glog.Warning(err)
					return
				}

				tid, err = ret.LastInsertId()
				if err != nil {
					glog.Warning(err)
					return
				}
			} else {
				glog.Warning(tag, hid, err)
				return
			}
		}

	}

	// host_id
	counter := e.id()

	cid = -1
	step = 0
	dstype = "nil"

	err = indexDb.QueryRow("SELECT id,step,type FROM counter WHERE "+
		"host_id = ? and counter = ?",
		hid, counter).Scan(&cid, &step, &dstype)
	if err != nil {
		statInc(ST_INDEX_COUNTER_MISS, 1)
		if err == sql.ErrNoRows || cid < 0 {
			statInc(ST_INDEX_COUNTER_INSERT, 1)
			ret, err := indexDb.Exec("INSERT INTO counter(host_id,counter,"+
				"step,type,ts,t_create) "+
				"VALUES (?,?,?,?,?,now()) "+
				"ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts),"+
				"step=VALUES(step),type=VALUES(type)",
				hid, counter, int(e.e.step), e.typ(),
				int64(e.e.lastTs))
			if err != nil {
				statInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
				glog.Warning(err)
				return
			}

			cid, err = ret.LastInsertId()
			if err != nil {
				glog.Warning(err)
				return
			}
		} else {
			glog.Warning(counter, hid, err)
			return
		}
	} /* else {
		if !(e.step == step && e.typ == dstype) {
			_, err := indexDb.Exec("UPDATE counter SET step = ?, "+
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

func indexTrashWorker() {
	var (
		e     *cacheEntry
		p, _p *list.ListHead
	)
	ticker := falconTicker(time.Second*INDEX_TRASH_LOOPTIME,
		indexConfig.Debug)
	q0 := &appCache.idx0q
	q2 := &appCache.idx2q
	for {
		select {
		case <-ticker:
			for p = q2.head.Next; p != &q2.head; p = p.Next {
				_p = p.Next
				e = list_idx_entry(p)
				if timeNow()-int64(e.e.lastTs) < INDEX_TIMEOUT {
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

func indexWorker() {
	var (
		e *cacheEntry
		p *list.ListHead
	)
	ticker := falconTicker(time.Second/INDEX_QPS, indexConfig.Debug)
	q0 := &appCache.idx0q
	q1 := &appCache.idx1q
	q2 := &appCache.idx2q

	for {
		select {
		case <-ticker:
			statInc(ST_INDEX_TICK, 1)

			p = q0.dequeue()
			if p != nil {
				// immediate update , enqueue q1
				e = list_idx_entry(p)
				goto out
			}

			if p = q1.dequeue(); p == nil {
				continue
			}

			e = list_idx_entry(p)

			if timeNow()-int64(e.e.idxTs) < INDEX_UPDATE_CYCLE_TIME {
				q1.addHead(p)
				continue
			}

			if timeNow()-int64(e.e.lastTs) > INDEX_TIMEOUT {
				// timeout entry move to q2
				q2.enqueue(p)
				continue
			}
		out:
			e.e.idxTs = C.int64_t(timeNow())
			q1.enqueue(p)
			indexUpdateCh <- e
		}
	}
}

func updateWorker(conn *sql.DB) {
	for {
		select {
		case e := <-indexUpdateCh:
			go indexUpdate(e)
		}
	}
}

func indexStart(config BackendOpts, p *specs.Process) {
	var err error

	indexConfig = config

	indexDb, err = sql.Open("mysql", storageConfig.Dsn)
	if err != nil {
		glog.Fatal(err)
	}

	indexDb.SetMaxIdleConns(storageConfig.DbMaxIdle)
	indexDb.SetMaxOpenConns(0)

	err = indexDb.Ping()
	if err != nil {
		glog.Fatal(err)
	}

	indexUpdateCh = make(chan *cacheEntry, INDEX_MAX_OPEN_CONNS)

	go indexWorker()
	go indexTrashWorker()
	go updateWorker(indexDb)

	glog.Info("indexStart ok")
}
