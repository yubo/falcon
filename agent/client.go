/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
)

type ClientModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func clientPut(client service.ServiceClient, items []*service.Item, timeout int) {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(items))

	glog.V(6).Infof("%s tx put %v", MODULE_NAME, items)

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	resp, err := client.Put(ctx, &service.PutRequest{Items: items})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	} else {
		statsInc(ST_TX_PUT_ERR_ITEMS, int(len(items)-int(resp.N)))
	}
}

func (p *ClientModule) txWorker(ch chan *putContext, upstream, host_ string,
	callTimeout, burstSize int) error {

	host := []byte(host_)
	if upstream == "stdout" {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case put := <-ch:

					for _, item := range put.items {
						fmt.Printf("%s TX PUT %10.4f %s %s\n",
							MODULE_NAME, item.Value,
							item.Metric, item.Tags)
					}
					if put.done != nil {
						put.done <- &PutResponse{N: int32(len(put.items))}
					}
				}
			}
		}()
		return nil
	}

	go func() {
		conn, _, err := falcon.DialRr(p.ctx, upstream, true)
		if err != nil {
			return
		}
		defer conn.Close()
		client := service.NewServiceClient(conn)

		items := make([]*service.Item, burstSize)
		i := 0

		for {
			select {
			case <-p.ctx.Done():
				return
			case put := <-ch:
				glog.V(5).Infof("%s tx put %d\n", MODULE_NAME, len(put.items))
				n := 0
				for _, item_ := range put.items {

					// TODO check
					item, err := item_.toServiceItem(host)
					if err != nil {
						break
					}

					glog.V(6).Infof("%s TX PUT %d %10.4f %s\n",
						MODULE_NAME, item.Timestamp,
						item.Value, item.Key)

					items[i] = item
					i++
					n++

					if i == burstSize {
						clientPut(client, items[:i], callTimeout)
						i = 0
					}
				}
				if put.done != nil {
					put.done <- &PutResponse{N: int32(n)}
				}
			case <-time.After(time.Second):
				if i > 0 {
					clientPut(client, items[:i], callTimeout)
					i = 0
				}
			}

		}
	}()
	return nil
}

func (p *ClientModule) prestart(agent *Agent) error {
	return nil
}

func (p *ClientModule) start(agent *Agent) error {

	upstream := agent.Conf.Configer.Str(C_UPSTREAM)
	callTimeout, _ := agent.Conf.Configer.Int(C_CALL_TIMEOUT)
	burstSize, _ := agent.Conf.Configer.Int(C_BURST_SIZE)
	p.ctx, p.cancel = context.WithCancel(context.Background())

	if err := p.txWorker(agent.putChan, upstream, agent.Conf.Host,
		callTimeout, burstSize); err != nil {
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
	return p.start(agent)
}
