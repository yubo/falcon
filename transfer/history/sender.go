/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
)

// BackendModule: transfer's module for banckend
// backendgroup: upstream container
// upstream: connection to the
var (
	_sender map[string]sender
)

type sender interface {
	new(*Transfer) sender
	start(string) error
	stop() error
	addClientChan(string, string, chan *falcon.Item) error
}

func registerSender(name string, u sender) {
	_sender[name] = u
}

func init() {
	_sender = make(map[string]sender)
	//registerSender("tsdb", &senderTsdb{})
	registerSender("falcon", &senderFalcon{})
}

type upstream struct {
	name     string
	upstream sender
}

/* upstream */
type BackendModule struct {
	shareUpsteam []*upstream
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *BackendModule) prestart(transfer *Transfer) error {

	p.shareUpsteam = make([]*upstream, 0)
	for _, v := range transfer.Conf.ShareMap {
		b := &upstream{name: v.Name}

		b.scheduler = newSchedConsistent()

		for shareid, addr := range transfer.Conf.ShareMap {
			ch := make(chan *falcon.Item)
			b.scheduler.addChan(node, ch)
			b.upstream.addClientChan(node, addr, ch)
		}
		p.backends = append(p.backends, b)
	}

	return nil
}

func (p *BackendModule) start(transfer *Transfer) error {

	glog.V(3).Infof(MODULE_NAME+"%s upstreamStart len(bs) %d", transfer.Conf.Name, len(p.backends))

	p.ctx, p.cancel = context.WithCancel(context.Background())

	appUpdateChan := transfer.appUpdateChan

	for _, b := range p.backends {
		b.upstream.start(b.name)
	}

	go func() {
		for {
			select {
			case <-p.ctx.Done():
				for _, b := range p.backends {
					b.upstream.stop()
				}
				return
			case items := <-appUpdateChan:
				for _, b := range p.backends {
					for _, item := range items {
						b.scheduler.sched(item.Id()) <- item
					}
				}
			}
		}
	}()

	return nil
}

func (p *BackendModule) stop(transfer *Transfer) error {
	p.cancel()
	return nil
}

func (p *BackendModule) reload(transfer *Transfer) error {
	p.stop(transfer)
	time.Sleep(time.Second)
	p.prestart(transfer)
	return p.start(transfer)
}
