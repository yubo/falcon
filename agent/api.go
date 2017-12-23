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

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiModule struct {
	enable     bool
	ctx        context.Context
	cancel     context.CancelFunc
	address    string
	updateChan chan []*falcon.Item
}

func (p *ApiModule) Get(ctx context.Context,
	in *falcon.GetRequest) (*falcon.GetResponse, error) {

	return &falcon.GetResponse{}, nil
}

func (p *ApiModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	p.updateChan <- in.Items
	return &falcon.PutResponse{int32(len(in.Items)), 0}, nil
}

func (p *ApiModule) prestart(agent *Agent) error {
	p.enable, _ = agent.Conf.Configer.Bool(C_GRPC_ENABLE)
	p.address = agent.Conf.Configer.Str(C_GRPC_ADDR)
	p.updateChan = agent.appUpdateChan
	return nil
}

func (p *ApiModule) start(agent *Agent) error {

	if !p.enable {
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service.RegisterServiceServer(server, &ApiModule{})

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
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(agent *Agent) error {
	if p.enable {
		p.stop(agent)
		time.Sleep(time.Second)
	}
	p.prestart(agent)
	return p.start(agent)
}
