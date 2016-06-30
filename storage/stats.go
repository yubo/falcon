/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"fmt"
	"sync/atomic"
)

const (
	ST_FETCH_S_SUCCESS = iota
	ST_FETCH_S_ERR
	ST_FETCH_S_ISNOTEXIST
	ST_SEND_S_SUCCESS
	ST_SEND_S_ERR
	ST_QUERY_S_SUCCESS
	ST_QUERY_S_ERR
	ST_CONN_S_ERR
	ST_CONN_S_DIAL
	ST_DISK_COUNTER
	ST_NET_COUNTER
	ST_IndexUpdateIncr
	ST_IndexUpdateIncrCnt
	ST_IndexUpdateIncrErrorCnt
	ST_IndexUpdateIncrDbEndpointSelectCnt
	ST_IndexUpdateIncrDbEndpointInsertCnt
	ST_IndexUpdateIncrDbTagEndpointSelectCnt
	ST_IndexUpdateIncrDbTagEndpointInsertCnt
	ST_IndexUpdateIncrDbEndpointCounterSelectCnt
	ST_IndexUpdateIncrDbEndpointCounterInsertCnt
	ST_IndexUpdateIncrDbEndpointCounterUpdateCnt
	ST_IndexUpdateAll // 索引全量更新
	ST_IndexUpdateAllCnt
	ST_IndexUpdateAllErrorCnt
	ST_IndexedItemCacheCnt // 索引缓存大小
	ST_UnIndexedItemCacheCnt
	ST_EndpointCacheCnt
	ST_CounterCacheCnt
	ST_GraphRpcRecvCnt // Rpc
	ST_GraphQueryCnt   // Query
	ST_GraphQueryItemCnt
	ST_GraphInfoCnt
	ST_GraphLastCnt
	ST_GraphLastRawCnt
	ST_GraphLoadDbCnt // load sth from db when query/info, tmp
	ST_STAT_SIZE
)

var (
	stat_name []string = []string{
		"FETCH_S_SUCCESS",
		"FETCH_S_ERR",
		"FETCH_S_ISNOTEXIST",
		"SEND_S_SUCCESS",
		"SEND_S_ERR",
		"QUERY_S_SUCCESS",
		"QUERY_S_ERR",
		"CONN_S_ERR",
		"CONN_S_DIAL",
		"DISK_COUNTER",
		"NET_COUNTER",
		"IndexUpdateIncr",
		"IndexUpdateIncrCnt",
		"IndexUpdateIncrErrorCnt",
		"IndexUpdateIncrDbEndpointSelectCnt",
		"IndexUpdateIncrDbEndpointInsertCnt",
		"IndexUpdateIncrDbTagEndpointSelectCnt",
		"IndexUpdateIncrDbTagEndpointInsertCnt",
		"IndexUpdateIncrDbEndpointCounterSelectCnt",
		"IndexUpdateIncrDbEndpointCounterInsertCnt",
		"IndexUpdateIncrDbEndpointCounterUpdateCnt",
		"IndexUpdateAll",
		"IndexUpdateAllCnt",
		"IndexUpdateAllErrorCnt",
		"IndexedItemCacheCnt",
		"UnIndexedItemCacheCnt",
		"EndpointCacheCnt",
		"CounterCacheCnt",
		"GraphRpcRecvCnt",
		"GraphQueryCnt",
		"GraphQueryItemCnt",
		"GraphInfoCnt",
		"GraphLastCnt",
		"GraphLastRawCnt",
		"GraphLoadDbCnt",
	}
)

var (
	stat_cnt [ST_STAT_SIZE]uint64
)

func stat_handle() (ret string) {
	for i := 0; i < ST_STAT_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n", stat_name[i], atomic.LoadUint64(&stat_cnt[i]))
	}
	return ret
}

func stat_inc(idx, n int) {
	atomic.AddUint64(&stat_cnt[idx], uint64(n))
}

func stat_set(idx, n int) {
	atomic.StoreUint64(&stat_cnt[idx], uint64(n))
}
