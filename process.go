/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
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
	Parse(text []byte, filename string, lino int) ModuleConf
	Stats(config interface{}) (string, error)
}

// reload not support add/del/disable module
type Process struct {
	Config *FalconConfig
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
	Modules = make(map[string]Module)
}

func NewProcess(c *FalconConfig) *Process {
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
	rand.Seed(time.Now().Unix())

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

func (p *Process) Stats(module string) error {
	for i := 0; i < len(p.Config.Conf); i++ {
		m, ok := ModuleTpls[GetType(p.Config.Conf[i])]
		if !ok {
			return fmt.Errorf("%s's module not support", GetType(p.Config.Conf[i]))
		}
		p.Module[i] = m.New(p.Config.Conf[i])
	}

	for i := 0; i < len(p.Module); i++ {
		m := p.Module[i]
		if module == "all" || module == "" || module == m.Name() {
			stats, err := m.Stats(p.Config.Conf[i])
			if err != nil {
				return err
			}
			fmt.Printf("%s\n%s\n", m.Name(), stats)
		}
	}
	return nil
}

func RegisterModule(m Module, name, tpl string) error {
	if _, ok := Modules[name]; ok {
		return ErrExist
	}

	if _, ok := ModuleTpls[tpl]; ok {
		return ErrExist
	}

	Modules[name] = m
	ModuleTpls[tpl] = m
	return nil
}

func SetGlog(c *FalconConfig) {
	glog.V(3).Infof("%s set glog %s, %d", MODULE_NAME, c.Log, c.Logv)
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
