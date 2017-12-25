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

var (
	collectorGroups map[string]map[string]Collector
)

type Collector interface {
	GName() string
	Name() string
	Start(*Agent) error
	Reset()
	Collect() ([]*falcon.Item, error)
}

func init() {
	collectorGroups = make(map[string]map[string]Collector)
}

func RegisterCollector(c Collector) {
	glog.V(4).Infof(MODULE_NAME+"register collector %s", c.Name())
	if _, ok := collectorGroups[c.GName()]; !ok {
		collectorGroups[c.GName()] = make(map[string]Collector)
	}
	collectorGroups[c.GName()][c.Name()] = c
}

type CollectModule struct {
	a []Collector

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *CollectModule) prestart(agent *Agent) error {

	p.a = []Collector{}

	keys := make(map[string]bool)

	plugins := strings.Split(agent.Conf.Configer.Str(C_PLUGINS), ",")

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if group, ok := collectorGroups[plugin]; ok {
			// skip if exists
			if keys[plugin] {
				continue
			}
			for _, c := range group {
				p.a = append(p.a, c)
			}
			keys[plugin] = true
		} else {
			glog.Infof("plugin %s miss", plugin)
		}
	}

	return nil
}

func (p *CollectModule) start(agent *Agent) error {

	p.ctx, p.cancel = context.WithCancel(context.Background())
	interval, _ := agent.Conf.Configer.Int(C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C

	for _, c := range p.a {
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
				for _, c := range p.a {
					if items, err := c.Collect(); err == nil {
						agent.appPutChan <- items
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
