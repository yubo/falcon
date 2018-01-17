/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"net/http"
	"time"

	"github.com/golang/glog"
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

func (p *ApiGwModule) prestart(transfer *Alarm) error {
	p.upstream = transfer.Conf.Configer.Str(C_API_ADDR)
	p.address = transfer.Conf.Configer.Str(C_HTTP_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	return nil
}

func (p *ApiGwModule) start(transfer *Alarm) error {
	if p.disable {
		glog.Info(MODULE_NAME + "http disable")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := falcon.Gateway(RegisterAlarmHandlerFromEndpoint, p.ctx, mux, p.upstream)
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

func (p *ApiGwModule) stop(transfer *Alarm) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiGwModule) reload(transfer *Alarm) error {
	return nil

	if !p.disable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
