/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcModule struct {
	enable  bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	b       *Backend
}

func (p *GrpcModule) Get(ctx context.Context,
	in *falcon.GetRequest) (resp *falcon.GetResponse, err error) {
	var (
		e *cacheEntry
	)

	statsInc(ST_RPC_SERV_QUERY, 1)

	resp = &falcon.GetResponse{
		Vs:       []*falcon.RRDData{},
		Endpoint: in.Endpoint,
		Metric:   in.Metric,
		Type:     in.Type,
	}

	if e = p.b.cache.get(in.Csum()); e == nil {
		err = falcon.ErrNoExits
		return
	}

	// TODO: get tsdb, rrd ?
	resp.Vs, _ = e._getData(e.commitId, e.dataId)

	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))

	return
}

func (p *GrpcModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	total, errors := p.b.handleItems(in.Items)
	return &falcon.PutResponse{Total: int32(total), Errors: int32(errors)}, nil
}

func (p *GrpcModule) prestart(backend *Backend) error {
	p.enable, _ = backend.Conf.Configer.Bool(C_GRPC_ENABLE)
	p.address = backend.Conf.Configer.Str(C_GRPC_ADDR)
	p.b = backend
	return nil
}

func (p *GrpcModule) start(backend *Backend) error {

	if !p.enable {
		glog.Info(MODULE_NAME + "grpc.Start not enabled")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterBackendServer(server, &GrpcModule{})

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

func (p *GrpcModule) stop(backend *Backend) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *GrpcModule) reload(backend *Backend) error {
	if p.enable {
		p.stop(backend)
		time.Sleep(time.Second)
	}
	p.prestart(backend)
	return p.start(backend)
}
