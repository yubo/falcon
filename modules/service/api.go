/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type apiModule struct {
	disable bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	service *Service
}

func (p *apiModule) Get(ctx context.Context,
	in *GetRequest) (res *GetResponse, err error) {

	glog.V(3).Infof("%s rx get len(Keys) %d", MODULE_NAME, len(in.Keys))
	res, err = p.service.tsdb.get(in)

	statsInc(ST_RX_GET_ITERS, 1)
	statsInc(ST_RX_GET_ITEMS, len(res.Data))
	return
}

func (p *apiModule) Put(ctx context.Context,
	in *PutRequest) (res *PutResponse, err error) {

	//TODO
	glog.V(4).Infof("%s rx put %v", MODULE_NAME, len(in.Data))
	res, err = p.service.tsdb.put(in)
	if err != nil {
		glog.V(4).Infof("%s rx put err %v", MODULE_NAME, err)
	}
	for i := int32(0); i < res.N; i++ {
		p.service.cache.put(in.Data[i])
	}

	statsInc(ST_RX_PUT_ITERS, 1)
	statsInc(ST_RX_PUT_ITEMS, int(res.N))
	statsInc(ST_RX_PUT_ERR_ITEMS, len(in.Data)-int(res.N))
	return
}

func (p *apiModule) GetStats(ctx context.Context, in *Empty) (*Stats, error) {
	return &Stats{Counter: statsGets()}, nil
}

func (p *apiModule) GetStatsName(ctx context.Context, in *Empty) (*StatsName, error) {
	return &StatsName{CounterName: statsCounterName}, nil
}

func (p *apiModule) prestart(service *Service) error {
	p.address = service.Conf.ApiAddr
	p.disable = core.AddrIsDisable(p.address)
	p.service = service
	return nil
}

func (p *apiModule) start(service *Service) error {

	if p.disable {
		glog.Info(MODULE_NAME + "api disable")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(core.CleanSockFile(core.ParseAddr(p.address)))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterServiceServer(server, p)

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

func (p *apiModule) stop(service *Service) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *apiModule) reload(service *Service) error {
	return nil

	if !p.disable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}
