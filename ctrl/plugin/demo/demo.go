/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package demo

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/config"
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

func (p *Demo) PreStart(c *config.Ctrl) error {
	glog.Info(MODULE_NAME + " prestart")
	return nil
}

func (p *Demo) Start(c *config.Ctrl) error {
	glog.Info(MODULE_NAME + " start")
	p.ctx, p.cancel = context.WithCancel(context.Background())

	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-ticker.C:
				glog.Info(MODULE_NAME + " .")
			}
		}
	}()
	return nil

}

func (p *Demo) Stop(c *config.Ctrl) error {
	glog.Info(MODULE_NAME + " stop")
	p.cancel()
	return nil
}

func (p *Demo) Reload(old, c *config.Ctrl) error {
	glog.Info(MODULE_NAME + " reload")
	p.Stop(c)
	time.Sleep(time.Second)
	p.PreStart(c)
	return p.Start(c)
}
