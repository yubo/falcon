/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"fmt"
	"sync/atomic"
)

const (
	ST_FETCH_SUCCESS = iota
	ST_FETCH_ERR
	ST_FETCH_ISNOTEXIST
	ST_SEND_SUCCESS
	ST_SEND_ERR
	ST_QUERY_SUCCESS
	ST_QUERY_ERR
	ST_CONN_ERR
	ST_CONN_DIAL
	ST_DISK_TASK_CNT
	ST_NET_TASK_CNT
	ST_COUNTER_CACHE_CNT
	ST_STORAGE_RPC_RECV_CNT // Rpc
	ST_STORAGE_QUERY_CNT    // Query
	ST_STORAGE_QUERY_ITEM_CNT
	ST_STORAGE_INFO_CNT
	ST_STORAGE_LAST_CNT
	ST_STORAGE_LAST_RAW_CNT
	ST_STORAGE_LOAD_DB_CNT // load sth from db when query/info, tmp
	ST_STAT_SIZE
)

var (
	statName []string = []string{
		"FETCH_SUCCESS",
		"FETCH_ERR",
		"FETCH_ISNOTEXIST",
		"SEND_SUCCESS",
		"SEND_ERR",
		"QUERY_SUCCESS",
		"QUERY_ERR",
		"CONN_ERR",
		"CONN_DIAL",
		"DISK_TASK_CNT",
		"NET_TASK_CNT",
		"COUNTER_CACHE_CNT",
		"STORAGE_RPC_RECV_CNT",
		"STORAGE_QUERY_CNT",
		"STORAGE_QUERY_ITEM_CNT",
		"STORAGE_INFO_CNT",
		"STORAGE_LAST_CNT",
		"STORAGE_LAST_RAW_CNT",
		"STORAGE_LOAD_DB_CNT",
	}
)

var (
	statCnt [ST_STAT_SIZE]uint64
)

func statHandle() (ret string) {
	for i := 0; i < ST_STAT_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n",
			statName[i], atomic.LoadUint64(&statCnt[i]))
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
