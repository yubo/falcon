/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package plugin

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

	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/specs"
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

func (this *cpuCoreSample) String() string {
	return fmt.Sprintf("<User:%d, Nice:%d, System:%d, Idle:%d, Iowait:%d,"+
		" Irq:%d, SoftIrq:%d, Steal:%d, Guest:%d, Total:%d>",
		this.data[PROC_STAT_USER],
		this.data[PROC_STAT_NICE],
		this.data[PROC_STAT_SYSTEM],
		this.data[PROC_STAT_IDLE],
		this.data[PROC_STAT_IOWAIT],
		this.data[PROC_STAT_IRQ],
		this.data[PROC_STAT_SOFTIRQ],
		this.data[PROC_STAT_STEAL],
		this.data[PROC_STAT_GUEST],
		this.data[PROC_STAT_TOTAL])
}

type cpuStatSample struct {
	cpu          *cpuCoreSample
	cpus         []*cpuCoreSample
	ctxt         uint64
	processes    uint64
	procsRunning uint64
	procsBlocked uint64
}

func (this *cpuStatSample) String() string {
	return fmt.Sprintf("<Cpu:%v, Cpus:%v, Ctxt:%d,"+
		" Processes:%d, ProcsRunning:%d, ProcsBlocking:%d>",
		this.cpu,
		this.cpus,
		this.ctxt,
		this.processes,
		this.procsRunning,
		this.procsBlocked)
}

var (
	_cpuCollector agent.Collector
)

func init() {
	_cpuCollector = &cpuCollector{}
	agent.Collector_Register(_cpuCollector)
}

type cpuCollector struct {
	cur  *cpuStatSample
	last *cpuStatSample
}

func (p *cpuCollector) Collect(step int,
	host string) (ret []*specs.MetaData, err error) {
	p.last = p.cur
	p.cur, err = p.collect()
	if err != nil {
		return nil, err
	}
	return p.stat(step, host)
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

func (p *cpuCollector) stat(step int,
	host string) (ret []*specs.MetaData, err error) {
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
		ret = append(ret, GaugeValue(step, host, procStatName[i],
			float64(p.cur.cpu.data[i]-p.last.cpu.data[i])*n))
	}
	ret = append(ret, GaugeValue(step, host, "cpu.busy",
		float64(100.0-(p.cur.cpu.data[PROC_STAT_IDLE]-
			p.last.cpu.data[PROC_STAT_IDLE]))))

	ret = append(ret, GaugeValue(step, host, "cpu.switches",
		float64(p.cur.ctxt)))

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
	sz := len(fields)
	for i := 1; i < sz; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}
		cu.data[PROC_STAT_TOTAL] += val
		cu.data[i] = val
	}
	return &cu
}
