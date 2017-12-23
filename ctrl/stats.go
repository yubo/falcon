/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/config"
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
	statName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
		"ST_UPSTREAM_RECONNECT",
		"ST_UPSTREAM_DIAL",
		"ST_UPSTREAM_DIAL_ERR",
		"ST_UPSTREAM_UPDATE",
		"ST_UPSTREAM_UPDATE_ITEM",
		"ST_UPSTREAM_UPDATE_ERR",
	}
)

var (
	statCnt [ST_ARRAY_SIZE]uint64
)

func statHandle() (ret string) {
	for i := 0; i < ST_ARRAY_SIZE; i++ {
		ret += fmt.Sprintf("%d ", atomic.LoadUint64(&statCnt[i]))
	}
	return ret[:len(ret)-1]
}

func statInc(idx, n int) {
	atomic.AddUint64(&statCnt[idx], uint64(n))
}

func statSet(idx, n int) {
	atomic.StoreUint64(&statCnt[idx], uint64(n))
}

func statGet(idx int) uint64 {
	return atomic.LoadUint64(&statCnt[idx])
}

type StatsModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *StatsModule) PreStart(c *config.Ctrl) error {
	return nil
}

func (p *StatsModule) Start(c *config.Ctrl) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	if c.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statHandle())
				}
			}
		}()
	}
	return nil
}

func (p *StatsModule) Stop(c *config.Ctrl) error {
	p.cancel()
	return nil
}

func (p *StatsModule) Reload(old, c *config.Ctrl) error {
	p.Stop(c)
	time.Sleep(time.Second)
	p.PreStart(c)
	return p.Start(c)
}
