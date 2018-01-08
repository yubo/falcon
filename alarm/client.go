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

// ClientModule: transfer's module for banckend
// servicegroup: upstream container
// upstream: connection to the

type ClientModule struct {
	putChan     chan *Event
	callTimeout int
	burstSize   int
	ctx         context.Context
	cancel      context.CancelFunc
}

func (p *ClientModule) prestart(transfer *Alarm) error {

	p.putChan = transfer.appPutChan
	p.callTimeout, _ = transfer.Conf.Configer.Int(C_CALL_TIMEOUT)

	return nil
}

func (p *ClientModule) start(transfer *Alarm) (err error) {

	glog.V(3).Infof(MODULE_NAME+"%s", transfer.Conf.Name)

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go putWorker(p.ctx, p.putChan)

	return nil
}

func (p *ClientModule) stop(transfer *Alarm) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(transfer *Alarm) error {
	p.stop(transfer)
	time.Sleep(time.Second)
	p.prestart(transfer)
	return p.start(transfer)
}

func putWorker(ctx context.Context, ch chan *Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			glog.V(3).Infof("%v", event)
		}
	}

}
