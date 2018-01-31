/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
)

type ApiGwModule struct {
	disable  bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *ApiGwModule) prestart(transfer *Transfer) error {
	return nil
}

func (p *ApiGwModule) start(transfer *Transfer) error {
	conf := &transfer.Conf.Configer
	p.upstream = conf.Str(C_API_ADDR)
	p.address = conf.Str(C_HTTP_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)

	if p.disable {
		glog.Info(MODULE_NAME + "http disable")
		return nil
	}
	glog.V(4).Infof("%s http addr %s", MODULE_NAME, p.address)

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := falcon.Gateway(service.RegisterServiceHandlerFromEndpoint, p.ctx, mux, p.upstream)
	if err != nil {
		return err
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

func (p *ApiGwModule) stop(transfer *Transfer) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiGwModule) reload(transfer *Transfer) error {
	return nil

	if !p.disable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
