/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/config"
)

const (
	ST_UPSTREAM_RECONNECT = iota
	ST_UPSTREAM_DIAL
	ST_UPSTREAM_DIAL_ERR
	ST_UPSTREAM_UPDATE
	ST_UPSTREAM_UPDATE_ITEM
	ST_UPSTREAM_UPDATE_ERR
	ST_ARRAY_SIZE
)

var (
	statsCounter     [ST_ARRAY_SIZE]uint64
	statsCounterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_upstream_reconnect",
		"st_upstream_dial",
		"st_upstream_dial_err",
		"st_upstream_update",
		"st_upstream_update_item",
		"st_upstream_update_err",
	}
)

func StatsHandle() (ret []uint64) {
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		ret = append(ret, atomic.LoadUint64(&statsCounter[i]))
	}
	return
}

func statInc(idx, n int) {
	atomic.AddUint64(&statsCounter[idx], uint64(n))
}

func statSet(idx, n int) {
	atomic.StoreUint64(&statsCounter[idx], uint64(n))
}

func statGet(idx int) uint64 {
	return atomic.LoadUint64(&statsCounter[idx])
}

func (p *Ctrl) Stats(conf interface{}) (s string) {
	// http api
	var stats []uint64

	url := conf.(*config.Ctrl).Ctrl.Str(C_HTTP_ADDR)
	if strings.HasPrefix(url, ":") {
		url = "http://localhost" + url
	}
	url += "/v1.0/pub/stats"

	if err := falcon.GetJson(url, &stats, time.Duration(200)*time.Millisecond); err != nil {
		return ""
	}

	for i := 0; i < ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n", statsCounterName[i], stats[i])
	}
	return
}
