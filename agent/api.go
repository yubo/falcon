/*
	address  string
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
*/
package agent

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type putContext struct {
	items []*Item
	done  chan *PutResponse
}

type ApiModule struct {
	disable bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	putChan chan *putContext
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (*PutResponse, error) {

	put := &putContext{
		items: in.Items,
		done:  make(chan *PutResponse),
	}

	p.putChan <- put
	resp := <-put.done
	return resp, nil
}

func (p *ApiModule) prestart(agent *Agent) error {
	p.address = agent.Conf.Configer.Str(C_API_ADDR)
	glog.Info(MODULE_NAME + "address " + p.address)
	p.disable = falcon.AddrIsDisable(p.address)
	p.putChan = agent.appPutChan
	return nil
}

func (p *ApiModule) start(agent *Agent) error {
	glog.V(3).Info(MODULE_NAME + "api start")
	if p.disable {
		return nil
	}
	glog.V(3).Info(MODULE_NAME + "api start")

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterAgentServer(server, p)

	// Register reflection service on gRPC server.
	reflection.Register(server)
	go func() {
		if err := server.Serve(ln); err != nil {
			p.cancel()
		}
	}()

	go func() {
		<-p.ctx.Done()
		server.Stop()
	}()

	return nil
}

func (p *ApiModule) stop(agent *Agent) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(agent *Agent) error {
	if !p.disable {
		p.stop(agent)
		time.Sleep(time.Second)
	}
	p.prestart(agent)
	return p.start(agent)
}
