/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm/config"
	"golang.org/x/net/context"
)

const (
	ST_RX_PUT_ITER = iota
	ST_RX_PUT_EVENT
	ST_RX_PUT_ERR_EVENT
	ST_PROCESS_EVENT
	ST_PROCESS_EVENT_ERR
	ST_ACTIONCHAN_IN
	ST_ACTIONCHAN_IN_ERR
	ST_ACTIONCHAN_OUT

	ST_EVENTENTRY
	ST_EVENTENTRY_EXPIRED

	ST_TX_PUT_ITER
	ST_TX_PUT_ITEMS
	ST_TX_PUT_ERR_ITERS
	ST_TX_PUT_ERR_ITEMS
	ST_ARRAY_SIZE
)

var (
	statsCounter      []uint64
	statsCounterName  [][]byte
	statsCounterName_ [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_rx_put_iters",
		"st_rx_put_event",
		"st_rx_put_event_err",
		"st_process_event",
		"st_process_event_err",
		"st_actionchan_in",
		"st_actionchan_in_err",
		"st_actionchan_out",
		"st_evententry",
		"st_evententry_expired",
		"st_tx_put_iters",
		"st_tx_put_items",
		"st_tx_put_err_iters",
		"st_tx_put_err_items",
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

func (p *Alarm) Stats(conf interface{}) (s string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(200)*time.Millisecond)
	conn, _, err := falcon.DialRr(ctx, conf.(*config.Alarm).Configer.Str(C_API_ADDR), false)
	if err != nil {
		return ""
	}
	defer conn.Close()

	client := NewAlarmClient(conn)
	stats, err := client.GetStats(context.Background(), &Empty{})
	if err != nil {
		return ""
	}

	for i := 0; i < ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n", statsCounterName_[i], stats.Counter[i])
	}
	return
}
