/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

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
	//putChan chan []*service.Item
	reqChan chan *reqPayload
}

func (p *ApiModule) Get(ctx context.Context,
	in *service.GetRequest) (res *service.GetResponse, err error) {

	req := &reqPayload{
		action: RPC_ACTION_GET,
		data:   in,
		done:   make(chan interface{}),
	}
	p.reqChan <- req
	res = (<-req.done).(*service.GetResponse)

	// directly forward
	statsInc(ST_RX_GET_ITERS, 1)
	statsInc(ST_RX_GET_ITEMS, len(res.Dps))
	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *service.PutRequest) (res *service.PutResponse, err error) {

	glog.V(5).Infof("%s rx put %v", MODULE_NAME, len(in.Items))

	res = &service.PutResponse{}

	for _, item := range in.Items {
		if err = item.Adjust(); err != nil {
			return
		}

		p.reqChan <- &reqPayload{
			action: RPC_ACTION_PUT,
			data:   item,
		}
		res.N++
	}

	statsInc(ST_RX_PUT_ITERS, 1)
	statsInc(ST_RX_PUT_ITEMS, int(res.N))
	statsInc(ST_RX_PUT_ERR_ITEMS, len(in.Items)-int(res.N))

	return
}

func (p *ApiModule) GetStats(ctx context.Context, in *service.Empty) (*service.Stats, error) {
	return &service.Stats{Counter: statsGets()}, nil
}

func (p *ApiModule) GetStatsName(ctx context.Context, in *service.Empty) (*service.StatsName, error) {
	return &service.StatsName{CounterName: statsCounterName}, nil
}

func (p *ApiModule) prestart(transfer *Transfer) error {
	p.address = transfer.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.reqChan = transfer.reqChan

	return nil
}

func (p *ApiModule) start(transfer *Transfer) (err error) {

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

func (p *ApiModule) stop(transfer *Transfer) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(transfer *Transfer) error {
	return nil

	if !p.disable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
