/*
 * Copyright 2016 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"errors"
	"math/rand"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type upstreamModule struct {
	running chan struct{}
	idx     uint32
	size    uint32
	pool    []string
}

func (p *upstreamModule) get() string {
	return p.pool[atomic.AddUint32(&p.idx, 1)%p.size]
}

func rpcDial(address string,
	timeout time.Duration) (*rpc.Client, error) {

	statsInc(ST_UPSTREAM_DIAL, 1)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
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
	return jsonrpc.NewClient(conn), err
}

func (p *upstreamModule) reconnection(client **rpc.Client, connTimeout int) {
	var err error

	statsInc(ST_UPSTREAM_RECONNECT, 1)
	addr := p.get()
	if *client != nil {
		(*client).Close()
	}

	*client, err = rpcDial(addr, time.Duration(connTimeout)*
		time.Millisecond)

	for err != nil {
		//danger!! block routine
		glog.Infof(MODULE_NAME+"%s", err)

		time.Sleep(time.Millisecond * 500)
		addr = p.get()
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

func (p *upstreamModule) putRpcStorageData(client **rpc.Client, items *[]*falcon.MetaData,
	connTimeout, callTimeout int) error {
	var (
		err  error
		resp *falcon.RpcResp
		i    int
	)

	statsInc(ST_UPSTREAM_UPDATE, 1)
	statsInc(ST_UPSTREAM_UPDATE_ITEM, len(*items))
	resp = &falcon.RpcResp{}

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "LB.Update", *items, resp,
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			glog.V(3).Infof(MODULE_NAME+"send success %d", len(*items))
			goto out
		}
		glog.V(3).Info(MODULE_NAME, err)
		if err == rpc.ErrShutdown {
			p.reconnection(client, connTimeout)
		}
	}
out:
	if err != nil {
		statsInc(ST_UPSTREAM_UPDATE_ERR, 1)
	}
	return err
}

func (p *upstreamModule) prestart(agent *Agent) error {
	p.running = make(chan struct{}, 0)
	return nil
}

func (p *upstreamModule) start(agent *Agent) error {

	backend := strings.Split(agent.Conf.Configer.Str(falcon.C_UPSTREAM), ",")
	connTimeout, _ := agent.Conf.Configer.Int(falcon.C_CONN_TIMEOUT)
	callTimeout, _ := agent.Conf.Configer.Int(falcon.C_CALL_TIMEOUT)
	payloadSize, _ := agent.Conf.Configer.Int(falcon.C_PAYLOADSIZE)
	debug := agent.Conf.Debug

	p.size = uint32(len(backend))
	p.idx = rand.Uint32() % p.size
	p.pool = make([]string, int(p.size))
	copy(p.pool, backend)

	if debug > 1 {
		go func() {
			for {
				select {
				case _, ok := <-p.running:
					if !ok {
						return
					}
				case items := <-agent.appUpdateChan:
					for k, v := range *items {
						glog.V(3).Infof(MODULE_NAME+"%d %s", k, v)
					}
				}
			}
		}()
		return nil
	}

	go func() {

		client, err := rpcDial(p.get(),
			time.Duration(connTimeout)*
				time.Millisecond)
		if err != nil {
			p.reconnection(&client, connTimeout)
		}

		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case items := <-agent.appUpdateChan:
				n := payloadSize
				i := 0
				for ; i < len(*items)-n; i += n {
					_items := (*items)[i : i+n]
					p.putRpcStorageData(&client, &_items,
						connTimeout,
						callTimeout)
				}
				if i < len(*items) {
					_items := (*items)[i:]
					p.putRpcStorageData(&client, &_items,
						connTimeout,
						callTimeout)
				}
			}

		}
	}()
	return nil
}

func (p *upstreamModule) stop(agent *Agent) error {
	close(p.running)
	return nil
}

func (p *upstreamModule) reload(agent *Agent) error {
	// TODO
	return nil
}
