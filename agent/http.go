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

func (p *Agent) push_handle(w http.ResponseWriter, req *http.Request) {
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

func (p *Agent) httpRoutes() {
	p.httpMux.HandleFunc("/push", p.push_handle)
}

func (p *Agent) httpStart() {
	if !p.Conf.Params.Http {
		glog.Info(MODULE_NAME + "http.Start warning, not enabled")
		return
	}

	addr := p.Conf.Params.HttpAddr
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof(MODULE_NAME+"%s http listening %s", p.Conf.Params.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.httpListener = ln.(*net.TCPListener)
	go s.Serve(tcpKeepAliveListener{p.httpListener})
}

func (p *Agent) httpStop() error {
	p.httpListener.Close()
	return nil
}
