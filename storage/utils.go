/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func renderDataJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, specs.Dto{Msg: "success", Data: data})
}

func renderMsgJson(w http.ResponseWriter, msg string) {
	renderJson(w, map[string]string{"msg": msg})
}

func autoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		renderMsgJson(w, err.Error())
		return
	}
	renderDataJson(w, data)
}

func dictedTagstring(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}
	s = strings.Replace(s, " ", "", -1)

	tag_dict := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {
			tag_dict[tag_pair[0]] = tag_pair[1]
		}
	}
	return tag_dict
}

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
func key2filename(baseDir string, key string) string {
	return fmt.Sprintf("%s/%s/%s.rrd", baseDir, key[0:2], key)
}

/*
func sortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	if size == 1 {
		for k, v := range tags {
			return fmt.Sprintf("%s=%s", k, v)
		}
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	ret := make([]string, size)
	for j, key := range keys {
		ret[j] = fmt.Sprintf("%s=%s", key, tags[key])
	}

	return strings.Join(ret, ",")
}
*/

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
			commitCache(false)

			glog.Info("rrd data commit complete")
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
