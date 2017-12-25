/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
)

type TimerModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *TimerModule) prestart(b *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return nil
}

func (p *TimerModule) start(b *Service) error {
	start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return

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

func (p *TimerModule) stop(b *Service) error {
	p.cancel()
	atomic.StoreInt64(&b.ts, 0)
	return nil
}

func (p *TimerModule) reload(b *Service) error {
	return nil
}
