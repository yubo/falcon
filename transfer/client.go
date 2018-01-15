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

const (
	RPC_ACTION_PUT = iota
	RPC_ACTION_GET
)

type ClientModule struct {
	reqChan      chan *reqPayload
	shardmap     []chan *reqPayload
	serviceChans map[string]chan *reqPayload
	clients      map[string]*rpcClient
	callTimeout  int
	burstSize    int
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *ClientModule) prestart(transfer *Transfer) error {

	p.reqChan = transfer.reqChan
	p.shardmap = make([]chan *reqPayload, len(transfer.Conf.Upstream))
	p.serviceChans = make(map[string]chan *reqPayload)
	p.clients = make(map[string]*rpcClient)
	p.callTimeout, _ = transfer.Conf.Configer.Int(C_CALL_TIMEOUT)
	p.burstSize, _ = transfer.Conf.Configer.Int(C_BURST_SIZE)

	for shardId, addr := range transfer.Conf.Upstream {
		ch := p.serviceChans[addr]
		if ch == nil {
			ch = make(chan *reqPayload, 144)
			p.serviceChans[addr] = ch
			p.clients[addr] = &rpcClient{addr: addr}
		}
		p.shardmap[shardId] = ch
	}

	return nil
}

func (p *ClientModule) start(transfer *Transfer) (err error) {

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

	go putWorker(p.ctx, p.reqChan, p.shardmap)

	for addr, ch := range p.serviceChans {
		go clientWorker(p.ctx, ch, p.clients[addr],
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

func putWorker(ctx context.Context, in chan *reqPayload, out []chan *reqPayload) {
	n := uint64(len(out))
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-in:
			switch req.action {
			case RPC_ACTION_PUT:
				item := req.data.(*service.Item)
				item.ShardId = int32(item.Sum64() % n)
				out[item.ShardId] <- req
			case RPC_ACTION_GET:
				r := req.data.(*service.GetRequest)
				r.ShardId = int32(falcon.Sum64(r.Key) % n)
				out[r.ShardId] <- req
			}
		}
	}
}

func clientPut(client service.ServiceClient, items []*service.Item,
	timeout int) *service.PutResponse {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(items))

	glog.V(5).Infof("%s tx put %v", MODULE_NAME, len(items))
	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	res, err := client.Put(ctx,
		&service.PutRequest{Items: items})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	}
	statsInc(ST_TX_PUT_ERR_ITEMS, int(len(items)-int(res.N)))
	return res
}

func clientGet(client service.ServiceClient, req *service.GetRequest,
	timeout int) *service.GetResponse {

	statsInc(ST_TX_GET_ITERS, 1)

	glog.V(5).Infof("%s tx get %v", MODULE_NAME, req)
	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	res, err := client.Get(ctx, req)
	if err != nil {
		statsInc(ST_TX_GET_ERR_ITERS, 1)
	}
	statsInc(ST_TX_GET_ITEMS, int(len(res.Dps)))
	return res
}

func clientWorker(ctx context.Context,
	ch chan *reqPayload, c *rpcClient, burstSize, timeout int) {

	var i int

	items := make([]*service.Item, burstSize)
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		case req := <-ch:
			switch req.action {
			case RPC_ACTION_PUT:
				if req.done != nil {
					req.done <- clientPut(c.cli, []*service.Item{
						req.data.(*service.Item)}, timeout)
					continue
				}
				items[i] = req.data.(*service.Item)
				i++
				if i == burstSize {
					clientPut(c.cli, items[:i], timeout)
					i = 0
				}
			case RPC_ACTION_GET:
				res := clientGet(c.cli, req.data.(*service.GetRequest), timeout)
				if req.done != nil {
					req.done <- res
				}
			}
		case <-time.After(time.Second):
			if i > 0 {
				clientPut(c.cli, items[:i], timeout)
				i = 0
			}
		}
	}
}
