/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"
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
	renderJson(w, Dto{Msg: "success", Data: data})
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

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
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

func httpStart() {
	if !config().Http {
		log.Println("http.Start warning, not enabled")
		return
	}

	addr := config().HttpAddr
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	log.Println("http listening", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	l := ln.(*net.TCPListener)
	registerExitChans(http_exit)

	go s.Serve(tcpKeepAliveListener{l})

	go func() {
		select {
		case done := <-http_exit:
			log.Println("http recv sigout and exit...")
			l.Close()
			done <- nil
			return
		}
	}()

}
