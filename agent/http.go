/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"encoding/json"
	"net"
	"net/http"
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
	httpListener  *net.TCPListener
	httpMux       *http.ServeMux
	appUpdateChan chan *[]*falcon.MetaData
}

func (p *httpModule) push(w http.ResponseWriter, req *http.Request) {
	if req.ContentLength == 0 {
		http.Error(w, "body is blank", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var meta []*falcon.MetaData
	err := decoder.Decode(&meta)
	if err != nil {
		http.Error(w, "connot decode body", http.StatusBadRequest)
		return
	}

	p.appUpdateChan <- &meta
	w.Write([]byte("success"))
}

func (p *httpModule) httpRoutes() {
	p.httpMux.HandleFunc("/push", p.push)
}

func (p *httpModule) prestart(agent *Agent) error {
	p.appUpdateChan = agent.appUpdateChan
	p.httpMux = http.NewServeMux()
	p.httpRoutes()
	return nil
}

func (p *httpModule) start(agent *Agent) error {
	enable, _ := agent.Conf.Configer.Bool(falcon.C_HTTP_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "http.Start warning, not enabled")
		return nil
	}

	addr := agent.Conf.Configer.Str(falcon.C_HTTP_ADDR)
	if addr == "" {
		return falcon.ErrParam
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof(MODULE_NAME+"%s http listening %s", agent.Conf.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.httpListener = ln.(*net.TCPListener)
	go s.Serve(tcpKeepAliveListener{p.httpListener})

	return nil
}

func (p *httpModule) stop(agent *Agent) error {
	p.httpListener.Close()
	return nil
}

func (p *httpModule) reload(agent *Agent) error {
	return nil
}
