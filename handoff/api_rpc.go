/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

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

type Falcon struct {
	handoff *Handoff
}

func (p *Falcon) Ping(req specs.Null, resp *specs.RpcResp) error {
	return nil
}

func (p *Falcon) Update(args []*specs.MetaData,
	reply *specs.HandoffResp) error {
	reply.Invalid = 0
	now := time.Now().Unix()

	items := []*specs.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.Name == "" || v.Host == "" {
			reply.Invalid += 1
			continue
		}

		if v.Type != specs.COUNTER &&
			v.Type != specs.GAUGE &&
			v.Type != specs.DERIVE {
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

		items = append(items, &specs.MetaData{
			Name: v.Name,
			Host: v.Host,
			Ts:   v.Ts,
			Step: v.Step,
			Type: v.Type,
			Tags: sortTags(v.Tags),
		})
	}

	p.handoff.appUpdateChan <- &items
	glog.V(3).Infof("recv %d", len(items))

	reply.Message = "ok"
	reply.Total = len(args)

	statInc(ST_RPC_UPDATE, 1)
	statInc(ST_RPC_UPDATE_CNT, len(items))
	statInc(ST_RPC_UPDATE_ERR, reply.Invalid)

	return nil
}

func (p *Handoff) rpcStart() (err error) {

	var addr *net.TCPAddr
	if !p.Rpc {
		return nil
	}
	falcon := &Falcon{handoff: p}
	rpc.Register(falcon)

	addr, err = net.ResolveTCPAddr("tcp", p.RpcAddr)
	if err != nil {
		glog.Fatalf("rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	p.rpcListener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf("rpc.Start error, listen %s failed, %s",
			p.RpcAddr, err)
	} else {
		glog.Infof("%s rpcStart ok, listening on %s", p.Name, p.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := p.rpcListener.Accept()
			if err != nil {
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

func (p *Handoff) rpcStop() (err error) {
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

/*
func rpcStart(config HandoffOpts, p *specs.Process) {
	var rpcListener *net.TCPListener

	rpc.Register(new(Falcon))
	p.RegisterEvent("rpc", rpcEvent)
	rpcConfig = config

	rpcStart(&rpcConfig, &rpcListener)

	go func() {
		select {
		case event := <-rpcEvent:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				rpcStop(&rpcConfig, rpcListener)
				event.Done <- nil

				return
			} else if event.Method == specs.ROUTINE_EVENT_M_RELOAD {
				rpcStop(&rpcConfig, rpcListener)

				glog.V(3).Infof("old:\n%s\n new:\n%s",
					rpcConfig, appConfig)
				rpcConfig = appConfig
				_rpcStart(&rpcConfig, &rpcListener)
			}
		}
	}()
}
*/
