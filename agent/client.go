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
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
)

type ClientModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func clientPut(client transfer.TransferClient, dps []*transfer.DataPoint, timeout int) {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(dps))

	glog.V(6).Infof("%s tx put %v", MODULE_NAME, dps)

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	resp, err := client.Put(ctx, &transfer.PutRequest{Data: dps})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	} else {
		statsInc(ST_TX_PUT_ERR_ITEMS, int(len(dps)-int(resp.N)))
	}
}

func (p *ClientModule) txWorker(ch chan *putContext, upstream string,
	callTimeout, burstSize int) error {

	if upstream == "stdout" {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case put := <-ch:

					for _, dp := range put.dps {
						fmt.Printf("%s TX PUT %10.4f %s\n",
							MODULE_NAME, dp.Value.Value,
							string(dp.Key))
					}
					if put.done != nil {
						put.done <- &transfer.PutResponse{N: int32(len(put.dps))}
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
		client := transfer.NewTransferClient(conn)

		dps := make([]*transfer.DataPoint, burstSize)
		i := 0

		for {
			select {
			case <-p.ctx.Done():
				return
			case put := <-ch:
				glog.V(5).Infof("%s tx put %d\n", MODULE_NAME, len(put.dps))
				n := 0
				for _, dp := range put.dps {

					// TODO check
					glog.V(6).Infof("%s TX PUT %d %10.4f %s\n",
						MODULE_NAME, dp.Value.Timestamp,
						dp.Value.Value, dp.Key)

					dps[i] = dp
					i++
					n++

					if i == burstSize {
						clientPut(client, dps[:i], callTimeout)
						i = 0
					}
				}
				if put.done != nil {
					put.done <- &transfer.PutResponse{N: int32(n)}
				}
			case <-time.After(time.Second):
				if i > 0 {
					clientPut(client, dps[:i], callTimeout)
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
	conf := &agent.Conf.Configer

	upstream := conf.Str(C_UPSTREAM)
	callTimeout, _ := conf.Int(C_CALL_TIMEOUT)
	burstSize, _ := conf.Int(C_BURST_SIZE)
	p.ctx, p.cancel = context.WithCancel(context.Background())

	if err := p.txWorker(agent.putChan, upstream,
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
