/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"time"

	"github.com/yubo/falcon/utils"
)

var (
	_collector []Collector
)

type Collector interface {
	Start(*Agent) error
	Collect(int, string) ([]*utils.MetaData, error)
	Reset()
}

func RegisterCollector(c Collector) {
	_collector = append(_collector, c)
}

type CollectModule struct {
	running chan struct{}
}

func (p *CollectModule) prestart(agent *Agent) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *CollectModule) start(agent *Agent) error {

	host := agent.Conf.Host
	i, _ := agent.Conf.Configer.Int(utils.C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(i)).C

	for _, c := range _collector {
		if err := c.Start(agent); err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case <-ticker:
				vs := []*utils.MetaData{}
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

func (p *CollectModule) stop(agent *Agent) error {
	close(p.running)
	return nil
}

func (p *CollectModule) reload(agent *Agent) error {
	return nil
}
