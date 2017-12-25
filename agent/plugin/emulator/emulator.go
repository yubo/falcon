/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package emulator

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/agent/utils"
)

func init() {
	agent.RegisterCollector(&emulator{
		gname: "emulator",
		name:  "main",
		tpl:   make(map[string]tpl),
	})
}

type tpl struct {
	file     string
	n        int
	interval int
	v        []float64
}

type emulator struct {
	name  string
	gname string
	tpl   map[string]tpl
	hn    int
	mn    int
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
	return p.v[idx] + float64(time%p.interval)*((p.v[(idx+1)%
		p.n]-p.v[idx])/float64(p.interval))
}

func (p *emulator) Name() string {
	return p.name
}

func (p *emulator) GName() string {
	return p.gname
}

func (p *emulator) Start(ag *agent.Agent) (err error) {

	dir := ag.Conf.Configer.Str(agent.C_EMU_TPL_DIR)

	fd, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer fd.Close()

	files, err := fd.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		n := len(file)
		if !(n > 4 && file[n-4:] == ".tpl") {
			continue
		}
		metric := file[:n-4]
		p.tpl[metric], err = readTpl(dir + "/" + file)
	}

	return nil
}

func (p *emulator) Collect() (ret []*falcon.Item, err error) {

	now := time.Now().Unix()
	ret = make([]*falcon.Item, len(p.tpl))

	n := 0
	for metric, _ := range p.tpl {
		tpl := p.tpl[metric]
		ret[n] = utils.GaugeValue(metric, tpl.emuValue(now))
		n++
	}

	return ret, nil
}

func (p *emulator) Reset() {
}
