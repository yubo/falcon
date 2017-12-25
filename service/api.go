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
	in *falcon.GetRequest) (resp *falcon.GetResponse, err error) {
	var (
		e *cacheEntry
	)

	statsInc(ST_RPC_SERV_QUERY, 1)

	resp = &falcon.GetResponse{
		Endpoint: in.Endpoint,
		Metric:   in.Metric,
		Type:     in.Type,
	}

	if e = p.service.cache.get(in.Key()); e == nil {
		err = falcon.ErrNoExits
		return
	}

	// TODO: get tsdb, rrd ?
	resp.Dps = e._getData(CACHE_SIZE)

	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Dps))

	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	total, errors := putItems(p.service, in.Items)
	return &falcon.PutResponse{Total: int32(total), Errors: int32(errors)}, nil
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

func putItems(service *Service, items []*falcon.Item) (total, errors int) {
	var (
		err  error
		e    *cacheEntry
		item *falcon.Item
	)

	if total = len(items); total == 0 {
		return
	}

	glog.V(4).Infof(MODULE_NAME+"recv %d", total)
	statsInc(ST_RPC_SERV_RECV, 1)
	statsInc(ST_RPC_SERV_RECV_ITEM, total)

	for i := 0; i < total; i++ {
		if item = items[i]; item == nil {
			errors++
			continue
		}
		key := item.Key()

		e = service.cache.get(key)
		if e == nil {
			e, err = service.createEntry(key, item)
			if err != nil {
				errors++
				continue
			}
		}

		e.put(item)
	}
	return
}
