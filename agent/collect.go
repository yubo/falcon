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
	Start(context.Context, *Agent) error
	Reset()
	Collect() ([]*Item, error)
}

func init() {
	collectorGroups = make(map[string]map[string]Collector)
}

func RegisterCollector(c Collector) {
	if _, ok := collectorGroups[c.GName()]; !ok {
		collectorGroups[c.GName()] = make(map[string]Collector)
	}
	collectorGroups[c.GName()][c.Name()] = c
}

type CollectModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *CollectModule) prestart(agent *Agent) error {
	return nil
}

func (p *CollectModule) start(agent *Agent) error {
	interval, _ := agent.Conf.Configer.Int(C_INTERVAL)
	collectors := getCollectors(strings.Split(agent.Conf.Configer.Str(C_PLUGINS), ","))
	putChan := agent.putChan

	p.ctx, p.cancel = context.WithCancel(context.Background())

	for _, c := range collectors {
		glog.V(4).Infof("%s plugins %s.Start()", MODULE_NAME, falcon.GetType(c))
		if err := c.Start(p.ctx, agent); err != nil {
			glog.V(4).Infof("%s plugins %s.Start() err %v", MODULE_NAME, falcon.GetType(c), err)
			return err
		}
	}

	go p.collectWorker(putChan, collectors, interval)

	return nil
}

func (p *CollectModule) stop(agent *Agent) error {
	p.cancel()
	return nil
}

func (p *CollectModule) reload(agent *Agent) error {
	p.stop(agent)
	time.Sleep(time.Second)
	return p.start(agent)
}

func (p *CollectModule) collectWorker(putChan chan *putContext, collectors []Collector, interval int) {
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			for _, c := range collectors {
				if items, err := c.Collect(); err == nil {
					putChan <- &putContext{items: items}
				}
			}
		}
	}

}

func getCollectors(plugins []string) []Collector {
	collectors := []Collector{}
	keys := make(map[string]bool)

	for _, plugin := range plugins {
		plugin = strings.TrimSpace(plugin)
		if group, ok := collectorGroups[plugin]; ok {
			// skip if exists
			if keys[plugin] {
				continue
			}
			for _, c := range group {
				glog.V(4).Infof("%s add plugin %s", MODULE_NAME, plugin)
				collectors = append(collectors, c)
			}
			keys[plugin] = true
		} else {
			glog.Infof("plugin %s miss", plugin)
		}
	}
	return collectors
}
