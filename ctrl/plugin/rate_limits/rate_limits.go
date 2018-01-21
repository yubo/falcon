/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package RateLimits

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/gotool/ratelimits"
)

const (
	C_RL_GC_INTERVAL = "rlgcinterval"
	C_RL_GC_TIMEOUT  = "rlgctimeout"
	C_RL_LIMIT       = "rllimit"
	C_RL_ACCURACY    = "rlaccuracy"
)

var (
	module = &RateLimitsModule{}
)

func init() {
	ctrl.RegisterModule(&RateLimitsModule{})
}

type RateLimitsModule struct {
	rl *ratelimits.RateLimits
}

func (p *RateLimitsModule) PreStart(c *ctrl.Ctrl) error {
	ctrl.Hooks.RateLimitsAccess = rateLimitsAccess
	module = p
	return nil
}

func (p *RateLimitsModule) Start(c *ctrl.Ctrl) error {
	conf := &c.Conf.Ctrl

	limit, _ := conf.Int(C_RL_LIMIT)
	accuracy, _ := conf.Int(C_RL_ACCURACY)

	if limit <= 0 {
		return nil
	}

	r, err := ratelimits.New(uint32(limit), uint32(accuracy))
	if err != nil {
		return err
	}

	timeout, _ := conf.Int(C_RL_GC_TIMEOUT)
	interval, _ := conf.Int(C_RL_GC_INTERVAL)
	err = r.GcStart(time.Duration(timeout)*time.Millisecond, time.Duration(interval)*time.Millisecond)
	if err != nil {
		return err
	}
	p.rl = r
	return nil
}

func (p *RateLimitsModule) Stop(c *ctrl.Ctrl) error {
	if p.rl != nil {
		p.rl.GcStop()
		p.rl = nil
	}
	return nil
}

func (p *RateLimitsModule) Reload(c *ctrl.Ctrl) error {
	p.Stop(c)
	p.Start(c)
	return nil
}

func rateLimitsAccess(ctx *context.Context) bool {
	if module.rl == nil {
		return true
	}

	ip := models.GetIPAdress(ctx.Request)
	if module.rl.Update(ip) {
		return true
	}

	http.Error(ctx.ResponseWriter, "Too Many Requests", 429)
	return false
}
