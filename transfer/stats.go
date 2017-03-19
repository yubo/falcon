/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

const (
	ST_RPC_UPDATE = iota
	ST_RPC_UPDATE_CNT
	ST_RPC_UPDATE_ERR
	ST_UPSTREAM_RECONNECT
	ST_UPSTREAM_DIAL
	ST_UPSTREAM_DIAL_ERR
	ST_UPSTREAM_PUT
	ST_UPSTREAM_PUT_ITEM
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

type statsModule struct {
	running chan struct{}
}

func (p *statsModule) prestart(L *Transfer) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *statsModule) start(L *Transfer) error {
	if L.Conf.Debug > 0 {
		statsTicker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case _, ok := <-p.running:
					if !ok {
						return
					}
				case <-statsTicker:
					glog.V(3).Info(MODULE_NAME + statsHandle())
				}
			}
		}()
	}
	return nil
}

func (p *statsModule) stop(L *Transfer) error {
	close(p.running)
	return nil
}

func (p *statsModule) reload(L *Transfer) error {
	return nil
}
