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
	"github.com/yubo/falcon/lib/core"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiModule struct {
	disable      bool
	ctx          context.Context
	cancel       context.CancelFunc
	address      string
	putEventChan chan *Event
}

func (p *ApiModule) Put(ctx context.Context,
	in *PutRequest) (*PutResponse, error) {

	res := &PutResponse{}

	glog.V(3).Infof("%s rx put %s", MODULE_NAME, in)
	for _, e := range in.Events {
		select {
		case p.putEventChan <- e:
			res.N++
		default:
		}
	}

	statsInc(ST_RX_PUT_ITER, 1)
	statsInc(ST_RX_PUT_EVENT, int(res.N))
	statsInc(ST_RX_PUT_ERR_EVENT, len(in.Events)-int(res.N))

	return res, nil
}

func (p *ApiModule) GetStats(ctx context.Context, in *Empty) (*Stats, error) {
	return &Stats{Counter: statsGets()}, nil
}

func (p *ApiModule) GetStatsName(ctx context.Context, in *Empty) (*StatsName, error) {
	return &StatsName{CounterName: statsCounterName}, nil
}

func (p *ApiModule) prestart(a *Alarm) error {
	p.address = a.Conf.ApiAddr
	p.disable = core.AddrIsDisable(p.address)
	p.putEventChan = a.putEventChan

	return nil
}

func (p *ApiModule) start(a *Alarm) (err error) {

	if p.disable {
		glog.Infof("%s api disable", MODULE_NAME)
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(core.CleanSockFile(core.ParseAddr(p.address)))
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
	return nil

	if !p.disable {
		p.stop(a)
		time.Sleep(time.Second)
	}
	p.prestart(a)
	return p.start(a)
}
