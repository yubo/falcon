/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package sys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/agent/utils"
)

const (
	_ = iota
	PROC_STAT_USER
	PROC_STAT_NICE
	PROC_STAT_SYSTEM
	PROC_STAT_IDLE
	PROC_STAT_IOWAIT
	PROC_STAT_IRQ
	PROC_STAT_SOFTIRQ
	PROC_STAT_STEAL
	PROC_STAT_GUEST
	PROC_STAT_TOTAL
	PROC_STAT_SIZE
)

var (
	procStatName []string = []string{
		"",
		"cpu.user",
		"cpu.nice",
		"cpu.system",
		"cpu.idle",
		"cpu.iowait",
		"cpu.irq",
		"cpu.softirq",
		"cpu.steal",
		"cpu.guest",
		"cpu.total",
		"PROC_STAT_SIZE",
	}
)

type cpuCoreSample struct {
	data [PROC_STAT_SIZE]uint64
}

func (p *cpuCoreSample) String() string {
	return fmt.Sprintf("User %d Nice %d System %d Idle %d Iowait %d"+
		" Irq %d SoftIrq %d Steal %d Guest %d Total %d",
		p.data[PROC_STAT_USER],
		p.data[PROC_STAT_NICE],
		p.data[PROC_STAT_SYSTEM],
		p.data[PROC_STAT_IDLE],
		p.data[PROC_STAT_IOWAIT],
		p.data[PROC_STAT_IRQ],
		p.data[PROC_STAT_SOFTIRQ],
		p.data[PROC_STAT_STEAL],
		p.data[PROC_STAT_GUEST],
		p.data[PROC_STAT_TOTAL])
}

type cpuStatSample struct {
	cpu          *cpuCoreSample
	cpus         []*cpuCoreSample
	ctxt         uint64
	processes    uint64
	procsRunning uint64
	procsBlocked uint64
}

func (p *cpuStatSample) String() string {
	return fmt.Sprintf("Cpu %v Cpus %v Ctxt %d"+
		" Processes %d ProcsRunning %d ProcsBlocking %d",
		p.cpu,
		p.cpus,
		p.ctxt,
		p.processes,
		p.procsRunning,
		p.procsBlocked)
}

func init() {
	agent.RegisterCollector(&cpuCollector{
		name:  "cpu",
		gname: "sys",
	})
}

type cpuCollector struct {
	cur   *cpuStatSample
	last  *cpuStatSample
	name  string
	gname string
}

func (p *cpuCollector) GName() string {
	return p.gname
}

func (p *cpuCollector) Name() string {
	return p.name
}

func (p *cpuCollector) Reset() {
	p.cur = nil
	p.last = nil
}

func (p *cpuCollector) Start(agent *agent.Agent) error {
	glog.V(5).Infof("start")
	return nil
}

func (p *cpuCollector) Collect() (ret []*falcon.Item, err error) {
	p.last = p.cur
	p.cur, err = p.collect()
	if err != nil {
		return nil, err
	}
	glog.V(5).Infof("%v", p.cur)
	return p.stat()
}

func (p *cpuCollector) collect() (*cpuStatSample, error) {

	bs, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	ps := &cpuStatSample{cpus: make([]*cpuCoreSample, runtime.NumCPU())}
	reader := bufio.NewReader(bytes.NewBuffer(bs))

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return ps, err
		}
		_cpuParseLine(line, ps)
	}
	return ps, nil
}

func (p *cpuCollector) stat() (ret []*falcon.Item, err error) {
	var n float64
	if p.last == nil {
		return nil, errors.New("no data")
	}

	if total := p.cur.cpu.data[PROC_STAT_TOTAL] -
		p.last.cpu.data[PROC_STAT_TOTAL]; total == 0 {
		n = 0.0
	} else {
		n = 100.0 / float64(total)
	}

	for i := 1; i < PROC_STAT_TOTAL; i++ {
		ret = append(ret, utils.GaugeValue(procStatName[i],
			float64(p.cur.cpu.data[i]-p.last.cpu.data[i])*n))
	}
	ret = append(ret, utils.GaugeValue("cpu.busy",
		float64(100.0-float64(p.cur.cpu.data[PROC_STAT_IDLE]-
			p.last.cpu.data[PROC_STAT_IDLE])*n)))

	ret = append(ret, utils.GaugeValue("cpu.switches",
		float64(p.cur.ctxt-p.last.ctxt)))

	return ret, nil
}

func _cpuParseLine(line string, ps *cpuStatSample) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return
	}

	fieldName := fields[0]
	if fieldName == "cpu" {
		ps.cpu = cpuParseFields(fields)
		return
	}

	if strings.HasPrefix(fieldName, "cpu") {
		idx, err := strconv.Atoi(fieldName[3:])
		if err != nil || idx >= len(ps.cpus) {
			return
		}

		ps.cpus[idx] = cpuParseFields(fields)
		return
	}

	if fieldName == "ctxt" {
		ps.ctxt, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "processes" {
		ps.processes, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_running" {
		ps.procsRunning, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}

	if fieldName == "procs_blocked" {
		ps.procsBlocked, _ = strconv.ParseUint(fields[1], 10, 64)
		return
	}
}

func cpuParseFields(fields []string) *cpuCoreSample {
	cu := cpuCoreSample{}
	for i := 1; i < PROC_STAT_TOTAL; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}
		cu.data[PROC_STAT_TOTAL] += val
		cu.data[i] = val
	}
	return &cu
}
