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
	"golang.org/x/net/context"
)

type HttpModule struct {
	enable   bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *HttpModule) prestart(transfer *Transfer) error {
	p.upstream = transfer.Conf.Configer.Str(C_GRPC_ADDR)
	p.address = transfer.Conf.Configer.Str(C_HTTP_ADDR)
	p.enable, _ = transfer.Conf.Configer.Bool(C_HTTP_ENABLE)
	return nil
}

func (p *HttpModule) start(transfer *Transfer) error {
	if !p.enable {
		glog.Info(MODULE_NAME + "http.Start not enabled")
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

func (p *HttpModule) stop(transfer *Transfer) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *HttpModule) reload(transfer *Transfer) error {
	if p.enable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
