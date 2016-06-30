/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"container/list"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type conn_list struct {
	sync.RWMutex
	list *list.List
}

func (l *conn_list) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}
func (l *conn_list) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

func rpcStart() {
	if !config().Rpc {
		log.Println("rpc.Start warning, not enabled")
		return
	}
	addr := config().RpcAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("rpc.Start error, listen %s failed, %s", addr, err)
	} else {
		log.Println("rpc.Start ok, listening on", addr)
	}

	rpc.Register(new(Graph))
	registerExitChans(rpc_exit)

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := listener.Accept()
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
				e := connects.insert(conn)
				defer connects.remove(e)
				rpc.ServeConn(conn)
			}()
		}
	}()

	go func() {
		select {
		case done := <-rpc_exit:
			log.Println("rpc, recv sigout and exiting...")
			listener.Close()
			done <- nil

			connects.Lock()
			for e := connects.list.Front(); e != nil; e = e.Next() {
				e.Value.(net.Conn).Close()
			}
			connects.Unlock()

			return
		}
	}()

}
