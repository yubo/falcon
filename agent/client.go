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

	"github.com/golang/glog"
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

func (p *ClientModule) put(items []*falcon.Item, timeout int) {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(items))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	resp, err := p.client.Put(ctx, &falcon.PutRequest{Items: items})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	}
	statsInc(ST_TX_PUT_ERR_ITEMS, int(resp.Errors))
}

func (client *ClientModule) mainLoop(agent *Agent) error {
	upstream := agent.Conf.Configer.Str(C_UPSTREAM)
	if upstream == "stdout" {
		go func() {
			for {
				select {
				case <-client.ctx.Done():
					return
				case get_items := <-agent.appPutChan:
					for _, item := range get_items {
						fmt.Printf("%s TX PUT %d %15s %10.4f %s\n",
							MODULE_NAME, item.Timestamp, item.Metric,
							item.Value, item.Tags)
					}
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

		client.client = service.NewServiceClient(conn)
		callTimeout, _ := agent.Conf.Configer.Int(C_CALL_TIMEOUT)
		burstSize, _ := agent.Conf.Configer.Int(C_BURST_SIZE)
		host := []byte(agent.Conf.Host)
		items := make([]*falcon.Item, burstSize)
		i := 0

		for {
			select {
			case <-client.ctx.Done():
				return
			case get_items := <-agent.appPutChan:
				glog.V(4).Infof("%s TX PUT %d\n", MODULE_NAME, len(get_items))
				for _, item := range get_items {

					glog.V(5).Infof("%s TX PUT %d %15s %10.4f %s\n",
						MODULE_NAME, item.Timestamp, item.Metric,
						item.Value, item.Tags)

					item.Endpoint = host
					items[i] = item
					i++
					if i == burstSize {
						client.put(items[:i], callTimeout)
						i = 0
					}

				}
			case <-time.After(time.Second):
				if i > 0 {
					client.put(items[:i], callTimeout)
					i = 0
				}
			}

		}
	}()
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
