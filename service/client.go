/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"time"

	"google.golang.org/grpc"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm"
	"golang.org/x/net/context"
)

// ClientModule: service's module for banckend
// servicegroup: upstream container
// upstream: connection to the

const (
	CLIENT_BURST_SIZE = 32
)

type rpcClient struct {
	addr string
	conn *grpc.ClientConn
	cli  alarm.AlarmClient
}

type ClientModule struct {
	conn   *grpc.ClientConn
	client alarm.AlarmClient
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ClientModule) put(events []*alarm.Event, timeout int) {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(events))

	glog.V(5).Infof("%s tx put %v", MODULE_NAME, events)
	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	res, err := p.client.Put(ctx, &alarm.PutRequest{Events: events})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	}
	statsInc(ST_TX_PUT_ERR_ITEMS, int(len(events)-int(res.N)))
}

func (p *ClientModule) mainLoop(upstreamAddr string, eventChan chan *alarm.Event,
	callTimeout, burstSize int) error {
	if upstreamAddr == "stdout" {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case e := <-eventChan:
					glog.Infof("%s TX Event %s\n", MODULE_NAME, e)
				}
			}
		}()
		return nil
	}

	go func() {
		conn, _, err := falcon.DialRr(p.ctx, upstreamAddr, true)
		if err != nil {
			glog.Fatalf("%s %v", MODULE_NAME, err)
			return
		}
		defer conn.Close()

		p.client = alarm.NewAlarmClient(conn)
		events := make([]*alarm.Event, burstSize)
		i := 0

		for {
			select {
			case <-p.ctx.Done():
				return
			case e := <-eventChan:
				glog.V(4).Infof("%s tx put %s\n", MODULE_NAME, e)

				events[i] = e
				i++

				if i == burstSize {
					p.put(events[:i], callTimeout)
					i = 0
				}
			case <-time.After(time.Second):
				if i > 0 {
					p.put(events[:i], callTimeout)
					i = 0
				}
			}

		}
	}()
	return nil
}

func (p *ClientModule) prestart(s *Service) error {
	return nil
}

func (p *ClientModule) start(s *Service) (err error) {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	upstreamAddr := s.Conf.Configer.Str(C_ALARM_ADDR)
	callTimeout, _ := s.Conf.Configer.Int(C_CALL_TIMEOUT)

	if err := p.mainLoop(upstreamAddr, s.eventChan, callTimeout,
		CLIENT_BURST_SIZE); err != nil {
		return err
	}
	return nil
}

func (p *ClientModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(s *Service) error {
	return nil
}
