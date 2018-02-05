/*
 * Copyright 2016,2017,2018 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package demo

import (
	"context"
	"math"
	"time"

	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/transfer"
)

const (
	x = (2 * math.Pi) / float64(3600)
)

func init() {
	agent.RegisterCollector(&cosCollector{})
}

type cosCollector struct{}

func (p *cosCollector) Name() (name, gname string) {
	return "cos", "demo"
}

func (p *cosCollector) Start(ctx context.Context, a *agent.Agent) error {
	interval, _ := a.Conf.Configer.Int(agent.C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C
	ch := a.PutChan

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				ch <- &agent.PutRequest{
					Dps: []*transfer.DataPoint{
						agent.GaugeValue("cos.ticker.demo",
							math.Cos(x*(float64(time.Now().Unix())+1800))),
					},
				}
			}
		}
	}()
	return nil
}

func (p *cosCollector) Collect() (ret []*transfer.DataPoint, err error) {
	return []*transfer.DataPoint{
		agent.GaugeValue("cos.collect.demo",
			math.Cos(x*float64(time.Now().Unix()))),
	}, nil

}
