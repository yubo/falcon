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
	ST_GET_ITERS = iota
	ST_GET_ITERS_ERR
	ST_GET_DPS
	ST_GET_DPS_ERR
	ST_ARRAY_SIZE
)

var (
	statsCounter     [ST_ARRAY_SIZE]uint64
	statsCounterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_get_iters",
		"st_get_iters_err",
		"st_get_dps",
		"st_get_dps_err",
	}
)

func statsDec(idx, n int) {
	atomic.AddUint64(&statsCounter[idx], ^uint64(n-1))
}

func statsInc(idx, n int) {
	atomic.AddUint64(&statsCounter[idx], uint64(n))
}

func statsSet(idx, n int) {
	atomic.StoreUint64(&statsCounter[idx], uint64(n))
}

func statsGet(idx int) uint64 {
	return atomic.LoadUint64(&statsCounter[idx])
}

func StatsGets() []uint64 {
	cnt := make([]uint64, ST_ARRAY_SIZE)
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		cnt[i] = atomic.LoadUint64(&statsCounter[i])
	}
	return cnt
}

func (p *Ctrl) Stats(conf interface{}) (s string, err error) {
	// http api
	var stats []uint64

	url := conf.(*config.Ctrl).Ctrl.Str(C_HTTP_ADDR)
	if strings.HasPrefix(url, ":") {
		url = "http://localhost" + url
	}
	url += "/v1.0/pub/stats"

	if err := falcon.GetJson(url, &stats, time.Duration(200)*time.Millisecond); err != nil {
		return "", err
	}

	for i := 0; i < ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n",
			statsCounterName[i], stats[i])
	}
	return s, nil
}
