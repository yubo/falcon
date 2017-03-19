/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

// backendModule: transfer's module for banckend
// backendgroup: upstream container
// upstream: connection to the
var (
	_sender map[string]sender
)

type sender interface {
	new(*Transfer) sender
	start(string) error
	stop() error
	addClientChan(string, string, chan *falcon.MetaData) error
}

func registerSender(name string, u sender) {
	_sender[name] = u
}

func init() {
	_sender = make(map[string]sender)
	registerSender("tsdb", &senderTsdb{})
	registerSender("falcon", &senderFalcon{})
}

type backend struct {
	name      string
	upstream  sender
	scheduler scheduler
}

/* upstream */
type backendModule struct {
	running  chan struct{}
	backends []*backend
}

func (p *backendModule) prestart(L *Transfer) error {
	p.running = make(chan struct{}, 0)

	p.backends = make([]*backend, 0)
	for _, v := range L.Conf.Backend {
		if v.Disabled {
			continue
		}
		b := &backend{name: v.Name}
		if st, ok := _sender[v.Type]; ok {
			b.upstream = st.new(L)
		} else {
			return falcon.ErrUnsupported
		}

		b.scheduler = newSchedConsistent()

		for node, addr := range v.Upstream {
			ch := make(chan *falcon.MetaData)
			b.scheduler.addChan(node, ch)
			b.upstream.addClientChan(node, addr, ch)
		}
		p.backends = append(p.backends, b)
	}

	return nil
}

func (p *backendModule) start(L *Transfer) error {

	glog.V(3).Infof(MODULE_NAME+"%s upstreamStart len(bs) %d", L.Conf.Name, len(p.backends))
	appUpdateChan := L.appUpdateChan

	for _, b := range p.backends {
		b.upstream.start(b.name)
	}

	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case items := <-appUpdateChan:
				for _, b := range p.backends {
					for _, item := range *items {
						ch := b.scheduler.sched(item.Id())
						ch <- item
					}
				}
			}
		}
	}()

	return nil
}

func (p *backendModule) stop(L *Transfer) error {

	close(p.running)
	for _, b := range p.backends {
		b.upstream.stop()
	}
	return nil
}

func (p *backendModule) reload(L *Transfer) error {
	return nil
}
