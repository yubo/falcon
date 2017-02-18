/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type backend struct {
	name      string
	streams   upstream
	scheduler lbScheduler
}

func (p *backend) String() string {
	return fmt.Sprintf("%s", p.name)
}

func (p *backend) start(o falcon.Backend) error {
	for node, addr := range o.Upstreams {
		ch := make(chan *falcon.MetaData)
		p.scheduler.addChan(node, ch)
		p.streams.addClientChan(node, addr, ch)
	}
	return nil
}

/* upstream */
type rpcClients struct {
	addr string
	cli  []*rpc.Client
}

type netClients struct {
	addr string
	cli  []net.Conn
}

type upstream interface {
	start(string) error
	stop() error
	addClientChan(string, string, chan *falcon.MetaData) error
}

type upstreamFalcon struct {
	name        string
	concurrency int
	connTimeout int
	callTimeout int
	payloadSize int
	clients     map[string]rpcClients
	chans       map[string]chan *falcon.MetaData
}

func newUpstreamFalcon(p *Lb) *upstreamFalcon {
	return &upstreamFalcon{
		name:        "falcon",
		concurrency: p.Conf.Params.Concurrency,
		connTimeout: p.Conf.Params.ConnTimeout,
		callTimeout: p.Conf.Params.CallTimeout,
		payloadSize: p.Conf.PayloadSize,
		clients:     make(map[string]rpcClients),
		chans:       make(map[string]chan *falcon.MetaData),
	}
}

func (tpl *upstreamFalcon) dial(address string,
	timeout int) (*rpc.Client, error) {
	return rpcDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *upstreamFalcon) addClientChan(key, addr string,
	ch chan *falcon.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = rpcClients{
		addr: addr,
		cli:  make([]*rpc.Client, p.concurrency),
	}
	for i := 0; i < p.concurrency; i++ {
		p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
		if err != nil {
			glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n",
				key, addr, err)
		}
	}
	return nil
}

func (p *upstreamFalcon) start(name string) error {
	for node, ch := range p.chans {
		for i := 0; i < len(p.clients[node].cli); i++ {
			go falconUpstreamWorker(name, i,
				ch, &p.clients[node].cli[i],
				p.clients[node].addr,
				p.connTimeout, p.callTimeout, p.payloadSize)
		}
	}
	return nil
}

func (p *upstreamFalcon) stop() error {
	for node, ch := range p.chans {
		close(ch)
		for i := 0; i < len(p.clients[node].cli); i++ {
			p.clients[node].cli[i].Close()
		}
	}
	return nil
}

func rpcDial(address string, timeout time.Duration) (*rpc.Client, error) {
	statInc(ST_UPSTREAM_DIAL, 1)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		statInc(ST_UPSTREAM_DIAL_ERR, 1)
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return rpc.NewClient(conn), err
}

func reconnection(client **rpc.Client, addr string, connTimeout int) {
	var err error

	statInc(ST_UPSTREAM_RECONNECT, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = rpcDial(addr, time.Duration(connTimeout)*time.Millisecond)

	for err != nil {
		//danger!! block routine
		glog.Infof(MODULE_NAME+"reconnection to %s %s", addr, err)
		time.Sleep(time.Millisecond * 500)
		*client, err = rpcDial(addr,
			time.Duration(connTimeout)*time.Millisecond)
	}
}

func netRpcCall(client *rpc.Client, method string, args interface{},
	reply interface{}, timeout time.Duration) error {
	done := make(chan *rpc.Call, 1)
	client.Go(method, args, reply, done)
	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[rpc]")
	case call := <-done:
		if call.Error == nil {
			return nil
		} else {
			return call.Error
		}
	}
}

func putRpcBackendData(client **rpc.Client, items []*falcon.RrdItem,
	addr string, connTimeout, callTimeout int) error {
	var (
		err  error
		resp *falcon.RpcResp
		i    int
	)

	resp = &falcon.RpcResp{}

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Backend.Put", items, resp,
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			glog.V(3).Infof(MODULE_NAME+"send %d %s", len(items), addr)
			goto out
		}
		glog.V(3).Infof(MODULE_NAME+"send to %s %s", addr, err)
		if err == rpc.ErrShutdown {
			reconnection(client, addr, connTimeout)
		}
	}
out:
	return err
}

func falconUpstreamWorker(name string, idx int, ch chan *falcon.MetaData,
	client **rpc.Client, addr string, connTimeout, callTimeout, payloadSize int) {
	var err error
	var i int
	rrds := make([]*falcon.RrdItem, payloadSize)
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return
			}
			if rrds[i], err = item.Rrd(); err != nil {
				continue
			}
			i++
			if i == payloadSize {
				statInc(ST_UPSTREAM_PUT, 1)
				statInc(ST_UPSTREAM_PUT_ITEM, payloadSize)
				if err = putRpcBackendData(client, rrds,
					addr, connTimeout, callTimeout); err != nil {
					statInc(ST_UPSTREAM_PUT_ERR, 1)
				}
				i = 0
			}
		}
	}
}

// TSDB
type upstreamTsdb struct {
	upstreamFalcon
	name        string
	concurrency int
	connTimeout int
	callTimeout int
	clients     map[string]netClients
	chans       map[string]chan *falcon.MetaData
}

func newUpstreamTsdb(p *Lb) *upstreamTsdb {
	return &upstreamTsdb{
		name:        "tsdb",
		concurrency: p.Conf.Params.Concurrency,
		connTimeout: p.Conf.Params.ConnTimeout,
		callTimeout: p.Conf.Params.CallTimeout,
		clients:     make(map[string]netClients),
		chans:       make(map[string]chan *falcon.MetaData),
	}
}

