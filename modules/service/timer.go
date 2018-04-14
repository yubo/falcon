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
	_timerTs int64
)

type timerModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *timerModule) prestart(s *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return nil
}

func (p *timerModule) start(s *Service) error {
	//start := time.Now().Unix()
	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				atomic.StoreInt64(&_timerTs, 0)
				return

			case <-ticker:
				now := time.Now().Unix()
				atomic.StoreInt64(&_timerTs, now)
				/*
					if s.Conf.Debug > 1 {
						atomic.StoreInt64(&_timerTs,
							start+(now-start)*DEBUG_MULTIPLES)
					} else {
						atomic.StoreInt64(&_timerTs, now)
					}
				*/
			}
		}
	}()
	return nil
}

func (p *timerModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *timerModule) reload(s *Service) error {
	return nil
}

func now() int64 {
	ts := atomic.LoadInt64(&_timerTs)
	if ts != 0 {
		return ts
	}
	return time.Now().Unix()
}
