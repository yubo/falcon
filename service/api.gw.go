/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"net/http"
	"time"

	"github.com/yubo/falcon"

	"golang.org/x/net/context"
)

type ApiGwModule struct {
	disable  bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *ApiGwModule) prestart(service *Service) error {
	p.upstream = service.Conf.Configer.Str(C_API_ADDR)
	p.address = service.Conf.Configer.Str(C_HTTP_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	return nil
}

func (p *ApiGwModule) start(service *Service) error {
	if p.disable {
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := falcon.Gateway(RegisterServiceHandlerFromEndpoint, p.ctx, mux, p.upstream)
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

func (p *ApiGwModule) stop(service *Service) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiGwModule) reload(service *Service) error {
	if !p.disable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}
