/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

type netClients struct {
	addr string
	cli  []net.Conn
}

type senderTsdb struct {
	name            string
	workerProcesses int
	connTimeout     int
	callTimeout     int
	payloadSize     int
	clients         map[string]netClients
	chans           map[string]chan *falcon.MetaData
}

func (p *senderTsdb) new(L *Transfer) sender {
	workerprocesses, _ := L.Conf.Configer.Int(C_WORKER_PROCESSES)
	conntimeout, _ := L.Conf.Configer.Int(C_CONN_TIMEOUT)
	calltimeout, _ := L.Conf.Configer.Int(C_CALL_TIMEOUT)
	payloadsize, _ := L.Conf.Configer.Int(C_PAYLOADSIZE)

	return &senderTsdb{
		name:            "tsdb",
		workerProcesses: workerprocesses,
		connTimeout:     conntimeout,
		callTimeout:     calltimeout,
		payloadSize:     payloadsize,
		clients:         make(map[string]netClients),
		chans:           make(map[string]chan *falcon.MetaData),
	}
}

func netDial(address string, timeout time.Duration) (net.Conn, error) {
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
	return conn, err
}

func (tpl *senderTsdb) dial(address string, timeout int) (net.Conn, error) {
	return netDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *senderTsdb) addClientChan(key, addr string,
	ch chan *falcon.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = netClients{
		addr: addr,
		cli:  make([]net.Conn, p.workerProcesses),
	}
	for i := 0; i < p.workerProcesses; i++ {
		p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
		if err != nil {
			glog.Fatalf(MODULE_NAME+"node:%s addr:%s err:%s\n", key, addr, err)
		}
	}
	return nil
}

func (p *senderTsdb) start(name string) error {
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

func (p *senderTsdb) stop() error {
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
				statsInc(ST_UPSTREAM_PUT, 1)
				statsInc(ST_UPSTREAM_PUT_ITEM, payloadSize)
				if err = putTsdbData(client, items,
					addr, connTimeout, callTimeout); err != nil {
					statsInc(ST_UPSTREAM_PUT_ERR, 1)
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

	statsInc(ST_UPSTREAM_RECONNECT, 1)
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
