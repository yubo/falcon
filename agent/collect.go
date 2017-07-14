/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type collector_t struct {
	a []Collector
	m map[string]Collector
}

var (
	collector collector_t
)

type Collector interface {
	Name() string
	Start(*Agent) error
	Collect(int, string) ([]*falcon.MetaData, error)
	Reset()
}

func init() {
	collector.m = make(map[string]Collector)
}

func RegisterCollector(c Collector) {
	glog.V(4).Infof(MODULE_NAME+"register collector %s", c.Name())
	collector.m[c.Name()] = c
}

type CollectModule struct {
	running chan struct{}
}

func (p *CollectModule) prestart(agent *Agent) error {
	p.running = make(chan struct{}, 0)
	keys := make(map[string]bool)

	plugins := strings.Split(agent.Conf.Configer.Str(C_PLUGINS), ",")

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if c, ok := collector.m[plugin]; ok {
			// skip if exists
			if keys[plugin] {
				continue
			}
			collector.a = append(collector.a, c)
			keys[plugin] = true
		}
	}

	return nil
}

func (p *CollectModule) start(agent *Agent) error {

	host := agent.Conf.Host
	i, _ := agent.Conf.Configer.Int(C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(i)).C
	for _, c := range collector.a {
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
				vs := []*falcon.MetaData{}
				for _, c := range collector.a {
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
