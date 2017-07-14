/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/config"
)

type Module interface {
	New(config interface{}) Module
	Prestart() error
	Start() error
	Stop() error
	Reload(config interface{}) error
	Signal(os.Signal) error
	String() string
	Name() string
	Parse(text []byte, filename string, lino int, debug bool) config.ModuleConf
}

// reload not support add/del/disable module
type Process struct {
	Config *config.FalconConfig
	Pid    int
	Status uint32
	Module []Module
}

const (
	APP_STATUS_INIT = iota
	APP_STATUS_PENDING
	APP_STATUS_RUNNING
	APP_STATUS_EXIT
	APP_STATUS_RELOAD
)

func init() {
	ModuleTpls = make(map[string]Module)
}

func NewProcess(c *config.FalconConfig) *Process {
	p := &Process{
		Config: c,
		Pid:    os.Getpid(),
		Status: APP_STATUS_PENDING,
		Module: make([]Module, len(c.Conf)),
	}
	return p
}

func (p *Process) Kill(sig syscall.Signal) error {
	if pid, err := ReadFileInt(p.Config.PidFile); err != nil {
		return err
	} else {
		glog.Infof("kill %d %s\n", pid, sig)
		return syscall.Kill(pid, sig)
	}
}

func (p *Process) Check() error {
	pid, err := ReadFileInt(p.Config.PidFile)
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

func (p *Process) Save() error {
	return ioutil.WriteFile(p.Config.PidFile,
		[]byte(fmt.Sprintf("%d", p.Pid)), 0644)
}

func (p *Process) Start() {

	SetGlog(p.Config)

	for i := 0; i < len(p.Config.Conf); i++ {
		m, ok := ModuleTpls[GetType(p.Config.Conf[i])]
		if !ok {
			glog.Exitf("%s's module not support", GetType(p.Config.Conf[i]))
		}
		p.Module[i] = m.New(p.Config.Conf[i])
	}

	for i := 0; i < len(p.Module); i++ {
		p.Module[i].Prestart()
	}

	for i := 0; i < len(p.Module); i++ {
		p.Module[i].Start()
	}
}

func RegisterModule(m Module, names ...string) error {
	for _, name := range names {
		if _, ok := ModuleTpls[name]; ok {
			return ErrExist
		} else {
			ModuleTpls[name] = m
		}
	}
	return nil
}

func SetGlog(c *config.FalconConfig) {
	glog.V(3).Infof(MODULE_NAME+"set glog %s, %d", c.Log, c.Logv)
	flag.Lookup("v").Value.Set(fmt.Sprintf("%d", c.Logv))

	if strings.ToLower(c.Log) == "stdout" {
		flag.Lookup("logtostderr").Value.Set("true")
		return
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
	}

	if fi, err := os.Stat(c.Log); err != nil || !fi.IsDir() {
		glog.Errorf("log dir %s does not exist or not dir", c.Log)
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
		flag.Lookup("log_dir").Value.Set(c.Log)
	}
}
