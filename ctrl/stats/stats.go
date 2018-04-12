/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package stats

import (
	"sync/atomic"
)

const (
	ST_GET_ITERS = iota
	ST_GET_ITERS_ERR
	ST_GET_DPS
	ST_GET_DPS_ERR
	ST_ARRAY_SIZE
)

var (
	Counter     [ST_ARRAY_SIZE]uint64
	CounterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_get_iters",
		"st_get_iters_err",
		"st_get_dps",
		"st_get_dps_err",
	}
)

func Dec(idx, n int) {
	atomic.AddUint64(&Counter[idx], ^uint64(n-1))
}

func Inc(idx, n int) {
	atomic.AddUint64(&Counter[idx], uint64(n))
}

func Set(idx, n int) {
	atomic.StoreUint64(&Counter[idx], uint64(n))
}

func Get(idx int) uint64 {
	return atomic.LoadUint64(&Counter[idx])
}

func Gets() []uint64 {
	cnt := make([]uint64, ST_ARRAY_SIZE)
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		cnt[i] = atomic.LoadUint64(&Counter[i])
	}
	return cnt
}
