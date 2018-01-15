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
	lru       *queue
	cleanChan chan *eventEntry
	ctx       context.Context
	cancel    context.CancelFunc
}

func (p *TaskModule) prestart(alarm *Alarm) error {
	p.lru = &alarm.lru
	p.cleanChan = alarm.delEventEntryChan
	return nil
}

func (p *TaskModule) start(alarm *Alarm) (err error) {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go lruCleanWorker(p.ctx, p.lru, p.cleanChan)
	return nil
}

func (p *TaskModule) stop(alarm *Alarm) error {
	p.cancel()
	return nil
}

func (p *TaskModule) reload(alarm *Alarm) error {
	p.stop(alarm)
	time.Sleep(time.Second)
	p.prestart(alarm)
	return p.start(alarm)
}

func lruCleanWorker(ctx context.Context, list *queue, ch chan *eventEntry) {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			glog.V(3).Infof("%s lruCleanWorker entering", MODULE_NAME)
			now := time.Now().Unix()

			l := list.dequeue()
			if l == nil {
				continue
			}

			e := list_entry(l)
			if now-e.lastTs < EVENT_TIMEOUT {
				list.addHead(l)
				continue
			}

			ch <- e
		}
	}
}
