/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
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

	//statsInc(ST_TX_PUT_EVENTS, 1)
	//statsInc(ST_TX_PUT_EVENTS, len(events))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	p.client.Put(ctx, &alarm.PutRequest{Events: events})
	//resp, err := p.client.Put(ctx, &alarm.PutRequest{Events: events})
	//if err != nil {
	//	statsInc(ST_TX_PUT_ERR_ITERS, 1)
	//}
	//statsInc(ST_TX_PUT_ERR_EVENTS, int(len(events)-int(resp.N)))
}

func (client *ClientModule) mainLoop(service *Service) error {
	upstream := service.Conf.Configer.Str(C_UPSTREAM)
	eventCh := service.appEventChan
	if upstream == "stdout" {
		go func() {
			for {
				select {
				case <-client.ctx.Done():
					return
				case e := <-eventCh:
					fmt.Printf("%s TX Event %s\n", MODULE_NAME, e)
				}
			}
		}()
		return nil
	}

	go func() {
		conn, _, err := falcon.DialRr(client.ctx, upstream, true)
		if err != nil {
			return
		}
		defer conn.Close()

		client.client = alarm.NewAlarmClient(conn)
		callTimeout, _ := service.Conf.Configer.Int(C_CALL_TIMEOUT)
		burstSize, _ := service.Conf.Configer.Int(C_BURST_SIZE)
		events := make([]*alarm.Event, burstSize)
		i := 0

		for {
			select {
			case <-client.ctx.Done():
				return
			case e := <-eventCh:
				glog.V(4).Infof("%s TX PUT %s\n", MODULE_NAME, e)

				events[i] = e
				i++

				if i == burstSize {
					client.put(events[:i], callTimeout)
					i = 0
				}
			case <-time.After(time.Second):
				if i > 0 {
					client.put(events[:i], callTimeout)
					i = 0
				}
			}

		}
	}()
	return nil
}

func (p *ClientModule) prestart(service *Service) error {
	return nil
}

func (p *ClientModule) start(service *Service) (err error) {

	p.ctx, p.cancel = context.WithCancel(context.Background())

	if err := p.mainLoop(service); err != nil {
		return err
	}
	return nil
}

func (p *ClientModule) stop(service *Service) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(service *Service) error {
	return nil
}
