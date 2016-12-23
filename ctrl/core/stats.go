/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

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
		ret += fmt.Sprintf("%s %d\n", statName[i],
			atomic.LoadUint64(&statCnt[i]))
	}
	return ret
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

func (p *Ctrl) statStart() {
	if p.Params.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-ticker:
					glog.V(3).Info(MODULE_NAME + statHandle())
				case _, ok := <-p.running:
					if !ok {
						return
					}
				}
			}
		}()
	}
}

func (p *Ctrl) statStop() {
}
