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
	putChan chan []*falcon.Item
}

func (p *ApiModule) Get(ctx context.Context,
	in *falcon.GetRequest) (*falcon.GetResponse, error) {
	return &falcon.GetResponse{}, nil
}

func (p *ApiModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

	glog.V(3).Infof("%s RX PUT %d", MODULE_NAME, len(in.Items))

	res := &falcon.PutResponse{Total: int32(len(in.Items))}
	items := []*falcon.Item{}
	now := time.Now().Unix()

	for _, item := range in.Items {
		if err := item.Adjust(now); err != nil {
			res.Errors += 1
			continue
		}
		items = append(items, item)
	}

	select {
	case p.putChan <- items:
	default:
		glog.V(4).Infof("%s RX PUT %d FAIL", MODULE_NAME, len(in.Items))
		res.Errors = res.Total
	}

	statsInc(ST_RX_PUT_ITERS, 1)
	statsInc(ST_RX_PUT_ITEMS, int(res.Total))
	statsInc(ST_RX_PUT_ERR_ITEMS, int(res.Errors))

	return res, nil
}

func (p *ApiModule) prestart(transfer *Transfer) error {
	p.address = transfer.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.putChan = transfer.appPutChan

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
	if !p.disable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
