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

// ClientModule: transfer's module for banckend
// backendgroup: upstream container
// upstream: connection to the

type sender interface {
	new(*Transfer) sender
	start(string) error
	stop() error
	addClientChan(string, string, chan *falcon.Item) error
}

func init() {
}

type rpcClients struct {
	addr string
	conn *grpc.ClientConn
	cli  backend.BackendClient
}

type upstream struct {
	name     string
	upstream sender
}

/* upstream */
type ClientModule struct {
	sharemap     []chan *falcon.Item
	serviceChans map[string]chan *falcon.Item
	clients      map[string]*rpcClients
	connTimeout  int
	callTimeout  int
	burstSize    int
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *ClientModule) prestart(transfer *Transfer) error {

	p.sharemap = make([]chan *falcon.Item, transfer.Conf.ShareCount)
	p.serviceChans = make(map[string]chan *falcon.Item)
	p.clients = make(map[string]*rpcClients)
	p.connTimeout, _ = transfer.Conf.Configer.Int(C_CONN_TIMEOUT)
	p.callTimeout, _ = transfer.Conf.Configer.Int(C_CALL_TIMEOUT)
	p.burstSize, _ = transfer.Conf.Configer.Int(C_BURST_SIZE)

	for shareid, service := range transfer.Conf.ShareMap {
		ch, ok := p.serviceChans[service]
		if !ok {
			ch = make(chan *falcon.Item)
			p.serviceChans[service] = ch
			p.clients[service] = &rpcClients{addr: service}
		}
		p.sharemap[shareid] = ch
	}

	return nil
}

func (p *ClientModule) start(transfer *Transfer) (err error) {

	glog.V(3).Infof(MODULE_NAME+"%s len(shareMap) %d", transfer.Conf.Name, len(p.sharemap))

	p.ctx, p.cancel = context.WithCancel(context.Background())

	for _, c := range p.clients {
		// FIXME: remove WithBlock, and reconnection when backend online
		c.conn, _, err = falcon.DialRr(p.ctx, c.addr, true)
		if err != nil {
			glog.Fatalf(MODULE_NAME+"addr:%s err:%s\n",
				c.addr, err)
		}
		c.cli = backend.NewBackendClient(c.conn)
	}

	for service, ch := range p.serviceChans {
		go clientWorker(p.ctx, ch, p.clients[service],
			p.burstSize)
	}

	return nil
}

func (p *ClientModule) stop(transfer *Transfer) error {
	p.cancel()
	return nil
}

func (p *ClientModule) reload(transfer *Transfer) error {
	p.stop(transfer)
	time.Sleep(time.Second)
	p.prestart(transfer)
	return p.start(transfer)
}

func clientPut(client backend.BackendClient, items []*falcon.Item) {
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

func clientWorker(ctx context.Context,
	ch chan *falcon.Item, c *rpcClients, burstSize int) {

	var i int

	items := make([]*falcon.Item, burstSize)
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		case item := <-ch:
			items[i] = item
			i++
			if i == burstSize {
				clientPut(c.cli, items[:i])
				i = 0
			}
		case <-time.After(time.Second):
			clientPut(c.cli, items[:i])
			i = 0
		}
	}
}
