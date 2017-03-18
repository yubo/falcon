/*
 * Copyright 2016 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

const (
	ST_UPSTREAM_RECONNECT = iota
	ST_UPSTREAM_DIAL
	ST_UPSTREAM_DIAL_ERR
	ST_UPSTREAM_UPDATE
	ST_UPSTREAM_UPDATE_ITEM
	ST_UPSTREAM_UPDATE_ERR
	ST_ARRAY_SIZE
)

var (
	counter     [ST_ARRAY_SIZE]uint64
	counterName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"ST_UPSTREAM_RECONNECT",
		"ST_UPSTREAM_DIAL",
		"ST_UPSTREAM_DIAL_ERR",
		"ST_UPSTREAM_UPDATE",
		"ST_UPSTREAM_UPDATE_ITEM",
		"ST_UPSTREAM_UPDATE_ERR",
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

func (p *statsModule) prestart(agent *Agent) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *statsModule) start(agent *Agent) error {
	if agent.Conf.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case _, ok := <-p.running:
					if !ok {
						return
					}
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statsHandle())
				}
			}
		}()
	}
	return nil
}

func (p *statsModule) stop(agent *Agent) error {
	close(p.running)
	return nil
}
func (p *statsModule) reload(agent *Agent) error {
	return nil
}
