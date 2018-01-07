/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	ST_RX_PUT_ITERS = iota
	ST_RX_PUT_ITEMS
	ST_RX_PUT_ERR_ITEMS
	ST_TX_PUT_ITERS
	ST_TX_PUT_ITEMS
	ST_TX_PUT_ERR_ITERS
	ST_TX_PUT_ERR_ITEMS
	ST_ARRAY_SIZE
)

var (
	counter     [ST_ARRAY_SIZE]uint64
	counterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"st_rx_put_iters",
		"st_rx_put_items",
		"st_rx_put_err_items",
		"st_tx_put_iters",
		"st_tx_put_items",
		"st_tx_put_err_iters",
		"st_tx_put_err_items",
	}
)

func statsInc(idx, n int) {
	atomic.AddUint64(&counter[idx], uint64(n))
}

func statsSet(idx, n int) {
	atomic.StoreUint64(&counter[idx], uint64(n))
}

func statsGet(idx int) uint64 {
	return atomic.LoadUint64(&counter[idx])
}

func statsHandle() (ret string) {
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		ret += fmt.Sprintf("%s %d\n", counterName[i],
			atomic.LoadUint64(&counter[i]))
	}
	return ret
}

type StatsModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *StatsModule) prestart(t *Transfer) error {
	return nil
}

func (p *StatsModule) start(t *Transfer) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	if t.Conf.Debug > 0 {
		statsTicker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case <-statsTicker:
					glog.V(3).Info(MODULE_NAME + statsHandle())
				}
			}
		}()
	}
	return nil
}

func (p *StatsModule) stop(t *Transfer) error {
	p.cancel()
	return nil
}

func (p *StatsModule) reload(t *Transfer) error {
	p.stop(t)
	time.Sleep(time.Second)
	p.prestart(t)
	return p.start(t)
}
