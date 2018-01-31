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
	"github.com/yubo/falcon/lib/tsdb"
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
	shardmap     []chan *reqPayload
	serviceChans map[string]chan *reqPayload
	clients      map[string]*rpcClient
	callTimeout  int
	burstSize    int
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *ClientModule) prestart(transfer *Transfer) error {
	p.shardmap = transfer.shardmap
	return nil
}

func (p *ClientModule) start(transfer *Transfer) (err error) {
	conf := &transfer.Conf.Configer

	p.serviceChans = make(map[string]chan *reqPayload)
	p.clients = make(map[string]*rpcClient)
	p.callTimeout, _ = conf.Int(C_CALL_TIMEOUT)
	p.burstSize, _ = conf.Int(C_BURST_SIZE)

	for shardId, addr := range transfer.Conf.ShardMap {
		ch := p.serviceChans[addr]
		if ch == nil {
			ch = make(chan *reqPayload, 144)
			p.serviceChans[addr] = ch
			p.clients[addr] = &rpcClient{addr: addr}
		}
		p.shardmap[shardId] = ch
	}

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

func clientPut(client service.ServiceClient, dps []*tsdb.DataPoint,
	timeout int) *service.PutResponse {

	statsInc(ST_TX_PUT_ITERS, 1)
	statsInc(ST_TX_PUT_ITEMS, len(dps))

	glog.V(5).Infof("%s tx put %v", MODULE_NAME, len(dps))
	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Millisecond)
	res, err := client.Put(ctx,
		&service.PutRequest{Data: dps})
	if err != nil {
		statsInc(ST_TX_PUT_ERR_ITERS, 1)
	}
	if res == nil {
		statsInc(ST_TX_PUT_ERR_ITEMS, int(len(dps)))
	} else {
		statsInc(ST_TX_PUT_ERR_ITEMS, int(len(dps)-int(res.N)))
	}
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
	statsInc(ST_TX_GET_ITEMS, int(len(res.Data)))
	return res
}

func clientWorker(ctx context.Context,
	ch chan *reqPayload, c *rpcClient, burstSize, timeout int) {

	var i int

	dps := make([]*tsdb.DataPoint, burstSize)
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		case req := <-ch:
			switch req.action {
			case RPC_ACTION_PUT:
				if req.done != nil {
					go func() {
						req.done <- clientPut(c.cli, []*tsdb.DataPoint{
							req.data.(*tsdb.DataPoint)}, timeout)
					}()
					continue
				}
				dps[i] = req.data.(*tsdb.DataPoint)
				i++
				if i == burstSize {
					clientPut(c.cli, dps[:i], timeout)
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
				clientPut(c.cli, dps[:i], timeout)
				i = 0
			}
		}
	}
}
