/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type UpstreamModule struct {
	conn   *grpc.ClientConn
	client transfer.TransferClient
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *UpstreamModule) update(items []*falcon.Item, timeout int) (err error) {

	statsInc(ST_UPSTREAM_UPDATE, 1)
	statsInc(ST_UPSTREAM_UPDATE_ITEM, len(items))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	_, err = p.client.Update(ctx, &falcon.UpdateRequest{Items: items})
	if err != nil {
		statsInc(ST_UPSTREAM_UPDATE_ERR, 1)
	}
	return err
}

func (p *UpstreamModule) mainLoop(agent *Agent) error {
	callTimeout, _ := agent.Conf.Configer.Int(C_CALL_TIMEOUT)
	payloadSize, _ := agent.Conf.Configer.Int(C_PAYLOADSIZE)

	conn, _, err := falcon.DialRr(p.ctx, agent.Conf.Configer.Str(C_UPSTREAM), true)
	if err != nil {
		return err
	}

	p.client = transfer.NewTransferClient(conn)
	ch := agent.appUpdateChan
	items := make([]*falcon.Item, payloadSize)
	i := 0

	go func() {
		defer conn.Close()
		for {
			select {
			case <-p.ctx.Done():
				return
			case _items := <-ch:
				for _, item := range _items {
					items[i] = item
					i++
					if i == payloadSize {
						p.update(items[:i], callTimeout)
						i = 0
					}

				}
			case <-time.After(time.Second):
				p.update(items[:i], callTimeout)
				i = 0

			}

		}
	}()
	return nil

}

func (p *UpstreamModule) debugLoop(agent *Agent) error {

	if agent.Conf.Debug > 1 {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case items := <-agent.appUpdateChan:
					for k, v := range items {
						glog.V(3).Infof(MODULE_NAME+"%d %s", k, v)
					}
				}
			}
		}()
	}
	return nil
}

func (p *UpstreamModule) prestart(agent *Agent) error {
	rand.Seed(time.Now().Unix())
	return nil
}

func (p *UpstreamModule) start(agent *Agent) error {

	p.ctx, p.cancel = context.WithCancel(context.Background())

	if err := p.debugLoop(agent); err != nil {
		return err
	}

	if err := p.mainLoop(agent); err != nil {
		return err
	}

	return nil
}

func (p *UpstreamModule) stop(agent *Agent) error {
	p.cancel()
	return nil
}

func (p *UpstreamModule) reload(agent *Agent) error {
	p.stop(agent)
	time.Sleep(time.Second)
	p.prestart(agent)
	return p.start(agent)
}
