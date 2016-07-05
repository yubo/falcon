/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"time"

	"github.com/yubo/falcon/specs"
)

type Collector interface {
	Collect(int, string) ([]*specs.MetaData, error)
}

var (
	collectConfig AgentOpts
	collectors    []Collector
)

func Collector_Register(c Collector) {
	collectors = append(collectors, c)
}

func collect(step int, host string) {
	t := time.NewTicker(time.Second * time.Duration(step)).C
	for {
		<-t
		vs := []*specs.MetaData{}
		for _, c := range collectors {
			if items, err := c.Collect(step, host); err == nil {
				vs = append(vs, items...)
			}
		}
		appUpdateChan <- &vs
	}
}

func collectStart(config AgentOpts) {
	collectConfig = config
	go collect(collectConfig.Interval, collectConfig.Host)
}
