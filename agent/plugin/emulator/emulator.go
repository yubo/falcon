/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package emulator

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/agent/utils"
)

func init() {
	agent.RegisterCollector(&emulator{})
}

type tpl struct {
	file     string
	n        int
	interval int
	v        []float64
}

type emulator struct {
	enable bool
	host   []string
	metric []string
	// tpl[host][metric] tpl
	tpl map[string]map[string]tpl
	hn  int
	mn  int
}

func readTpl(filePath string) (tpl, error) {
	tpl := tpl{file: filePath}
	fd, err := os.Open(filePath)
	if err != nil {
		glog.Fatalf("open %s: %v", tpl.file, err)
	}
	defer fd.Close()

	_, err = fmt.Fscanf(fd, "%d %d", &tpl.n, &tpl.interval)
	if err != nil {
		glog.Fatalf("open %s: %v", tpl.file, err)
	}

	tpl.v = make([]float64, tpl.n)
	for i := 0; i < tpl.n; i++ {
		_, err = fmt.Fscanf(fd, "%f", &tpl.v[i])
		if err != nil {
			glog.Fatalf("open %s: %v", tpl.file, err)
		}
	}
	return tpl, nil
}

func (p *tpl) emuValue(ts int64) float64 {
	time := int(ts % int64(p.n*p.interval))
	idx := time / p.interval
	return p.v[idx] + float64(time%p.interval)*((p.v[(idx+1)%p.n]-p.v[idx])/float64(p.interval))
}

func (p *emulator) Name() string {
	return "emulator"
}

func (p *emulator) Start(ag *agent.Agent) (err error) {
	var (
		host   string
		metric string
		tp     string
		tplnum int
	)
	p.enable, err = ag.Conf.Configer.Bool(agent.C_EMU_ENABLE)
	if !p.enable {
		return nil
	}

	host = ag.Conf.Configer.Str(agent.C_EMU_HOST)
	p.hn, _ = ag.Conf.Configer.Int(agent.C_EMU_HOSTNUM)
	metric = ag.Conf.Configer.Str(agent.C_EMU_METRIC)
	p.mn, _ = ag.Conf.Configer.Int(agent.C_EMU_METRICNUM)
	tp = ag.Conf.Configer.Str(agent.C_EMU_TPL)
	tplnum, _ = ag.Conf.Configer.Int(agent.C_EMU_TPLNUM)

	p.tpl = make(map[string]map[string]tpl, p.hn)

	for n, i := 0, 0; i < p.hn; i++ {
		_host := fmt.Sprintf(host, i)
		p.tpl[_host] = make(map[string]tpl)

		for j := 0; j < p.mn; j++ {
			_metric := fmt.Sprintf(metric, j)

			p.tpl[_host][_metric], err = readTpl(fmt.Sprintf(tp, n%tplnum))
			if err != nil {
				glog.Fatal(err)
			}
			n++
		}
	}

	return nil
}

func (p *emulator) Collect(step int,
	host string) (ret []*falcon.MetaData, err error) {
	if !p.enable {
		return nil, errors.New("not enable")
	}

	now := time.Now().Unix()
	ret = make([]*falcon.MetaData, p.hn*p.mn)

	n := 0
	for host, _ := range p.tpl {
		for metric, _ := range p.tpl[host] {
			tpl := p.tpl[host][metric]
			ret[n] = utils.GaugeValue(tpl.interval, host, metric,
				tpl.emuValue(now))
		}
	}

	return ret, nil
}

func (p *emulator) Reset() {
}
