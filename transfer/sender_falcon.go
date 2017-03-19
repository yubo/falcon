/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"errors"
	"net"
	"net/rpc"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type rpcClients struct {
	addr string
	cli  []*rpc.Client
}

type senderFalcon struct {
	name            string
	workerProcesses int
	connTimeout     int
	callTimeout     int
	payloadSize     int
	clients         map[string]rpcClients
	chans           map[string]chan *falcon.MetaData
}

func (p *senderFalcon) new(L *Transfer) sender {
	workerprocesses, _ := L.Conf.Configer.Int(falcon.C_WORKER_PROCESSES)
	conntimeout, _ := L.Conf.Configer.Int(falcon.C_CONN_TIMEOUT)
	calltimeout, _ := L.Conf.Configer.Int(falcon.C_CALL_TIMEOUT)
	payloadsize, _ := L.Conf.Configer.Int(falcon.C_PAYLOADSIZE)

	return &senderFalcon{
		name:            "falcon",
		workerProcesses: workerprocesses,
		connTimeout:     conntimeout,
		callTimeout:     calltimeout,
		payloadSize:     payloadsize,
		clients:         make(map[string]rpcClients),
		chans:           make(map[string]chan *falcon.MetaData),
	}
}

func (tpl *senderFalcon) dial(address string,
	timeout int) (*rpc.Client, error) {
	return rpcDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *senderFalcon) addClientChan(key, addr string,
	ch chan *falcon.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = rpcClients{
		addr: addr,
		cli:  make([]*rpc.Client, p.workerProcesses),
	}
	/*
		for i := 0; i < p.workerProcesses; i++ {
			p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
			if err != nil {
				glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n",
					key, addr, err)
			}
		}
	*/
	return nil
}

func (p *senderFalcon) start(name string) (err error) {
	for key, _ := range p.clients {
		for i := 0; i < p.workerProcesses; i++ {
			p.clients[key].cli[i], err = p.dial(p.clients[key].addr, p.connTimeout)
			if err != nil {
				glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n",
					key, p.clients[key].addr, err)
			}
		}
	}

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

func (p *senderFalcon) stop() error {
	for node, ch := range p.chans {
		close(ch)
		for i := 0; i < len(p.clients[node].cli); i++ {
			p.clients[node].cli[i].Close()
		}
	}
	return nil
}

func rpcDial(address string, timeout time.Duration) (*rpc.Client, error) {
	statsInc(ST_UPSTREAM_DIAL, 1)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial(falcon.ParseAddr(address))
	if err != nil {
		statsInc(ST_UPSTREAM_DIAL_ERR, 1)
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

	statsInc(ST_UPSTREAM_RECONNECT, 1)
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
		err = netRpcCall(*client, "Bkd.Put", items, resp,
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

// falcon {{{
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
				statsInc(ST_UPSTREAM_PUT, 1)
				statsInc(ST_UPSTREAM_PUT_ITEM, payloadSize)
				if err = putRpcBackendData(client, rrds,
					addr, connTimeout, callTimeout); err != nil {
					statsInc(ST_UPSTREAM_PUT_ERR, 1)
				}
				i = 0
			}
		}
	}
}

//}}}
