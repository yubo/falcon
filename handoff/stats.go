/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
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
	statName [ST_ARRAY_SIZE]string = [ST_ARRAY_SIZE]string{
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

func statStart(config HandoffOpts, p *specs.Process) {
	if config.Debug > 0 {
		ticker := time.NewTicker(time.Second * DEBUG_STAT_STEP).C
		go func() {
			for {
				select {
				case <-ticker:
					glog.V(3).Info(statHandle())
				}
			}
		}()
	}
}
