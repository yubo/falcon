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
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiModule struct {
	disable bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	service *Service
}

func (p *ApiModule) Get(ctx context.Context,
	in *GetRequest) (resp *GetResponse, err error) {

	statsInc(ST_RPC_SERV_QUERY, 1)

	resp = &GetResponse{Key: in.Key}

	resp.Dps, err = p.service.shard.get(in)
	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Dps))
	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (resp *PutResponse, err error) {

	n := 0
	size := len(in.Items)

	glog.V(4).Infof(MODULE_NAME+"recv %d", size)
	statsInc(ST_RPC_SERV_RECV, 1)
	statsInc(ST_RPC_SERV_RECV_ITEM, size)

	for i := 0; i < size; i++ {
		item := in.Items[i]
		if item == nil {
			continue
		}

		if _, err = p.service.shard.put(item); err != nil {
			n++
			continue
		}
	}

	return &PutResponse{N: int32(n)}, nil
}

func (p *ApiModule) prestart(service *Service) error {
	p.address = service.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.service = service
	return nil
}

func (p *ApiModule) start(service *Service) error {

	if p.disable {
		glog.Info(MODULE_NAME + "api disable")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
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

func (p *ApiModule) stop(service *Service) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(service *Service) error {
	if !p.disable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}
