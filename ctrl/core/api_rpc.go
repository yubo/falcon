/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"container/list"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
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

type CTRL struct {
	ctrl *Ctrl
}

func (p *CTRL) Ping(req specs.Null, resp *specs.RpcResp) error {
	return nil
}

func (p *CTRL) ListLb(req specs.Null, resp *[]string) error {
	*resp = p.ctrl.Lbs
	return nil
}

func (p *CTRL) ListBackend(req specs.Null, resp *[]specs.Backend) error {
	*resp = p.ctrl.Backends
	return nil
}

func (p *CTRL) ListMigrate(req specs.Null, resp *specs.Migrate) error {
	*resp = p.ctrl.Migrate
	return nil
}

func (p *Ctrl) rpcStart() (err error) {

	var addr *net.TCPAddr
	if !p.Params.Rpc {
		return nil
	}
	ctrl := &CTRL{ctrl: p}
	rpc.Register(ctrl)

	addr, err = net.ResolveTCPAddr("tcp", p.Params.RpcAddr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, net.ResolveTCPAddr failed, %s",
			err)
	}

	p.rpcListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf(MODULE_NAME+"rpc.Start error, listen %s failed, %s",
			p.Params.RpcAddr, err)
	} else {
		glog.Infof(MODULE_NAME+"%s rpcStart ok, listening on %s",
			p.Params.Name, p.Params.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
				if p.status == specs.APP_STATUS_EXIT {
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

func (p *Ctrl) rpcStop() (err error) {
	if p.rpcListener == nil {
		return specs.ErrNoent
	}

	p.rpcListener.Close()
	p.rpcConnects.Lock()
	for e := p.rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	p.rpcConnects.Unlock()

	return nil
}
