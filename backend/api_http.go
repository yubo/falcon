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
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

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

type httpModule struct {
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
	b            *Backend
}

func (p *httpModule) count_handler(w http.ResponseWriter, r *http.Request) {
	ts := p.b.timeNow() - 3600
	count := 0

	for _, v := range p.b.cache.hash {
		if int64(v.lastTs) > ts {
			count++
		}
	}
	w.Write([]byte(fmt.Sprintf("%d\n", count)))
}

func (p *httpModule) recv_hanlder(w http.ResponseWriter, r *http.Request) {
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

	item := &falcon.MetaData{
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

	p.b.handleItems([]*falcon.RrdItem{gitem})
	renderDataJson(w, "ok")
}

func (p *httpModule) recv2_handler(w http.ResponseWriter, r *http.Request) {
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

	item := &falcon.MetaData{
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

	p.b.handleItems([]*falcon.RrdItem{gitem})
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
	renderDataJson(w, statsHandle())
}

func (p *httpModule) httpRoutes() {
	p.httpMux.HandleFunc("/count", p.count_handler)

	p.httpMux.HandleFunc("/api/recv/", p.recv_hanlder)

	p.httpMux.HandleFunc("/v2/api/recv", p.recv2_handler)

	p.httpMux.HandleFunc("/index/updateAll", updateAll_handler)

	// 获取索引全量更新的并行数
	p.httpMux.HandleFunc("/index/updateAll/concurrent", updateAll_concurrent_handler)

	p.httpMux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	p.httpMux.HandleFunc("/version",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(falcon.VERSION))
		})

	p.httpMux.HandleFunc("/config",
		func(w http.ResponseWriter, r *http.Request) {
			renderDataJson(w, p)
		})

	p.httpMux.HandleFunc("/stat", stat_handler)

}

func (p *httpModule) prestart(b *Backend) error {
	p.httpMux = http.NewServeMux()
	p.httpRoutes()
	p.b = b
	return nil
}

func (p *httpModule) start(b *Backend) error {
	enable, _ := b.Conf.Configer.Bool(falcon.C_HTTP_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "http not enabled")
		return nil
	}

	addr := b.Conf.Configer.Str(falcon.C_HTTP_ADDR)
	if addr == "" {
		return falcon.ErrParam
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof(MODULE_NAME+"%s http listening %s", b.Conf.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.httpListener = ln.(*net.TCPListener)
	go s.Serve(tcpKeepAliveListener{p.httpListener})
	return nil
}

func (p *httpModule) stop(b *Backend) error {
	p.httpListener.Close()
	return nil
}

func (p *httpModule) reload(b *Backend) error {
	return nil
}
