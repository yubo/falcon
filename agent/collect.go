/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"time"

	"github.com/yubo/falcon"
)

type Collector interface {
	Collect(int, string) ([]*falcon.MetaData, error)
}

var (
	collectors []Collector
)

func Collector_Register(c Collector) {
	collectors = append(collectors, c)
}

func (p *Agent) collectStart() {
	ticker := time.NewTicker(time.Second *
		time.Duration(p.Conf.Interval)).C
	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case <-ticker:
				vs := []*falcon.MetaData{}
				for _, c := range collectors {
					if items, err := c.Collect(p.Conf.Interval,
						p.Conf.Params.Host); err == nil {
						vs = append(vs, items...)
					}
				}
				p.appUpdateChan <- &vs
			}
		}
	}()
}

func (p *Agent) collectStop() {
}
