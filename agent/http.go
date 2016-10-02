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
	"github.com/yubo/falcon/specs"
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

func (p *Agent) httpRoutes() {
	p.httpMux.HandleFunc("/push", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var meta []*specs.MetaData
		err := decoder.Decode(&meta)
		if err != nil {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}

		p.appUpdateChan <- &meta
		w.Write([]byte("success"))
	})
}

func (p *Agent) httpStart() {
	if !p.Http {
		glog.Info("http.Start warning, not enabled")
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
	glog.Infof("%s http listening %s", p.Name, addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(err)
	}

	p.httpListener = ln.(*net.TCPListener)
	go s.Serve(tcpKeepAliveListener{p.httpListener})
}

func (p *Agent) httpStop() error {
	p.httpListener.Close()
	return nil
}
