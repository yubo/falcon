/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/modules/ctrl/api/models"
	"github.com/yubo/falcon/modules/ctrl/api/routers"
	"github.com/yubo/gotool/ratelimits"
)

const (
	C_RL_GC_INTERVAL = "rlgcinterval"
	C_RL_GC_TIMEOUT  = "rlgctimeout"
	C_RL_LIMIT       = "rllimit"
	C_RL_ACCURACY    = "rlaccuracy"
)

var (
	_rl *ratelimits.RateLimits
)

func init() {
	routers.RateLimitsAccess = rateLimitsAccess
}

type rateLimitsModule struct {
	rl *ratelimits.RateLimits
}

func (p *rateLimitsModule) PreStart(c *Ctrl) error {
	return nil
}

func (p *rateLimitsModule) Start(c *Ctrl) error {
	conf := c.Conf.HttpRateLimit

	if !conf.Enable || conf.Limit <= 0 {
		return nil
	}

	r, err := ratelimits.New(uint32(conf.Limit), uint32(conf.Accuracy))
	if err != nil {
		return err
	}

	err = r.GcStart(time.Duration(conf.GcTimeout)*time.Millisecond,
		time.Duration(conf.GcInterval)*time.Millisecond)
	if err != nil {
		return err
	}
	_rl = r
	return nil
}

func (p *rateLimitsModule) Stop(c *Ctrl) error {
	if _rl != nil {
		_rl.GcStop()
		_rl = nil
	}
	return nil
}

func (p *rateLimitsModule) Reload(c *Ctrl) error {
	p.Stop(c)
	p.Start(c)
	return nil
}

func rateLimitsAccess(ctx *context.Context) bool {
	if _rl == nil {
		return true
	}

	ip := models.GetIPAdress(ctx.Request)
	if _rl.Update(ip) {
		return true
	}

	http.Error(ctx.ResponseWriter, "Too Many Requests", 429)
	return false
}
