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

type ProcEvent struct {
	Method int
	Done   chan error
}

type procEvent struct {
	name  string
	event chan ProcEvent
}

type hookFunc struct {
	name string
	fun  func(interface{})
	arg  interface{}
}

type Process struct {
	PidFile   string
	Pid       int
	status    uint32
	events    []*procEvent
	postHooks []hookFunc
}

func NewProcess(pidfile string) *Process {
	return &Process{
		PidFile:   pidfile,
		Pid:       os.Getpid(),
		status:    APP_STATUS_PENDING,
		events:    []*procEvent{},
		postHooks: []hookFunc{},
	}
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

func (p *Process) RegisterEvent(name string, ch chan ProcEvent) {
	glog.V(3).Infof("register process event chan '%s'", name)
	p.events = append(p.events, &procEvent{name: name, event: ch})
}

func (p *Process) RegisterPostHook(name string, fun func(interface{}), arg interface{}) {
	glog.V(3).Infof("register post hook '%s'", name)
	p.postHooks = append(p.postHooks, hookFunc{name: name, fun: fun, arg: arg})
}

func (p *Process) StartSignal() {
	sigs := make(chan os.Signal, 1)
	glog.Infof("[%d] register signal notify", p.Pid)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	atomic.StoreUint32(&p.status, APP_STATUS_RUNING)

	for {
		s := <-sigs
		glog.Infof("recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Info("gracefull shut down")
			atomic.StoreUint32(&p.status, APP_STATUS_EXIT)
			event := ProcEvent{
				Method: ROUTINE_EVENT_M_EXIT,
				Done:   make(chan error),
			}
			for i := len(p.events) - 1; i >= 0; i-- {
				glog.V(3).Infof("send exit signal to %s",
					p.events[i].name)
				p.events[i].event <- event
				if err := <-event.Done; err != nil {
					glog.Info(err)
				}
				glog.V(3).Infof("%s done", p.events[i].name)
			}
			glog.V(3).Infof("begin post hooks start")

			for i := len(p.postHooks) - 1; i >= 0; i-- {
				glog.V(3).Infof("call hook %d %s", i, p.postHooks[i].name)
				p.postHooks[i].fun(p.postHooks[i].arg)
			}

			glog.Infof("pid:%d exit", p.Pid)
			os.Remove(p.PidFile)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Info("relod shut down")
			atomic.StoreUint32(&p.status, APP_STATUS_RELOAD)
			event := ProcEvent{
				Method: ROUTINE_EVENT_M_RELOAD,
				Done:   make(chan error),
			}
			for i := len(p.events) - 1; i >= 0; i-- {
				glog.V(3).Infof("send reload signal to %s",
					p.events[i].name)
				p.events[i].event <- event
				if err := <-event.Done; err != nil {
					glog.Info(err)
				}
				glog.V(3).Infof("%s done", p.events[i].name)
			}
			atomic.StoreUint32(&p.status, APP_STATUS_RUNING)
		}
	}
}
