/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
)

type ApiGwModule struct {
	disable  bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *ApiGwModule) prestart(agent *Agent) error {
	p.upstream = agent.Conf.Configer.Str(C_API_ADDR)
	p.address = agent.Conf.Configer.Str(C_HTTP_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	return nil
}

func (p *ApiGwModule) start(agent *Agent) error {
	if p.disable {
		return nil
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := falcon.Gateway(transfer.RegisterTransferHandlerFromEndpoint, p.ctx, mux, p.upstream)
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

func (p *ApiGwModule) stop(agent *Agent) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiGwModule) reload(agent *Agent) error {
	return nil

	if !p.disable {
		p.stop(agent)
		time.Sleep(time.Second)
	}
	p.prestart(agent)
	return p.start(agent)
}
