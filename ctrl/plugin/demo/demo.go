/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package demo

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl"
	"golang.org/x/net/context"
)

const (
	MODULE_NAME = "\x1B[33m[CTRL_PLUGIN_DEMO]\x1B[0m"
)

func init() {
	ctrl.RegisterModule(&Demo{})
}

type Demo struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *Demo) PreStart(c *ctrl.Ctrl) error {
	return nil
}

func (p *Demo) Start(c *ctrl.Ctrl) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())

	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-ticker.C:
				glog.V(5).Info(MODULE_NAME + " .")
			}
		}
	}()
	return nil

}

func (p *Demo) Stop(c *ctrl.Ctrl) error {
	p.cancel()
	return nil
}

func (p *Demo) Reload(c *ctrl.Ctrl) error {
	return nil
}
