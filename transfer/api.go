/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"net"

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
	shardmap []chan *reqPayload
}

func (p *ApiModule) Get(ctx context.Context,
	in *GetRequest) (res *GetResponse, err error) {

	var rs []*service.GetResponse

	reqs, err := in.Adjust()
	if err != nil {
		return nil, err
	}

	for shardId, req := range reqs {
		req := &reqPayload{
			action: RPC_ACTION_GET,
			data:   req,
			done:   make(chan interface{}),
		}
		p.shardmap[shardId] <- req
		rs = append(rs, (<-req.done).(*service.GetResponse))
	}

	res = &GetResponse{}
	for _, r := range rs {
		for _, v := range r.Data {
			res.Data = append(res.Data, &DataPoints{
				Key:    v.Key.Key,
				Values: v.Values,
			})
		}
	}

	// directly forward
	statsInc(ST_RX_GET_ITERS, 1)
	statsInc(ST_RX_GET_ITEMS, len(res.Data))
	return
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (res *PutResponse, err error) {

	glog.V(5).Infof("%s rx put %v", MODULE_NAME, len(in.Data))

	res = &PutResponse{}

	for _, dp := range in.Data {
		tdp, err := dp.Adjust()
		if err != nil {
			return nil, err
		}

		req := &reqPayload{
			action: RPC_ACTION_PUT,
			data:   tdp,
		}
		p.shardmap[tdp.Key.ShardId] <- req
		res.N++
	}

	statsInc(ST_RX_PUT_ITERS, 1)
	statsInc(ST_RX_PUT_ITEMS, int(res.N))
	statsInc(ST_RX_PUT_ERR_ITEMS, len(in.Data)-int(res.N))

	return
}

func (p *ApiModule) GetStats(ctx context.Context, in *Empty) (*Stats, error) {
	return &Stats{Counter: statsGets()}, nil
}

func (p *ApiModule) GetStatsName(ctx context.Context, in *Empty) (*StatsName, error) {
	return &StatsName{CounterName: statsCounterName}, nil
}

func (p *ApiModule) prestart(transfer *Transfer) error {
	p.shardmap = transfer.shardmap
	return nil
}

func (p *ApiModule) start(transfer *Transfer) (err error) {

	if p.disable {
		glog.Info(MODULE_NAME + "api disable")
		return nil
	}

	p.address = transfer.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterTransferServer(server, p)

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
}
