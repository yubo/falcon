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
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiModule struct {
	disable bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	putChan chan []*falcon.Item
}

func (p *ApiModule) Get(ctx context.Context,
	in *falcon.GetRequest) (*falcon.GetResponse, error) {

	return &falcon.GetResponse{}, nil
}

func (p *ApiModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	p.putChan <- in.Items
	return &falcon.PutResponse{int32(len(in.Items)), 0}, nil
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
	service.RegisterServiceServer(server, p)

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
