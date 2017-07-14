/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"container/list"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type connList struct {
	sync.RWMutex
	list *list.List
}

func (l *connList) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}
func (l *connList) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

type LB struct {
	appUpdateChan chan *[]*falcon.MetaData // upstreams
}

func (p *LB) Ping(req falcon.Null, resp *falcon.RpcResp) error {
	return nil
}

func (p *LB) Update(args []*falcon.MetaData,
	reply *falcon.LbResp) error {
	reply.Invalid = 0
	now := time.Now().Unix()

	items := []*falcon.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.Name == "" || v.Host == "" {
			reply.Invalid += 1
			continue
		}

		if v.Type != falcon.COUNTER &&
			v.Type != falcon.GAUGE &&
			v.Type != falcon.DERIVE {
			reply.Invalid += 1
			continue
		}

		if v.Step <= 0 {
			reply.Invalid += 1
			continue
		}

		if len(v.Name)+len(v.Tags) > 510 {
			reply.Invalid += 1
			continue
		}

		if v.Ts <= 0 || v.Ts > now*2 {
			v.Ts = now
		}

		items = append(items, &falcon.MetaData{
			Name: v.Name,
			Host: v.Host,
			Ts:   v.Ts,
			Step: v.Step,
			Type: v.Type,
			Tags: sortTags(v.Tags),
		})
	}

	p.appUpdateChan <- &items
	glog.V(3).Infof(MODULE_NAME+"recv %d", len(items))

	reply.Message = "ok"
	reply.Total = len(args)

	statsInc(ST_RPC_UPDATE, 1)
	statsInc(ST_RPC_UPDATE_CNT, len(items))
	statsInc(ST_RPC_UPDATE_ERR, reply.Invalid)

	return nil
}

type RpcModule struct {
	rpcConnects connList
	rpcListener net.Listener
}

func (p *RpcModule) prestart(L *Transfer) error {
	p.rpcConnects = connList{list: list.New()}
	return nil
}

func (p *RpcModule) start(L *Transfer) (err error) {

	enable, _ := L.Conf.Configer.Bool(C_RPC_ENABLE)
	if !enable {
		glog.Info(MODULE_NAME + "rpc.Start warning, not enabled")
		return nil
	}

	rpc.Register(&LB{appUpdateChan: L.appUpdateChan})

	addr := L.Conf.Configer.Str(C_RPC_ADDR)
	p.rpcListener, err = net.Listen(falcon.ParseAddr(addr))
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, listen %s failed, %s",
			addr, err)
	} else {
		glog.Infof(MODULE_NAME+"%s rpcStart ok, listening on %s", L.Conf.Name, addr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					time.Sleep(tempDelay)
					continue
				}
				return
			}
			tempDelay = 0
			go func() {
				e := p.rpcConnects.insert(conn)
				defer p.rpcConnects.remove(e)
				jsonrpc.ServeConn(conn)
			}()
		}
	}()
	return err
}

func (p *RpcModule) stop(L *Transfer) error {
	if p.rpcListener == nil {
		return falcon.ErrNoExits
	}

	p.rpcListener.Close()
	p.rpcConnects.Lock()
	for e := p.rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	p.rpcConnects.Unlock()
	p.rpcConnects = connList{list: list.New()}

	return nil
}

func (p *RpcModule) reload(L *Transfer) error {
	// TODO
	return nil
}
