/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
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

func (p *Handoff) httpRoutes() {
}

func (p *Handoff) httpStart() {
	if !p.Http {
		return
	}

	addr := p.HttpAddr
	if addr == "" {
		return
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof("%s httpStart listening %s", p.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(err)
	}

	p.httpListener = ln.(*net.TCPListener)

	p.httpRoutes()

	go s.Serve(tcpKeepAliveListener{p.httpListener})
}

func (p *Handoff) httpStop() error {
	p.httpListener.Close()
	return nil
}
