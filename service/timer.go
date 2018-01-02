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

var (
	timer *TimerModule
)

func init() {
	timer = &TimerModule{}
}

type TimerModule struct {
	ctx    context.Context
	cancel context.CancelFunc
	ts     int64
}

func (p *TimerModule) now() int64 {
	ts := atomic.LoadInt64(&p.ts)
	if ts != 0 {
		return ts
	}
	return time.Now().Unix()
}

func (p *TimerModule) prestart(s *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	timer = p
	return nil
}

func (p *TimerModule) start(s *Service) error {
	//start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return

			case <-ticker:
				now := time.Now().Unix()
				atomic.StoreInt64(&p.ts, now)
				/*
					if s.Conf.Debug > 1 {
						atomic.StoreInt64(&p.ts,
							start+(now-start)*DEBUG_MULTIPLES)
					} else {
						atomic.StoreInt64(&p.ts, now)
					}
				*/
			}
		}
	}()
	return nil
}

func (p *TimerModule) stop(s *Service) error {
	p.cancel()
	atomic.StoreInt64(&p.ts, 0)
	return nil
}

func (p *TimerModule) reload(s *Service) error {
	return nil
}
