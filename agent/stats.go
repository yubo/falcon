/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent/config"
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
)

const (
	ST_TX_PUT_ITERS = iota
	ST_TX_PUT_ITEMS
	ST_TX_PUT_ERR_ITERS
	ST_TX_PUT_ERR_ITEMS
	ST_RX_PUT_FMT_ERR
	ST_ARRAY_SIZE
)

var (
	statsCounter      []uint64
	statsCounterName  [][]byte
	statsCounterName_ [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
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

func (p *Agent) Stats(conf interface{}) (s string, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(200)*time.Millisecond)
	conn, _, err := falcon.DialRr(ctx, conf.(*config.Agent).Configer.Str(C_API_ADDR), false)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := transfer.NewTransferClient(conn)
	stats, err := client.GetStats(context.Background(), &transfer.Empty{})
	if err != nil {
		return "", err
	}

	for i := 0; i < ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n", statsCounterName_[i], stats.Counter[i])
	}
	return s, nil
}
