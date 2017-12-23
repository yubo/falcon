/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"time"

	"google.golang.org/grpc"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend"
	"golang.org/x/net/context"
)

type rpcClients struct {
	addr string
	conn *grpc.ClientConn
	cli  backend.BackendClient
}

type senderFalcon struct {
	name            string
	workerProcesses int
	connTimeout     int
	callTimeout     int
	payloadSize     int
	clients         map[string]rpcClients
	chans           map[string]chan *falcon.Item
	ctx             context.Context
	cancel          context.CancelFunc
}

func (p *senderFalcon) new(L *Transfer) sender {
	workerprocesses, _ := L.Conf.Configer.Int(C_WORKER_PROCESSES)
	conntimeout, _ := L.Conf.Configer.Int(C_CONN_TIMEOUT)
	calltimeout, _ := L.Conf.Configer.Int(C_CALL_TIMEOUT)
	payloadsize, _ := L.Conf.Configer.Int(C_PAYLOADSIZE)

	return &senderFalcon{
		name:            "falcon",
		workerProcesses: workerprocesses,
		connTimeout:     conntimeout,
		callTimeout:     calltimeout,
		payloadSize:     payloadsize,
		clients:         make(map[string]rpcClients),
		chans:           make(map[string]chan *falcon.Item),
	}
}

func (p *senderFalcon) addClientChan(key, addr string,
	ch chan *falcon.Item) (err error) {
	p.chans[key] = ch
	p.clients[key] = rpcClients{addr: addr}
	return nil
}

func (p *senderFalcon) start(name string) (err error) {
	p.ctx, p.cancel = context.WithCancel(context.Background())

	for k, c := range p.clients {
		// FIXME: remove WithBlock, and reconnection when backend online
		c.conn, _, err = falcon.DialRr(p.ctx, c.addr, true)
		if err != nil {
			glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n",
				k, c.addr, err)
		}
		c.cli = backend.NewBackendClient(c.conn)
	}

	for node, ch := range p.chans {
		go falconUpstreamWorker(p.ctx, name,
			ch, p.clients[node],
			p.clients[node].addr,
			p.payloadSize)
	}
	return nil
}

func (p *senderFalcon) stop() error {
	p.cancel()
	return nil
}

// falcon {{{
func falconUpstreamUpdate(client backend.BackendClient, items []*falcon.Item) {
	if res, err := client.Put(context.Background(),
		&falcon.PutRequest{Items: items}); err != nil {
		glog.Error(err)
		statsInc(ST_UPSTREAM_PUT_ERR, 1)
	} else {
		statsInc(ST_UPSTREAM_PUT, 1)
		statsInc(ST_UPSTREAM_PUT_ITEM_TOTAL, int(res.Total))
		statsInc(ST_UPSTREAM_PUT_ITEM_ERR, int(res.Errors))
	}
}

func falconUpstreamWorker(ctx context.Context, name string,
	ch chan *falcon.Item, c rpcClients, addr string, payloadSize int) {

	var i int

	items := make([]*falcon.Item, payloadSize)
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		case item := <-ch:
			items[i] = item
			i++
			if i == payloadSize {
				falconUpstreamUpdate(c.cli, items[:i])
				i = 0
			}
		case <-time.After(time.Second):
			falconUpstreamUpdate(c.cli, items[:i])
			i = 0
		}
	}
}

//}}}
