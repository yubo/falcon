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
	in *GetRequest) (res *GetResponse, err error) {

	res, err = p.service.tsdb.get(in)

	statsInc(ST_RX_GET_ITERS, 1)
	statsInc(ST_RX_GET_ITEMS, len(res.Dps))
	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (res *PutResponse, err error) {

	glog.V(5).Infof("%s rx put %v", MODULE_NAME, len(in.Items))
	res, err = p.service.tsdb.put(in)

	statsInc(ST_RX_PUT_ITERS, 1)
	statsInc(ST_RX_PUT_ITEMS, int(res.N))
	statsInc(ST_RX_PUT_ERR_ITEMS, len(in.Items)-int(res.N))
	return

	/*
		res = &PutResponse{}
		for i := 0; i < len(in.Items); i++ {
			item := in.Items[i]
			if item == nil {
				continue
			}

			if _, err = p.service.tsdb.put(item); err != nil {
				res.N++
				continue
			}
		}

		statsInc(ST_RX_PUT_ITERS, 1)
		statsInc(ST_RX_PUT_ITEMS, int(res.N))
		statsInc(ST_RX_PUT_ERR_ITEMS, len(in.Items)-int(res.N))
	*/
}

func (p *ApiModule) GetStats(ctx context.Context, in *Empty) (*Stats, error) {
	return &Stats{Counter: statsGets()}, nil
}

func (p *ApiModule) GetStatsName(ctx context.Context, in *Empty) (*StatsName, error) {
	return &StatsName{CounterName: statsCounterName}, nil
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
	return nil

	if !p.disable {
		p.stop(service)
		time.Sleep(time.Second)
	}
	p.prestart(service)
	return p.start(service)
}
