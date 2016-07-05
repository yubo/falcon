/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"sync/atomic"
)

const (
	ST_SUCCESS = iota
	ST_CONN_ERR
	ST_CONN_DIAL
	ST_PUT_SUCCESS
	ST_PUT_ERR
	ST_STAT_SIZE
)

var (
	statName []string = []string{
		"SUCCESS",
		"conn error",
		"conn dail",
		"put success",
		"put err",
	}
)

var (
	statCnt [ST_STAT_SIZE]uint64
)

func statHandle() (ret string) {
	for i := 0; i < ST_STAT_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n", statName[i],
			atomic.LoadUint64(&statCnt[i]))
	}
	return ret
}

func statInc(idx, n int) {
	atomic.AddUint64(&statCnt[idx], uint64(n))
}

func statSet(idx, n int) {
	atomic.StoreUint64(&statCnt[idx], uint64(n))
}

func statGet(idx int) uint64 {
	return atomic.LoadUint64(&statCnt[idx])
}
