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
	ST_RPC_UPDATE = iota
	ST_RPC_UPDATE_CNT
	ST_RPC_UPDATE_ERR
	ST_RPC_DROP_CNT
	ST_UPSTREAM_RECONNECT
	ST_UPSTREAM_DIAL
	ST_UPSTREAM_DIAL_ERR
	ST_UPSTREAM_PUT
	ST_UPSTREAM_PUT_ITEM_TOTAL
	ST_UPSTREAM_PUT_ITEM_ERR
	ST_UPSTREAM_PUT_ERR
	ST_ARRAY_SIZE
)

var (
	counter     [ST_ARRAY_SIZE]uint64
	counterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"ST_RPC_UPDATE",
		"ST_RPC_UPDATE_CNT",
		"ST_RPC_UPDATE_ERR",
		"ST_UPSTREAM_RECONNECT",
		"ST_UPSTREAM_DIAL",
		"ST_UPSTREAM_DIAL_ERR",
		"ST_UPSTREAM_PUT",
		"ST_UPSTREAM_PUT_ITEM",
		"ST_UPSTREAM_PUT_ERR",
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
