/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"

	"golang.org/x/net/context"
)

type apiGwModule struct {
	disable  bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *apiGwModule) prestart(service *Service) error {
	p.upstream = service.Conf.ApiAddr
	p.address = service.Conf.HttpAddr
	p.disable = core.AddrIsDisable(p.address)
	return nil
}

func (p *apiGwModule) start(service *Service) error {
	if p.disable {
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := core.Gateway(RegisterServiceHandlerFromEndpoint, p.ctx, mux, p.upstream)
	if err != nil {
		return nil
	}

	server := &http.Server{Addr: p.address, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			glog.Errorf("%s ListenAndServ %s err %v", MODULE_NAME, p.address, err)
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

func (p *apiGwModule) stop(service *Service) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *apiGwModule) reload(service *Service) error {
	return nil

	if !p.disable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}
