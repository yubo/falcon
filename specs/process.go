/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package specs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
)

type procEvent struct {
}

type hookFunc struct {
	name string
	fun  func(interface{})
	arg  interface{}
}

type Module interface {
	Init() error
	Start() error
	Stop() error
	Reload() error
	Signal(os.Signal) error
	String() string
	Desc() string
}
type Process struct {
	PidFile   string
	Pid       int
	status    uint32
	events    []*procEvent
	postHooks []hookFunc
	modules   []Module
}

func NewProcess(pidfile string, modules []Module) *Process {
	p := &Process{
		PidFile: pidfile,
		Pid:     os.Getpid(),
		status:  APP_STATUS_PENDING,
	}
	p.modules = make([]Module, len(modules))
	copy(p.modules, modules)
	return p
}

func (p *Process) Status() uint32 {
	return atomic.LoadUint32(&p.status)
}

func (p *Process) Kill(sig syscall.Signal) error {

	if _, err := os.Stat(p.PidFile); os.IsNotExist(err) {
		return errors.New("pid file not exist")
	}
	if data, err := ioutil.ReadFile(p.PidFile); err != nil {
		return err
	} else {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err != nil {
			return err
		} else {
			return syscall.Kill(pid, sig)
		}
	}
}

func (p *Process) Check() error {
	if _, err := os.Stat(p.PidFile); os.IsNotExist(err) {
		return nil
	}
	if data, err := ioutil.ReadFile(p.PidFile); err != nil {
		return err
	} else {
		pid := strings.TrimSpace(string(data))
		if _, err := os.Stat(fmt.Sprintf("/proc/%s",
			pid)); os.IsNotExist(err) {
			return nil
		} else {
			return fmt.Errorf("proccess %s exist", pid)
		}
	}
}

func (p *Process) Save() error {
	return ioutil.WriteFile(p.PidFile,
		[]byte(fmt.Sprintf("%d", p.Pid)), 0644)
}

func (p *Process) Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	atomic.StoreUint32(&p.status, APP_STATUS_RUNING)

	for _, m := range p.modules {
		m.Init()
		m.Start()
	}

	glog.Infof(MODULE_NAME+"[%d] register signal notify", p.Pid)

	for {
		s := <-sigs
		glog.Infof(MODULE_NAME+"recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			pidfile := fmt.Sprintf("%s.%d", p.PidFile, p.Pid)
			glog.Info(MODULE_NAME + "exiting")
			atomic.StoreUint32(&p.status, APP_STATUS_EXIT)
			os.Rename(p.PidFile, pidfile)

			for _, m := range p.modules {
				m.Stop()
			}

			glog.Infof(MODULE_NAME+"pid:%d exit", p.Pid)
			os.Remove(pidfile)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Info(MODULE_NAME + "reload")
			atomic.StoreUint32(&p.status, APP_STATUS_RELOAD)
			for _, m := range p.modules {
				m.Reload()
			}
			atomic.StoreUint32(&p.status, APP_STATUS_RUNING)
		default:
			for _, m := range p.modules {
				m.Signal(s)
			}
		}
	}
}
