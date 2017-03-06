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

var (
	_collector []Collector
)

type Collector interface {
	Collect(int, string) ([]*falcon.MetaData, error)
	Reset()
}

func RegisterCollector(c Collector) {
	_collector = append(_collector, c)
}

type collectModule struct {
	running chan struct{}
}

func (p *collectModule) prestart(agent *Agent) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *collectModule) start(agent *Agent) error {

	host := agent.Conf.Host
	i, _ := agent.Conf.Configer.Int(falcon.C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(i)).C

	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case <-ticker:
				vs := []*falcon.MetaData{}
				for _, c := range _collector {
					if items, err := c.Collect(i,
						host); err == nil {
						vs = append(vs, items...)
					}
				}
				agent.appUpdateChan <- &vs
			}
		}
	}()
	return nil
}

func (p *collectModule) stop(agent *Agent) error {
	close(p.running)
	return nil
}

func (p *collectModule) reload(agent *Agent) error {
	return nil
}
