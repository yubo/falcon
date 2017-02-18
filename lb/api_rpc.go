/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

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

var ()

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
	lb *Lb
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

	p.lb.appUpdateChan <- &items
	glog.V(3).Infof(MODULE_NAME+"recv %d", len(items))

	reply.Message = "ok"
	reply.Total = len(args)

	statInc(ST_RPC_UPDATE, 1)
	statInc(ST_RPC_UPDATE_CNT, len(items))
	statInc(ST_RPC_UPDATE_ERR, reply.Invalid)

	return nil
}

func (p *Lb) rpcStart() (err error) {

	var addr *net.TCPAddr
	if !p.Conf.Params.Rpc {
		return nil
	}
	lb := &LB{lb: p}
	rpc.Register(lb)

	addr, err = net.ResolveTCPAddr("tcp", p.Conf.Params.RpcAddr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	p.rpcListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, listen %s failed, %s",
			p.Conf.Params.RpcAddr, err)
	} else {
		glog.Infof(MODULE_NAME+"%s rpcStart ok, listening on %s", p.Conf.Params.Name, p.Conf.Params.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
				if p.status == falcon.APP_STATUS_EXIT {
					return
				}
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

func (p *Lb) rpcStop() (err error) {
	if p.rpcListener == nil {
		return falcon.ErrNoent
	}

	p.rpcListener.Close()
	p.rpcConnects.Lock()
	for e := p.rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	p.rpcConnects.Unlock()

	return nil
}
