/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"sync/atomic"
	"time"
)

type timerModule struct {
	running chan struct{}
}

func (p *timerModule) prestart(b *Backend) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *timerModule) start(b *Backend) error {
	start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}

			case <-ticker:
				now := time.Now().Unix()
				if b.Conf.Debug > 1 {
					atomic.StoreInt64(&b.ts,
						start+(now-start)*DEBUG_MULTIPLES)
				} else {
					atomic.StoreInt64(&b.ts, now)
				}
			}
		}
	}()
	return nil
}

func (p *timerModule) stop(b *Backend) error {
	close(p.running)
	atomic.StoreInt64(&b.ts, 0)
	return nil
}

func (p *timerModule) reload(b *Backend) error {
	return nil
}
