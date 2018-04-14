/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package sys

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
)

type SysConfig struct {
	Dev      bool   `json:"dev"`
	LogFile  string `json:"log_file"`
	LogLevel int    `json:"log_level"`
	Cpu      bool   `json:"cpu"`
	CpuFile  string `json:"cpu_file"`
	Heap     bool   `json:"heap"`
	HeapFile string `json:"heap_file"`
	Http     bool   `json:"http"`
	HttpPort string `json:"http_port"`
	PidFile  string `json:"pid_file"`
}

type Sys struct {
	config *SysConfig

	// runtime
	server   *http.Server
	cpuProf  *os.File
	heapProf *os.File
	pid      int
}

const (
	MODULE_NAME = "sys"
)

func init() {
	core.RegisterModule(&Sys{})
}

/* called by app.Start() befor app.modules.prestart() */
func (p *Sys) ReadConfig(conf *core.Configer) error {
	p.config = &SysConfig{}

	err := conf.Read(MODULE_NAME, p.config)
	if err != nil {
		return err
	}

	p.pid = os.Getpid()
	setGlog(p.config)

	return err
}

func (p *Sys) Prestart() (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := p.config

	if c.Cpu {
		if p.cpuProf, err = os.OpenFile(c.CpuFile, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return
		}
		pprof.StartCPUProfile(p.cpuProf)
	}
	if c.Heap {
		if p.heapProf, err = os.OpenFile(c.HeapFile, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return
		}
		pprof.WriteHeapProfile(p.heapProf)
	}

	return nil
}

func (p *Sys) Start() (err error) {
	c := p.config

	if err = procCheck(c); err != nil {
		return err
	}

	if err = procSave(c.PidFile, p.pid); err != nil {
		return err
	}

	if c.Http {
		p.server = &http.Server{Addr: c.HttpPort}
		go p.server.ListenAndServe()
	}
	return nil
}

func (p *Sys) Stop() error {
	c := p.config

	pidfile := fmt.Sprintf("%s.%d", c.PidFile, p.pid)
	glog.Infof("process %d exiting", p.pid)
	os.Rename(c.PidFile, pidfile)

	if p.server != nil {
		p.server.Shutdown(nil)
		p.server = nil
	}

	if c.Cpu {
		pprof.StopCPUProfile()
		p.cpuProf.Close()
	}
	if c.Heap {
		p.heapProf.Close()
	}

	glog.Infof("%s pid:%d exit", MODULE_NAME, p.pid)
	os.Remove(pidfile)

	return nil
}

func (p *Sys) Reload(conf *core.Configer) error {
	// TODO
	/*
		newConfig := &SysConfig{}
		err := conf.Read(MODULE_NAME, &newConfig)
		if err != nil {
			return nil, err
		}

		p.Stop()
		time.Sleep(time.Second)

		p.config = newConfig
		p.PreStart()
		return p.Start()
	*/
	return nil
}

func (p *Sys) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}

func (p *Sys) Stats(conf *core.Configer) (string, error) {
	return "", nil
}

func setGlog(c *SysConfig) {
	glog.V(3).Infof("set glog %s, %d", c.LogFile, c.LogLevel)
	flag.Lookup("v").Value.Set(fmt.Sprintf("%d", c.LogLevel))

	if strings.ToLower(c.LogFile) == "stdout" {
		flag.Lookup("logtostderr").Value.Set("true")
		return
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
	}

	if fi, err := os.Stat(c.LogFile); err != nil || !fi.IsDir() {
		glog.Errorf("log dir %s does not exist or not dir", c.LogFile)
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
		flag.Lookup("log_dir").Value.Set(c.LogFile)
	}
}

func procCheck(c *SysConfig) error {
	pid, err := core.ReadFileInt(c.PidFile)
	if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%s", pid))
	if os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("proccess %s exist", pid)
}

func procSave(pidfile string, pid int) error {
	return ioutil.WriteFile(pidfile,
		[]byte(fmt.Sprintf("%d", pid)), 0644)
}
