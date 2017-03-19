/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"io/ioutil"
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
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
}

func echo_handle(w http.ResponseWriter, req *http.Request) {
	if ctx, err := ioutil.ReadAll(req.Body); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(ctx)
	}
}

func (p *httpModule) httpRoutes() {
	p.httpMux.HandleFunc("/echo", echo_handle)
}

func (p *httpModule) prestart(L *Transfer) error {
	p.httpMux = http.NewServeMux()
	p.httpRoutes()
	return nil
}

func (p *httpModule) start(L *Transfer) error {
	enable, _ := L.Conf.Configer.Bool(falcon.C_HTTP_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "http.Start warning, not enabled")
		return nil
	}

	network, addr := falcon.ParseAddr(L.Conf.Configer.Str(falcon.C_HTTP_ADDR))
	if addr == "" {
		return falcon.ErrParam
	}

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof(MODULE_NAME+"%s httpStart listening %s", L.Conf.Name, addr)

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

func (p *httpModule) stop(L *Transfer) error {
	p.httpListener.Close()
	return nil
}

func (p *httpModule) reload(L *Transfer) error {
	return nil
}
