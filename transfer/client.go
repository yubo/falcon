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
	"github.com/yubo/falcon/service"
	"golang.org/x/net/context"
)

// ClientModule: transfer's module for banckend
// servicegroup: upstream container
// upstream: connection to the

type rpcClient struct {
	addr string
	conn *grpc.ClientConn
	cli  service.ServiceClient
}

type ClientModule struct {
	putChan      chan []*service.Item
	shardmap     []chan *service.Item
	serviceChans map[string]chan *service.Item
	clients      map[string]*rpcClient
	callTimeout  int
	burstSize    int
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *ClientModule) prestart(transfer *Transfer) error {

	p.putChan = transfer.appPutChan
	p.shardmap = make([]chan *service.Item, len(transfer.Conf.Upstream))
	p.serviceChans = make(map[string]chan *service.Item)
	p.clients = make(map[string]*rpcClient)
	p.callTimeout, _ = transfer.Conf.Configer.Int(C_CALL_TIMEOUT)
	p.burstSize, _ = transfer.Conf.Configer.Int(C_BURST_SIZE)

	for k, v := range transfer.Conf.Upstream {
		ch, ok := p.serviceChans[v]
		if !ok {
			ch = make(chan *service.Item)
			p.serviceChans[v] = ch
			p.clients[v] = &rpcClient{addr: v}
		}
		p.shardmap[k] = ch
	}

	return nil
}

func (p *ClientModule) start(transfer *Transfer) (err error) {

	glog.V(3).Infof(MODULE_NAME+"%s len(shardMap) %d", transfer.Conf.Name, len(p.shardmap))

	p.ctx, p.cancel = context.WithCancel(context.Background())

	for _, c := range p.clients {
		// FIXME: remove WithBlock, and reconnection when service online
		c.conn, _, err = falcon.DialRr(p.ctx, c.addr, true)
		if err != nil {
			glog.Errorf("%s addr:%s err:%s\n",
				MODULE_NAME, c.addr, err)
			return err
		}
		c.cli = service.NewServiceClient(c.conn)
	}

	go putWorker(p.ctx, p.putChan, p.shardmap)

	for k, v := range p.serviceChans {
		go clientWorker(p.ctx, v, p.clients[k],
			p.burstSize, p.callTimeout)
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

func putWorker(ctx context.Context, in chan []*service.Item, out []chan *service.Item) {
	n := uint64(len(out))
	for {
		select {
		case <-ctx.Done():
			return
		case items := <-in:
			for _, item := range items {
				item.ShardId = int32(item.Sum64() % n)
				select {
				case out[item.ShardId] <- item:
				default:
				}
			}
		}
	}

}

func clientPut(client service.ServiceClient, items []*service.Item,
	timeout int) {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(items))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	res, err := client.Put(ctx,
		&service.PutRequest{Items: items})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	}
	statsInc(ST_TX_PUT_ERR_ITEMS, int(len(items)-int(res.N)))
}

func clientWorker(ctx context.Context,
	ch chan *service.Item, c *rpcClient, burstSize, timeout int) {

	var i int

	items := make([]*service.Item, burstSize)
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		case item := <-ch:
			items[i] = item
			i++
			if i == burstSize {
				clientPut(c.cli, items[:i], timeout)
				i = 0
			}
		case <-time.After(time.Second):
			clientPut(c.cli, items[:i], timeout)
			i = 0
		}
	}
}
