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

// TaskModule: alarm's module for sms/email gateway
// servicegroup: upstream container
// upstream: connection to the

type TaskModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *TaskModule) prestart(alarm *Alarm) error {
	return nil
}

func (p *TaskModule) start(alarm *Alarm) (err error) {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	expireTime, _ := alarm.Conf.Configer.Int64(C_EVENT_EXPIRE_TIME)

	go p.cleanWorker(alarm.delEventEntryChan, &alarm.lru, expireTime)
	return nil
}

func (p *TaskModule) stop(alarm *Alarm) error {
	p.cancel()
	return nil
}

func (p *TaskModule) reload(alarm *Alarm) error {
	p.stop(alarm)
	time.Sleep(time.Second)
	return p.start(alarm)
}

func (p *TaskModule) cleanWorker(ch chan *eventEntry, list *queue, expireTime int64) {

	ticker := time.NewTicker(time.Second * EVENT_CLEAN_INTERVAL).C
	ctx := p.ctx

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			glog.V(5).Infof("%s cleanWorker entering", MODULE_NAME)
			now := time.Now().Unix()

			for l := list.dequeue(); l != nil; l = list.dequeue() {
				e := list_entry(l)
				if now-e.lastTs < expireTime {
					list.addHead(l)
					break
				}
				ch <- e
			}
		}
	}
}
