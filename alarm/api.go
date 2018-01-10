/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

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
	putChan chan *Event
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (*PutResponse, error) {

	glog.V(3).Infof("%s RX PUT %s", MODULE_NAME, in)

	return &PutResponse{N: int32(len(in.Events))}, nil
}

func (p *ApiModule) prestart(a *Alarm) error {
	p.address = a.Conf.Configer.Str(C_API_ADDR)
	p.disable = falcon.AddrIsDisable(p.address)
	p.putChan = a.appPutChan

	return nil
}

func (p *ApiModule) start(a *Alarm) (err error) {

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
	RegisterAlarmServer(server, p)

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

func (p *ApiModule) stop(a *Alarm) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ApiModule) reload(a *Alarm) error {
	if !p.disable {
		p.stop(a)
		time.Sleep(time.Second)
	}
	p.prestart(a)
	return p.start(a)
}
