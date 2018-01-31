/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

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

func indexUpdateEndpointCounter(metric, tags, typ string, ts int64, db orm.Ormer, hid int64) (err error) {

	_, err = db.Raw("INSERT INTO endpoint_counter(endpoint_id,counter,ts,t_create) VALUES (?,?,?,now()) ON DUPLICATE KEY UPDATE ts=?, t_modify=now()", hid, metric+"/"+tags+"/"+typ, ts, ts).Exec()
	if err != nil {
		statsInc(ST_INDEX_COUNTER_INSERT_ERR, 1)
	}
	return
}

func indexUpdate(e *cacheEntry, db orm.Ormer) {

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
