/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package core

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"syscall"
	"time"

	"github.com/golang/glog"
)

type Process struct {
	Configer   *Configer
	ValuesFile string
	ConfigFile string
	Pid        int
	Status     uint32
}

const (
	APP_STATUS_INIT = iota
	APP_STATUS_PENDING
	APP_STATUS_RUNNING
	APP_STATUS_EXIT
	APP_STATUS_RELOAD
)

type Module interface {
	ReadConfig(config *Configer) error
	Prestart() error
	Start() error
	Stop() error
	Reload(config *Configer) error
	Signal(os.Signal) error
	Stats(config *Configer) (string, error)
}

var (
	modules     []Module
	moduleNames = make(map[string]bool)
)

func RegisterModule(m Module) error {
	if moduleNames[GetName(m)] {
		return ErrExist
	}
	modules = append(modules, m)
	return nil
}

func NewProcess(configFile, baseConf string, valueFiles []string) (*Process, error) {

	configer, err := NewConfiger(configFile, baseConf, valueFiles)
	if err != nil {
		return nil, err
	}

	err = configer.Parse()
	if err != nil {
		return nil, err
	}

	p := &Process{
		Configer: configer,
		Pid:      os.Getpid(),
		Status:   APP_STATUS_PENDING,
		//Module: make([]Module, len(c.Conf)),
	}

	return p, nil
}

func (p *Process) Kill(sig syscall.Signal) error {
	pidFile := p.Configer.GetStr("pid_file")
	if pid, err := ReadFileInt(pidFile); err != nil {
		return err
	} else {
		glog.Infof("kill %d %s\n", pid, sig)
		return syscall.Kill(pid, sig)
	}
}

func (p *Process) Check() error {
	pidFile := p.Configer.GetStr("pid_file")
	pid, err := ReadFileInt(pidFile)
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
	pidFile := p.Configer.GetStr("pid_file")
	return ioutil.WriteFile(pidFile,
		[]byte(fmt.Sprintf("%d", p.Pid)), 0644)
}

// only be called once
func (p *Process) Start() error {

	rand.Seed(time.Now().Unix())

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("readconfig: %s\n", GetName(modules[i]))
		err := modules[i].ReadConfig(p.Configer)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("prestart: %s\n", GetName(modules[i]))
		if err := modules[i].Prestart(); err != nil {
			return err
		}
	}

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("start: %s\n", GetName(modules[i]))
		if err := modules[i].Start(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Process) Stop() error {
	for i, n := 0, len(modules); i < n; i++ {
		glog.V(4).Infof("stop: %s\n", GetName(modules[i]))
		modules[n-i-1].Stop()
	}
	return nil
}

func (p *Process) Reload() error {
	err := p.Configer.Parse()
	if err != nil {
		return err
	}

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("reload: %s\n", GetName(modules[i]))
		if err := modules[i].Reload(p.Configer); err != nil {
			glog.Errorf("reload: %v\n", err)
		}
	}
	return nil
}

func (p *Process) Stats(module string) error {

	for i := 0; i < len(modules); i++ {
		err := modules[i].ReadConfig(p.Configer)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(modules); i++ {
		m := modules[i]
		if module == "all" || module == "" || module == GetName(m) {
			stats, err := m.Stats(p.Configer)
			if err != nil {
				fmt.Printf("%s %s\n", GetName(m), err.Error())
				return err
			}
			fmt.Printf("%s\n%s\n", GetName(m), stats)
		}
	}
	return nil
}

func (p *Process) Signal(s os.Signal) error {

	for i := 0; i < len(modules); i++ {
		modules[i].Signal(s)
	}
	return nil
}
