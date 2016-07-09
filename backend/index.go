/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

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
				e.host, e.lastTs)
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
			glog.Warning(e.host, err)
			return
		}
	}

	// tag
	tags := strings.Split(e.tags, ",")
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
					tag, hid, e.lastTs)
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
	counter := e.Id()

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
				hid, counter, e.step, e.typ, e.lastTs)
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

func index1Worker() {
	var (
		entry *cacheEntry
		p, _p *list.ListHead
	)
	ticker := falconTicker(time.Second*INDEX_TRASH_LOOPTIME,
		indexConfig.Debug)
	q0 := &appCache.idx0q
	q1 := &appCache.idx1q
	for {
		select {
		case <-ticker:
			for p = q1.head.Next; p != &q1.head; p = p.Next {
				_p = p.Next
				entry = list_idx_entry(p)
				if timeNow()-entry.lastTs < INDEX_TIMEOUT {
					q1.Lock()
					p.Del()
					q1.size--
					q1.Unlock()

					entry.idxTs = 0
					statInc(ST_INDEX_TRASH_PICKUP, 1)
					q0.addHead(p)
				}
				p = _p
			}
		}
	}
}

func index0Worker() {
	var (
		entry   *cacheEntry
		pending *list.ListHead
		p       *list.ListHead
	)
	ticker := falconTicker(time.Second/INDEX_QPS, indexConfig.Debug)
	q0 := &appCache.idx0q
	q1 := &appCache.idx1q
	pending = &q0.head

	for {
		select {
		case <-ticker:
			statInc(ST_INDEX_TICK, 1)
			if q0.size == 0 {
				continue
			}

			q0.Lock()
			if pending == &q0.head {
				for p = q0.head.Next; p != &q0.head; p = p.Next {
					entry = list_idx_entry(p)
					if entry.idxTs == 0 {
						pending = p
					} else {
						break
					}
				}
			}

			if pending != &q0.head {
				p = pending
				pending = p.Prev
			}

			entry = list_idx_entry(p)

			if timeNow()-entry.idxTs > INDEX_UPDATE_CYCLE_TIME {
				p.Del()
				q0.size--
				q0.Unlock()
			} else {
				q0.Unlock()
				continue
			}

			if timeNow()-entry.lastTs > INDEX_TIMEOUT {
				statInc(ST_INDEX_TIMEOUT, 1)
				q1.enqueue(p)
				continue
			}

			entry.idxTs = timeNow()
			q0.enqueue(p)
			indexUpdateCh <- entry
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

	indexDb, err = sql.Open("mysql", rrdConfig.Dsn)
	if err != nil {
		glog.Fatal(err)
	}

	indexDb.SetMaxIdleConns(rrdConfig.DbMaxIdle)
	indexDb.SetMaxOpenConns(0)

	err = indexDb.Ping()
	if err != nil {
		glog.Fatal(err)
	}

	indexUpdateCh = make(chan *cacheEntry, INDEX_MAX_OPEN_CONNS)

	go index0Worker()
	go index1Worker()
	go updateWorker(indexDb)

	glog.Info("indexStart ok")
}
