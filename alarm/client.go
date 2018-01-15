/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

// ClientModule: alarm's module for sms/email gateway
// servicegroup: upstream container
// upstream: connection to the

type ClientModule struct {
	actionChan      chan *Action
	workerProcesses int
	callTimeout     int
	burstSize       int
	ctx             context.Context
	cancel          context.CancelFunc
}

func (p *ClientModule) prestart(alarm *Alarm) error {
	return nil
}

func (p *ClientModule) start(alarm *Alarm) (err error) {

	p.callTimeout, _ = alarm.Conf.Configer.Int(C_CALL_TIMEOUT)
	p.workerProcesses, _ = alarm.Conf.Configer.Int(C_WORKER_PROCESSES)
	p.actionChan = alarm.actionChan

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go actionWorker(p.ctx, p.actionChan, p.workerProcesses)

	return nil
}

func (p *ClientModule) stop(alarm *Alarm) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(alarm *Alarm) error {
	p.stop(alarm)
	time.Sleep(time.Second)
	p.prestart(alarm)
	return p.start(alarm)
}

func actionWorker(ctx context.Context, ch chan *Action, n int) {
	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case action := <-ch:
					glog.V(3).Infof("%s >>action<< %v", MODULE_NAME, action)
				}
			}
		}()
	}
}
