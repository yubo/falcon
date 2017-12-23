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
	enable  bool
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
		Vs:       []*falcon.RRDData{},
		Endpoint: in.Endpoint,
		Metric:   in.Metric,
		Type:     in.Type,
	}

	if e = p.service.cache.get(in.Csum()); e == nil {
		err = falcon.ErrNoExits
		return
	}

	// TODO: get tsdb, rrd ?
	resp.Vs, _ = e._getData(e.commitId, e.dataId)

	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))

	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	total, errors := handleItems(p.service, in.Items)
	return &falcon.PutResponse{Total: int32(total), Errors: int32(errors)}, nil
}

func (p *ApiModule) prestart(service *Service) error {
	p.enable, _ = service.Conf.Configer.Bool(C_GRPC_ENABLE)
	p.address = service.Conf.Configer.Str(C_GRPC_ADDR)
	p.service = service
	return nil
}

func (p *ApiModule) start(service *Service) error {

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
	RegisterServiceServer(server, &ApiModule{})

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
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(service *Service) error {
	if p.enable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}

func handleItems(service *Service, items []*falcon.Item) (total, errors int) {
	var (
		err error
		e   *cacheEntry
	)

	if total = len(items); total == 0 {
		return
	}

	glog.V(4).Infof(MODULE_NAME+"recv %d", total)
	statsInc(ST_RPC_SERV_RECV, 1)
	statsInc(ST_RPC_SERV_RECV_ITEM, total)

	for i := 0; i < total; i++ {
		if items[i] == nil {
			errors++
			continue
		}
		key := items[i].Csum()

		e = service.cache.get(key)
		if e == nil {
			e, err = service.createEntry(key, items[i])
			if err != nil {
				errors++
				continue
			}
		}

		/*
			if DATA_TIMESTAMP_REGULATE {
				items[i].Ts = items[i].Ts -
					items[i].Ts%int64(items[i].Step)
			}

			if items[i].Ts <= e.lastTs || items[i].Ts <= 0 {
				errors++
				continue
			}
		*/

		e.put(items[i])
	}
	return
}
