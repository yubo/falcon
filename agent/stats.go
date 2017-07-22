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

	"github.com/golang/glog"
	"golang.org/x/net/context"
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

type StatsModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *StatsModule) prestart(agent *Agent) error {
	return nil
}

func (p *StatsModule) start(agent *Agent) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	if agent.Conf.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statsHandle())
				}
			}
		}()
	}
	return nil
}

func (p *StatsModule) stop(agent *Agent) error {
	p.cancel()
	return nil
}
func (p *StatsModule) reload(agent *Agent) error {
	p.stop(agent)
	time.Sleep(time.Second)
	p.prestart(agent)
	return p.start(agent)
}
