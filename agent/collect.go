/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
)

type collector_t struct {
	a []Collector
	//m map[string]Collector
	groups map[string]map[string]Collector
}

var (
	collector collector_t
)

type Collector interface {
	GName() string
	Name() string
	Start(*Agent) error
	Reset()
	Collect(int, string) ([]*falcon.Item, error)
}

func init() {
	//collector.m = make(map[string]Collector)
	collector.groups = make(map[string]map[string]Collector)
}

func RegisterCollector(c Collector) {
	glog.V(4).Infof(MODULE_NAME+"register collector %s", c.Name())
	//collector.m[c.Name()] = c
	if _, ok := collector.groups[c.GName()]; !ok {
		collector.groups[c.GName()] = make(map[string]Collector)
	}
	collector.groups[c.GName()][c.Name()] = c
}

type CollectModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *CollectModule) prestart(agent *Agent) error {
	keys := make(map[string]bool)

	plugins := strings.Split(agent.Conf.Configer.Str(C_PLUGINS), ",")

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if group, ok := collector.groups[plugin]; ok {
			// skip if exists
			if keys[plugin] {
				continue
			}
			for _, c := range group {
				collector.a = append(collector.a, c)
			}
			keys[plugin] = true
		}
	}

	return nil
}

func (p *CollectModule) start(agent *Agent) error {

	p.ctx, p.cancel = context.WithCancel(context.Background())
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
			case <-p.ctx.Done():
				return
			case <-ticker:
				for _, c := range collector.a {
					if items, err := c.Collect(i,
						host); err == nil {
						agent.appUpdateChan <- items
					}
				}
			}
		}
	}()
	return nil
}

func (p *CollectModule) stop(agent *Agent) error {
	p.cancel()
	return nil
}

func (p *CollectModule) reload(agent *Agent) error {
	p.stop(agent)
	time.Sleep(time.Second)
	p.prestart(agent)
	return p.start(agent)
}
