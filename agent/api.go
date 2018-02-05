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
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PutRequest struct {
	Dps  []*transfer.DataPoint
	Done chan *transfer.PutResponse
}

type ApiModule struct {
	disable bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	putChan chan *PutRequest
}

func (p *ApiModule) Put(ctx context.Context, in *transfer.PutRequest) (*transfer.PutResponse,
	error) {

	glog.V(5).Infof("%s rx put %v", MODULE_NAME, in)
	put := &PutRequest{
		Dps:  in.Data,
		Done: make(chan *transfer.PutResponse),
	}

	p.putChan <- put
	resp := <-put.Done
	return resp, nil
}

func (p *ApiModule) Get(ctx context.Context, in *transfer.GetRequest) (*transfer.GetResponse,
	error) {
	return nil, falcon.ErrUnsupported
}

func (p *ApiModule) GetStats(ctx context.Context, in *transfer.Empty) (*transfer.Stats,
	error) {
	return &transfer.Stats{Counter: statsGets()}, nil
}

func (p *ApiModule) GetStatsName(ctx context.Context, in *transfer.Empty) (*transfer.StatsName,
	error) {
	return &transfer.StatsName{CounterName: statsCounterName}, nil
}

func (p *ApiModule) prestart(agent *Agent) error {
	p.address = agent.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.putChan = agent.PutChan
	return nil
}

func (p *ApiModule) start(agent *Agent) error {
	if p.disable {
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.CleanSockFile(falcon.ParseAddr(p.address)))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	transfer.RegisterTransferServer(server, p)

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
	return nil

	if !p.disable {
		p.stop(agent)
		time.Sleep(time.Second)
	}
	p.prestart(agent)
	return p.start(agent)
}