func netDial(address string, timeout time.Duration) (net.Conn, error) {
	statInc(ST_UPSTREAM_DIAL, 1)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		statInc(ST_UPSTREAM_DIAL_ERR, 1)
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, err
}

func (tpl *upstreamTsdb) dial(address string, timeout int) (net.Conn, error) {
	return netDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *upstreamTsdb) addClientChan(key, addr string,
	ch chan *falcon.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = netClients{
		addr: addr,
		cli:  make([]net.Conn, p.concurrency),
	}
	for i := 0; i < p.concurrency; i++ {
		p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
		if err != nil {
			glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n", key, addr, err)
		}
	}
	return nil
}

func (p *upstreamTsdb) start(name string) error {
	for node, ch := range p.chans {
		for i := 0; i < len(p.clients[node].cli); i++ {
			go tsdbUpstreamWorker(name, i, ch,
				&p.clients[node].cli[i],
				p.clients[node].addr,
				p.connTimeout,
				p.callTimeout,
				p.payloadSize)
		}
	}
	return nil
}

func (p *upstreamTsdb) stop() error {
	for node, ch := range p.chans {
		close(ch)
		for i := 0; i < len(p.clients[node].cli); i++ {
			p.clients[node].cli[i].Close()
		}
	}
	return nil
}

func tsdbUpstreamWorker(name string, idx int,
	ch chan *falcon.MetaData, client *net.Conn,
	addr string, connTimeout, callTimeout, payloadSize int) {
	var err error
	var i int
	items := make([]*falcon.MetaData, payloadSize)
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return
			}
			items[i] = item
			i++
			if i == payloadSize {
				statInc(ST_UPSTREAM_PUT, 1)
				statInc(ST_UPSTREAM_PUT_ITEM, payloadSize)
				if err = putTsdbData(client, items,
					addr, connTimeout, callTimeout); err != nil {
					statInc(ST_UPSTREAM_PUT_ERR, 1)
				}
			}
		}
	}
}

func putTsdbData(client *net.Conn, items []*falcon.MetaData,
	addr string, connTimeout, callTimeout int) (err error) {
	var (
		i          int
		tsdbBuffer bytes.Buffer
	)

	for _, item := range items {
		tsdbBuffer.WriteString(item.Tsdb().TsdbString())
		tsdbBuffer.WriteString("\n")
	}

	for i = 0; i < CONN_RETRY; i++ {
		err = tsdbSend(*client, tsdbBuffer.Bytes(),
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			goto out
		} else {
			tsdbReconnection(client, addr, connTimeout)
		}
	}
out:
	return err

}

func tsdbSend(client net.Conn, data []byte, timeout time.Duration) (err error) {
	done := make(chan error)
	go func() {
		_, err = client.Write(data)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[tsdb]")
	case err = <-done:
		if err != nil {
			return fmt.Errorf("call failed, err %v", err)
		}
		return nil
	}
}

func tsdbReconnection(client *net.Conn, addr string, connTimeout int) {
	var err error

	statInc(ST_UPSTREAM_RECONNECT, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = netDial(addr,
		time.Duration(connTimeout)*time.Millisecond)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = netDial(addr,
			time.Duration(connTimeout)*time.Millisecond)
	}
}

func (p *Lb) upstreamWorkerStart() error {
	for _, b := range p.bs {
		b.streams.start(b.name)
	}
	return nil
}

func (p *Lb) upstreamWorkerStop() error {
	for _, b := range p.bs {
		b.streams.stop()
	}
	return nil
}

func (p *Lb) loadBalancerWorkerStart() {
	var (
		ok    bool
		b     *backend
		item  *falcon.MetaData
		items *[]*falcon.MetaData
		ch    chan *falcon.MetaData
	)

	// upstreams
	p.appUpdateChan = make(chan *[]*falcon.MetaData, 16)

	go func() {
		for {
			if items, ok = <-p.appUpdateChan; !ok {
				return
			}
			for _, b = range p.bs {
				for _, item = range *items {
					ch = b.scheduler.sched(item.Id())
					ch <- item
				}
			}
		}
	}()
}

func (p *Lb) loadBalancerWorkerStop() {
	close(p.appUpdateChan)
}

func (p *Lb) upstreamStart() error {
	p.bs = make([]*backend, 0)
	for _, v := range p.Conf.Backends {
		if v.Disabled {
			continue
		}
		b := &backend{name: v.Name}
		if v.Type == "falcon" {
			b.streams = newUpstreamFalcon(p)
		} else if v.Type == "tsdb" {
			b.streams = newUpstreamTsdb(p)
		} else {
			glog.Fatal(MODULE_NAME, falcon.ErrUnsupported)
		}

		//if v.Sched == "consistent" {
		b.scheduler = newSchedConsistent()
		//} else {
		//	glog.Fatal(falcon.ErrUnsupported)
		//}

		b.start(v)
		p.bs = append(p.bs, b)
	}

	glog.V(3).Infof(MODULE_NAME+"%s upstreamStart len(bs) %d", p.Conf.Params.Name, len(p.bs))

	p.upstreamWorkerStart()
	p.loadBalancerWorkerStart()
	return nil
}

func (p *Lb) upstreamStop() error {
	p.loadBalancerWorkerStop()
	p.upstreamWorkerStop()
	return nil
}
