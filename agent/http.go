/*
 * Copyright 2016 falcon Author. All rights reserved.
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

type HttpModule struct {
	httpListener  *net.TCPListener
	httpMux       *http.ServeMux
	appUpdateChan chan *[]*falcon.MetaData
}

func (p *HttpModule) push(w http.ResponseWriter, req *http.Request) {
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

func (p *HttpModule) httpRoutes() {
	p.httpMux.HandleFunc("/push", p.push)
}

func (p *HttpModule) prestart(agent *Agent) error {
	p.appUpdateChan = agent.appUpdateChan
	p.httpMux = http.NewServeMux()
	p.httpRoutes()
	return nil
}

func (p *HttpModule) start(agent *Agent) error {
	enable, _ := agent.Conf.Configer.Bool(C_HTTP_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "http.Start warning, not enabled")
		return nil
	}

	network, addr := falcon.ParseAddr(agent.Conf.Configer.Str(C_HTTP_ADDR))
	if addr == "" {
		return falcon.ErrParam
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof(MODULE_NAME+"%s http listening %s", agent.Conf.Name, addr)

	ln, err := net.Listen(network, addr)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	if network == "tcp" {
		p.httpListener = ln.(*net.TCPListener)
		go s.Serve(tcpKeepAliveListener{p.httpListener})
	} else {
		go s.Serve(ln)
	}

	return nil
}

func (p *HttpModule) stop(agent *Agent) error {
	p.httpListener.Close()
	return nil
}

func (p *HttpModule) reload(agent *Agent) error {
	return nil
}
