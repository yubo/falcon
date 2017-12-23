/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ClientModule struct {
	conn   *grpc.ClientConn
	client service.ServiceClient
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ClientModule) update(items []*falcon.Item, timeout int) (err error) {

	statsInc(ST_UPSTREAM_UPDATE, 1)
	statsInc(ST_UPSTREAM_UPDATE_ITEM, len(items))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	_, err = p.client.Put(ctx, &falcon.PutRequest{Items: items})
	if err != nil {
		statsInc(ST_UPSTREAM_UPDATE_ERR, 1)
	}
	return err
}

func stdoutLoop(p *ClientModule, ch chan []*falcon.Item) {
	for {
		select {
		case <-p.ctx.Done():
			return
		case get_items := <-ch:
			for k, v := range get_items {
				fmt.Printf(MODULE_NAME+"%d %s\n", k, v)
			}
		}
	}
}

func socketLoop(p *ClientModule, ch chan []*falcon.Item, agent *Agent) {
	callTimeout, _ := agent.Conf.Configer.Int(C_CALL_TIMEOUT)
	burstSize, _ := agent.Conf.Configer.Int(C_BURST_SIZE)
	items := make([]*falcon.Item, burstSize)
	i := 0

	for {
		select {
		case <-p.ctx.Done():
			return
		case get_items := <-ch:
			for _, item := range get_items {
				items[i] = item
				i++
				if i == burstSize {
					p.update(items[:i], callTimeout)
					i = 0
				}

			}
		case <-time.After(time.Second):
			p.update(items[:i], callTimeout)
			i = 0

		}

	}

}

func (p *ClientModule) mainLoop(agent *Agent) error {

	upstream := agent.Conf.Configer.Str(C_UPSTREAM)
	if upstream == "stdout" {
		go stdoutLoop(p, agent.appUpdateChan)
		return nil
	}

	conn, _, err := falcon.DialRr(p.ctx, upstream, true)
	if err != nil {
		return err
	}
	defer conn.Close()

	p.client = service.NewServiceClient(conn)
	go socketLoop(p, agent.appUpdateChan, agent)

	return nil
}

func (p *ClientModule) prestart(agent *Agent) error {
	rand.Seed(time.Now().Unix())
	return nil
}

func (p *ClientModule) start(agent *Agent) error {

	p.ctx, p.cancel = context.WithCancel(context.Background())

	if err := p.mainLoop(agent); err != nil {
		return err
	}

	return nil
}

func (p *ClientModule) stop(agent *Agent) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(agent *Agent) error {
	p.stop(agent)
	time.Sleep(time.Second)
	p.prestart(agent)
	return p.start(agent)
}
