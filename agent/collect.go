/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
)

var (
	collectorGroups map[string]map[string]Collector
	hostName        string
)

type Collector interface {
	GName() string
	Name() string
	Start(context.Context, *Agent) error
	Reset()
	Collect() ([]*transfer.DataPoint, error)
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
	hostName = agent.Conf.Host

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
				if dps, err := c.Collect(); err == nil {
					putChan <- &putContext{dps: dps}
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

func NewMetricValue(metric string,
	val float64, typ string, tags ...string) *transfer.DataPoint {
	var tags_ string

	if len(tags) > 0 {
		sort.Strings(tags)
		tags_ = strings.Join(tags, ",")
	}

	return &transfer.DataPoint{
		Key: falcon.AttrKey(hostName, metric, tags_, typ),
		Value: &tsdb.TimeValuePair{
			Timestamp: time.Now().Unix(),
			Value:     val,
		},
	}
}

func GaugeValue(metric string, val float64, tags ...string) *transfer.DataPoint {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val float64, tags ...string) *transfer.DataPoint {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}
