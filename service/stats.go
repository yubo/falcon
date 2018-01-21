/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service/config"
	"golang.org/x/net/context"
)

const (
	ST_RX_PUT_ITERS = iota
	ST_RX_PUT_ITEMS
	ST_RX_PUT_ERR_ITEMS
	ST_RX_GET_ITERS
	ST_RX_GET_ITEMS

	ST_TX_PUT_ITERS
	ST_TX_PUT_ITEMS
	ST_TX_PUT_ERR_ITERS
	ST_TX_PUT_ERR_ITEMS

	ST_INDEX_HOST_INSERT_ERR
	ST_INDEX_TAG_INSERT_ERR
	ST_INDEX_COUNTER_INSERT_ERR
	ST_INDEX_TICK
	ST_INDEX_UPDATE

	ST_JUDGE_CNT
	ST_EVENT_CNT
	ST_EVENT_DROP_CNT

	ST_ARRAY_SIZE
)

var (
	statsCounter      []uint64
	statsCounterName  [][]byte
	statsCounterName_ [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_rx_put_iters",
		"st_rx_put_items",
		"st_rx_put_err_items",
		"st_rx_get_iters",
		"st_rx_get_items",
		"st_tx_put_iters",
		"st_tx_put_items",
		"st_tx_put_err_iters",
		"st_tx_put_err_items",
		"st_index_host_insert_err",
		"st_index_tag_insert_err",
		"st_index_counter_insert_err",
		"st_index_tick",
		"st_index_update",
		"st_judge_cnt",
		"st_event_cnt",
		"st_event_drop_cnt",
	}
)

func init() {
	statsCounter = make([]uint64, ST_ARRAY_SIZE)
	statsCounterName = make([][]byte, ST_ARRAY_SIZE)
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		statsCounterName[i] = []byte(statsCounterName_[i])
	}
}

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

func statsGets() []uint64 {
	cnt := make([]uint64, ST_ARRAY_SIZE)
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		cnt[i] = atomic.LoadUint64(&statsCounter[i])
	}
	return cnt
}

func (p *Service) Stats(conf interface{}) (s string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(200)*time.Millisecond)
	conn, _, err := falcon.DialRr(ctx, conf.(*config.Service).Configer.Str(C_API_ADDR), false)
	if err != nil {
		return ""
	}
	defer conn.Close()

	client := NewServiceClient(conn)
	stats, err := client.GetStats(context.Background(), &Empty{})
	if err != nil {
		return ""
	}

	for i := 0; i < ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n", statsCounterName_[i], stats.Counter[i])
	}
	return
}
