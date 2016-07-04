/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, specs.Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	RenderDataJson(w, data)
}

func registerEventChan(e *specs.RoutineEvent) {
	glog.V(3).Infof("register exit chan '%s'", e.Name)
	appEvents = append(appEvents, e)
}

func startSignal() {
	pid := os.Getpid()
	sigs := make(chan os.Signal, 1)
	glog.Infof("[%d] register signal notify", pid)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		glog.Infof("recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			glog.Info("gracefull shut down")
			atomic.StoreUint32(&appStatus, specs.APP_STATUS_EXIT)
			event := specs.REvent{
				Method: specs.ROUTINE_EVENT_M_EXIT,
				Done:   make(chan error),
			}
			for i := len(appEvents) - 1; i >= 0; i-- {
				glog.V(3).Infof("send exit signal to %s",
					appEvents[i].Name)
				appEvents[i].E <- event
				if err := <-event.Done; err != nil {
					glog.Info(err)
				}
				glog.V(3).Infof("%s done", appEvents[i].Name)
			}

			glog.Infof("pid:%d exit", pid)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Info("relod shut down")
			atomic.StoreUint32(&appStatus, specs.APP_STATUS_RELOAD)
			event := specs.REvent{
				Method: specs.ROUTINE_EVENT_M_RELOAD,
				Done:   make(chan error),
			}
			for i := len(appEvents) - 1; i >= 0; i-- {
				glog.V(3).Infof("send reload signal to %s",
					appEvents[i].Name)
				appEvents[i].E <- event
				if err := <-event.Done; err != nil {
					glog.Info(err)
				}
				glog.V(3).Infof("%s done", appEvents[i].Name)
			}
			atomic.StoreUint32(&appStatus, specs.APP_STATUS_RUNING)
		}
	}
}

func sortTags(s string) string {
	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		return ""
	}

	tags := strings.Split(s, ",")
	sort.Strings(tags)
	return strings.Join(tags, ",")
}
