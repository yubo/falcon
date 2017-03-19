/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

type statsIdx struct {
	begin, end int
}

const (
	ST_RPC_SERV_QUERY = iota
	ST_RPC_SERV_QUERY_ITEM
	ST_RPC_SERV_RECV
	ST_RPC_SERV_RECV_ITEM
	ST_RPC_SERV_GETRRD
	ST_RPC_SERV_GETRRD_ERR

	ST_RRD_CREAT
	ST_RRD_CREAT_ERR
	ST_RRD_UPDATE
	ST_RRD_UPDATE_ERR
	ST_RRD_FETCH
	ST_RRD_FETCH_ERR

	ST_RPC_CLI_SEND
	ST_RPC_CLI_SEND_ERR
	ST_RPC_CLI_QUERY
	ST_RPC_CLI_QUERY_ERR
	ST_RPC_CLI_FETCH
	ST_RPC_CLI_FETCH_ERR
	ST_RPC_CLI_FETCH_ERR_NOEXIST
	ST_RPC_CLI_RECONNECT
	ST_RPC_CLI_DIAL
	ST_RPC_CLI_DIAL_ERR

	ST_CACHE_MISS
	ST_CACHE_CREATE
	ST_CACHE_OVERRUN

	ST_INDEX_HOST_MISS
	ST_INDEX_HOST_INSERT
	ST_INDEX_HOST_INSERT_ERR
	ST_INDEX_TAG_MISS
	ST_INDEX_TAG_INSERT
	ST_INDEX_TAG_INSERT_ERR
	ST_INDEX_COUNTER_MISS
	ST_INDEX_COUNTER_INSERT
	ST_INDEX_COUNTER_INSERT_ERR
	ST_INDEX_TICK
	ST_INDEX_TIMEOUT
	ST_INDEX_TRASH_PICKUP
	ST_INDEX_UPDATE

	ST_ARRAY_SIZE
)

const (
	ST_M_RPC_SERV uint32 = 1 << iota
	ST_M_RRD
	ST_M_RPC_CLI
	ST_M_CACHE
	ST_M_INDEX
	ST_M_SIZE int = iota
)

const (
	DEBUG_STAT_MODULE = ST_M_CACHE | ST_M_INDEX
)

var (
	stat_module [ST_M_SIZE]statsIdx = [ST_M_SIZE]statsIdx{
		{
			ST_RPC_SERV_QUERY,
			ST_RRD_CREAT,
		},
		{
			ST_RRD_CREAT,
			ST_RPC_CLI_SEND,
		},
		{
			ST_RPC_CLI_SEND,
			ST_CACHE_MISS,
		},
		{
			ST_CACHE_MISS,
			ST_INDEX_HOST_MISS,
		},
		{
			ST_INDEX_HOST_MISS,
			ST_ARRAY_SIZE,
		},
	}

	statName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"ST_RPC_SERV_QUERY",
		"ST_RPC_SERV_QUERY_ITEM",
		"ST_RPC_SERV_RECV",
		"ST_RPC_SERV_RECV_ITEM",
		"ST_RPC_SERV_GETRRD",
		"ST_RPC_SERV_GETRRD_ERR",
		"ST_RRD_CREAT",
		"ST_RRD_CREAT_ERR",
		"ST_RRD_UPDATE",
		"ST_RRD_UPDATE_ERR",
		"ST_RRD_FETCH",
		"ST_RRD_FETCH_ERR",
		"ST_RPC_CLI_SEND",
		"ST_RPC_CLI_SEND_ERR",
		"ST_RPC_CLI_QUERY",
		"ST_RPC_CLI_QUERY_ERR",
		"ST_RPC_CLI_FETCH",
		"ST_RPC_CLI_FETCH_ERR",
		"ST_RPC_CLI_FETCH_ERR_NOEXIST",
		"ST_RPC_CLI_RECONNECT",
		"ST_RPC_CLI_DIAL",
		"ST_RPC_CLI_DIAL_ERR",
		"ST_CACHE_MISS",
		"ST_CACHE_CREATE",
		"ST_CACHE_OVERRUN",
		"ST_INDEX_HOST_MISS",
		"ST_INDEX_HOST_INSERT",
		"ST_INDEX_HOST_INSERT_ERR",
		"ST_INDEX_TAG_MISS",
		"ST_INDEX_TAG_INSERT",
		"ST_INDEX_TAG_INSERT_ERR",
		"ST_INDEX_COUNTER_MISS",
		"ST_INDEX_COUNTER_INSERT",
		"ST_INDEX_COUNTER_INSERT_ERR",
		"ST_INDEX_TICK",
		"ST_INDEX_TIMEOUT",
		"ST_INDEX_PICKUP",
		"ST_INDEX_UPDATE",
	}
)

var (
	statCnt [ST_ARRAY_SIZE]uint64
)

func statsHandle() (ret string) {
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n",
			statName[i], atomic.LoadUint64(&statCnt[i]))
	}
	return ret
}

func statsModuleHandle(m uint32) (ret string) {
	for i := 0; i < ST_M_SIZE; i++ {
		if m&1 == 1 {
			for j := stat_module[i].begin; j < stat_module[i].end; j++ {
				ret += fmt.Sprintf("%s %d\n",
					statName[j], atomic.LoadUint64(&statCnt[j]))
			}
		}
		m = m >> 1
	}
	return ret
}

func statsInc(idx, n int) {
	atomic.AddUint64(&statCnt[idx], uint64(n))
}

func statsSet(idx, n int) {
	atomic.StoreUint64(&statCnt[idx], uint64(n))
}

func statsGet(idx int) uint64 {
	return atomic.LoadUint64(&statCnt[idx])
}

func statsRrd() {
	for i := ST_RRD_CREAT; i <= ST_RRD_FETCH_ERR; i++ {
		glog.V(3).Infof("%s %d", statName[i], statCnt[i])
	}
}

type statsModule struct {
	running chan struct{}
}

func (p *statsModule) prestart(b *Backend) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *statsModule) start(b *Backend) error {
	if b.Conf.Debug > 0 {
		ticker := falconTicker(time.Second*DEBUG_STAT_STEP, b.Conf.Debug)
		go func() {
			for {
				select {
				case _, ok := <-p.running:
					if !ok {
						return
					}
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statsModuleHandle(DEBUG_STAT_MODULE))
				}
			}
		}()
	}
	return nil
}

func (p *statsModule) stop(b *Backend) error {
	close(p.running)
	return nil
}

func (p *statsModule) reload(b *Backend) error {
	return nil
}
