/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func echo_handle(w http.ResponseWriter, req *http.Request) {
	if ctx, err := ioutil.ReadAll(req.Body); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(ctx)
	}
}

func (p *Ctrl) httpRoutes() {
	p.httpMux.HandleFunc("/echo", echo_handle)
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

func (p *Ctrl) httpStart() {
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
	glog.Infof(MODULE_NAME+"%s httpStart listening %s", p.Conf.Params.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(MODULE_NAME, err)
	}

	p.httpListener = ln.(*net.TCPListener)
	go s.Serve(tcpKeepAliveListener{p.httpListener})
}

func (p *Ctrl) httpStop() error {
	p.httpListener.Close()
	return nil
}
