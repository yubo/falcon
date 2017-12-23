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
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServiceModule struct {
	enable     bool
	ctx        context.Context
	cancel     context.CancelFunc
	address    string
	updateChan chan []*falcon.Item
}

func (p *ServiceModule) Get(ctx context.Context,
	in *falcon.GetRequest) (*falcon.GetResponse, error) {
	return &falcon.GetResponse{}, nil
}

func (p *ServiceModule) Put(ctx context.Context,
	in *falcon.PutRequest) (*falcon.PutResponse, error) {

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

	p.updateChan <- in.Items
	glog.V(3).Infof(MODULE_NAME+"recv %d", len(items))

	statsInc(ST_RPC_UPDATE, 1)
	statsInc(ST_RPC_UPDATE_CNT, len(items))
	statsInc(ST_RPC_UPDATE_ERR, int(res.Errors))

	return res, nil
}

func (p *ServiceModule) prestart(transfer *Transfer) error {
	p.enable, _ = transfer.Conf.Configer.Bool(C_GRPC_ENABLE)
	p.address = transfer.Conf.Configer.Str(C_GRPC_ADDR)
	p.updateChan = transfer.appUpdateChan
	return nil
}

func (p *ServiceModule) start(transfer *Transfer) (err error) {

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
	RegisterTransferServer(server, &ServiceModule{})

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

func (p *ServiceModule) stop(transfer *Transfer) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ServiceModule) reload(transfer *Transfer) error {
	if p.enable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}
