/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type HttpModule struct {
	enable   bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *HttpModule) prestart(backend *Backend) error {
	p.upstream = backend.Conf.Configer.Str(C_GRPC_ADDR)
	p.address = backend.Conf.Configer.Str(C_HTTP_ADDR)
	p.enable, _ = backend.Conf.Configer.Bool(C_HTTP_ENABLE)
	return nil
}

func (p *HttpModule) start(backend *Backend) error {
	if !p.enable {
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := Gateway(p.ctx, mux, p.upstream)
	if err != nil {
		return nil
	}

	server := &http.Server{Addr: p.address, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			p.cancel()
		}
		return
	}()

	go func() {
		<-p.ctx.Done()
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		server.Shutdown(ctx)
	}()
	return nil
}

func (p *HttpModule) stop(backend *Backend) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *HttpModule) reload(backend *Backend) error {
	if p.enable {
		p.stop(backend)
		time.Sleep(time.Second)
	}
	p.prestart(backend)
	return p.start(backend)
}
