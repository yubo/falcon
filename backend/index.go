/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
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
	"github.com/yubo/gotool/list"
)

type IndexModule struct {
	indexUpdateCh chan *cacheEntry
	indexDb       *sql.DB
	running       chan struct{}
}

func indexUpdate(e *cacheEntry, db *sql.DB) {
	var (
		err           error
		hid, cid, tid int64
		dstype        string
		step          int
	)

	statsInc(ST_INDEX_UPDATE, 1)
	hid = -1
	err = db.QueryRow("SELECT id FROM host WHERE host = ?", e.host).Scan(&hid)
	if err != nil {
		statsInc(ST_INDEX_HOST_MISS, 1)
		if err == sql.ErrNoRows || hid < 0 {
			statsInc(ST_INDEX_HOST_INSERT, 1)
			ret, err := db.Exec("INSERT INTO host(host, ts, t_create) "+
				"VALUES (?, ?, now()) ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id), ts=VALUES(ts)",
				e.host, e.lastTs)
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
			glog.Warning(MODULE_NAME+string(e.host), err)
			return
		}
	}

	// tag
	tags := strings.Split(string(e.tags), ",")
	for _, tag := range tags {

		tid = -1
		err := db.QueryRow("SELECT id FROM tag WHERE tag = ? and "+
			"host_id = ?", tag, hid).Scan(&tid)
		if err != nil {
			statsInc(ST_INDEX_TAG_MISS, 1)
			if err == sql.ErrNoRows || tid < 0 {
				statsInc(ST_INDEX_TAG_INSERT, 1)
				ret, err := db.Exec("INSERT INTO tag(tag, host_id, "+
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

	// host_id
	counter := e.id()

	cid = -1
	step = 0
	dstype = "nil"

	err = db.QueryRow("SELECT id,step,type FROM counter WHERE "+
		"host_id = ? and counter = ?",
		hid, counter).Scan(&cid, &step, &dstype)
	if err != nil {
		statsInc(ST_INDEX_COUNTER_MISS, 1)
		if err == sql.ErrNoRows || cid < 0 {
			statsInc(ST_INDEX_COUNTER_INSERT, 1)
			ret, err := db.Exec("INSERT INTO counter(host_id,counter,"+
				"step,type,ts,t_create) "+
				"VALUES (?,?,?,?,?,now()) "+
				"ON DUPLICATE KEY "+
				"UPDATE id=LAST_INSERT_ID(id),ts=VALUES(ts),"+
				"step=VALUES(step),type=VALUES(type)",
				hid, counter, e.step, e.typ,
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

func (this *IndexModule) indexTrashWorker(b *Backend) {
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
		case _, ok := <-this.running:
			if !ok {
				return
			}
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

func (p *IndexModule) indexWorker(b *Backend) {
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
		case _, ok := <-p.running:
			if !ok {
				return
			}
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
		case _, ok := <-p.running:
			if !ok {
				return
			}
		case e := <-p.indexUpdateCh:
			go indexUpdate(e, p.indexDb)
		}
	}
}

// indexModule depend on cacheModule
func (p *IndexModule) prestart(b *Backend) error {
	p.running = make(chan struct{}, 0)
	p.indexUpdateCh = make(chan *cacheEntry, INDEX_MAX_OPEN_CONNS)
	return nil
}

func (p *IndexModule) start(b *Backend) error {
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

func (p *IndexModule) stop(b *Backend) error {
	close(p.running)
	return nil
}

func (p *IndexModule) reload(b *Backend) error {
	return nil
}
