/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/specs"
)

var (
	httpEvent  chan specs.ProcEvent
	httpConfig BackendOpts
)

func count_handler(w http.ResponseWriter, r *http.Request) {
	ts := timeNow() - 3600
	count := 0

	for _, v := range appCache.hash {
		if int64(v.e.lastTs) > ts {
			count++
		}
	}
	w.Write([]byte(fmt.Sprintf("%d\n", count)))
}

func recv_hanlder(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Path[len("/api/recv/"):]
	args := strings.Split(urlParam, "/")

	argsLen := len(args)
	if !(argsLen == 6 || argsLen == 7) {
		renderDataJson(w, "bad args")
		return
	}

	host := args[0]
	k := args[1]
	ts, _ := strconv.ParseInt(args[2], 10, 64)
	step, _ := strconv.ParseInt(args[3], 10, 32)
	typ := args[4]
	value, _ := strconv.ParseFloat(args[5], 64)
	tags := ""
	if argsLen == 7 {
		tags = args[6]
	}

	item := &specs.MetaData{
		Host:  host,
		Name:  k,
		Ts:    ts,
		Step:  step,
		Type:  typ,
		Value: value,
		Tags:  tags,
	}
	gitem, err := item.Rrd()
	if err != nil {
		renderDataJson(w, err)
		return
	}

	handleItems([]*specs.RrdItem{gitem})
	renderDataJson(w, "ok")
}

func recv2_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0 &&
		len(r.Form["v"]) > 0 && len(r.Form["ts"]) > 0 &&
		len(r.Form["step"]) > 0 && len(r.Form["type"]) > 0) {
		renderDataJson(w, "bad args")
		return
	}
	host := r.Form["e"][0]
	k := r.Form["m"][0]
	value, _ := strconv.ParseFloat(r.Form["v"][0], 64)
	ts, _ := strconv.ParseInt(r.Form["ts"][0], 10, 64)
	step, _ := strconv.ParseInt(r.Form["step"][0], 10, 32)
	typ := r.Form["type"][0]

	tags := r.Form["t"][0]

	item := &specs.MetaData{
		Host:  host,
		Name:  k,
		Ts:    ts,
		Step:  step,
		Type:  typ,
		Value: value,
		Tags:  tags,
	}
	gitem, err := item.Rrd()
	if err != nil {
		renderDataJson(w, err)
		return
	}

	handleItems([]*specs.RrdItem{gitem})
	renderDataJson(w, "ok")
}

func updateAll_handler(w http.ResponseWriter, r *http.Request) {
	//	go UpdateIndexAllByDefaultStep()
	renderDataJson(w, "ok")
}

func updateAll_concurrent_handler(w http.ResponseWriter, r *http.Request) {
	//renderDataJson(w, GetConcurrentOfUpdateIndexAll())
	renderDataJson(w, "ok")
}

func stat_handler(w http.ResponseWriter, r *http.Request) {
	renderDataJson(w, statHandle())
}

func httpRoutes() {
	http.HandleFunc("/count", count_handler)

	http.HandleFunc("/api/recv/", recv_hanlder)

	http.HandleFunc("/v2/api/recv", recv2_handler)

	http.HandleFunc("/index/updateAll", updateAll_handler)

	// 获取索引全量更新的并行数
	http.HandleFunc("/index/updateAll/concurrent", updateAll_concurrent_handler)

	http.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	http.HandleFunc("/version",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(falcon.VERSION))
		})

	http.HandleFunc("/config",
		func(w http.ResponseWriter, r *http.Request) {
			renderDataJson(w, appConfig)
		})

	http.HandleFunc("/stat", stat_handler)

}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func httpStart(config BackendOpts, p *specs.Process) {
	if !config.Http {
		glog.Info("http.Start warning, not enabled")
		return
	}

	httpConfig = config

	addr := httpConfig.HttpAddr
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof("http listening %s", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(err)
	}

	l := ln.(*net.TCPListener)
	p.RegisterEvent("http", httpEvent)

	go s.Serve(tcpKeepAliveListener{l})

	go func() {
		select {
		case event := <-httpEvent:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				l.Close()
				event.Done <- nil
				return
			}
		}
	}()

}
