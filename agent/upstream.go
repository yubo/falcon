/*
 * Copyright 2016 yubo. All rights reserved.
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
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	upstreamConfig AgentOpts
	streamPool     upstreamPool
	idxClient      int
)

type upstreamPool struct {
	idx  uint32
	size uint32
	pool []string
}

func (p *upstreamPool) get() string {
	return p.pool[atomic.AddUint32(&p.idx, 1)%p.size]
}

func rpcDial(address string,
	timeout time.Duration) (*rpc.Client, error) {

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
	return jsonrpc.NewClient(conn), err
}

func reconnection(client **rpc.Client, connTimeout int) {
	var err error

	statInc(ST_UPSTREAM_RECONNECT, 1)
	addr := streamPool.get()
	if *client != nil {
		(*client).Close()
	}

	*client, err = rpcDial(addr, time.Duration(connTimeout)*
		time.Millisecond)

	for err != nil {
		//danger!! block routine
		glog.Infof("%s", err)

		time.Sleep(time.Millisecond * 500)
		addr = streamPool.get()
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

func putRpcStorageData(client **rpc.Client, items *[]*specs.MetaData,
	connTimeout, callTimeout int) error {
	var (
		err  error
		resp *specs.RpcResp
		i    int
	)

	statInc(ST_UPSTREAM_UPDATE, 1)
	statInc(ST_UPSTREAM_UPDATE_ITEM, len(*items))
	resp = &specs.RpcResp{}

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Falcon.Update", *items, resp,
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			glog.V(3).Infof("send success %d", len(*items))
			goto out
		}
		glog.V(3).Info(err)
		if err == rpc.ErrShutdown {
			reconnection(client, connTimeout)
		}
	}
out:
	if err != nil {
		statInc(ST_UPSTREAM_UPDATE_ERR, 1)
	}
	return err
}

func upstreamStart(config AgentOpts, p *specs.Process) {
	var (
		client *rpc.Client
		err    error
		i      int
	)
	upstreamConfig = config

	streamPool = upstreamPool{}
	streamPool.size = uint32(len(upstreamConfig.Handoff.Upstreams))
	streamPool.idx = rand.Uint32() % streamPool.size
	streamPool.pool = make([]string, int(streamPool.size))
	copy(streamPool.pool, upstreamConfig.Handoff.Upstreams)

	if upstreamConfig.Debug > 1 {
		go func() {
			for {
				items := <-appUpdateChan
				for k, v := range *items {
					glog.V(3).Infof("%d %s", k, v)
				}
			}
		}()
		return
	}

	go func() {

		client, err = rpcDial(streamPool.get(),
			time.Duration(upstreamConfig.Handoff.ConnTimeout)*
				time.Millisecond)
		if err != nil {
			reconnection(&client, upstreamConfig.Handoff.ConnTimeout)
		}

		for {
			items := <-appUpdateChan

			n := upstreamConfig.Handoff.Batch
			for i = 0; i < len(*items)-n; i += n {
				_items := (*items)[i : i+n]
				putRpcStorageData(&client, &_items,
					upstreamConfig.Handoff.ConnTimeout,
					upstreamConfig.Handoff.CallTimeout)
			}
			if i < len(*items) {
				_items := (*items)[i:]
				putRpcStorageData(&client, &_items,
					upstreamConfig.Handoff.ConnTimeout,
					upstreamConfig.Handoff.CallTimeout)
			}

		}
	}()

}
