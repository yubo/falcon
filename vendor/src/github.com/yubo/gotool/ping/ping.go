/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/* addr [4]byte use big endian */
package ping

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/gotool/list"
)

type task_entry struct {
	task   *Task
	s_list list.ListHead // for server.ips
	t_list list.ListHead // for task
	recv   *bool
}

type Task struct {
	sync.RWMutex               // for entry_list
	t_list       list.ListHead // server.task_list
	e_list       list.ListHead // Task_entry.t_list
	startTs      int64
	lastTs       int64
	Timeout      int64
	retry        int
	RetryLimit   int
	Ips          [][4]byte
	Ret          []bool
	Done         chan *Task
	Error        error
}

type serv struct {
	sync.RWMutex
	run     bool
	running chan struct{}
	tx_fd   int
	rx_fd   int
	rx_fp   *os.File
	rate    uint32
	ticker  *time.Ticker
	t_list  list.ListHead // task list
	ips     map[[4]byte]*list.ListHead
}

var (
	ErrAcces = errors.New("Permission denied (you must be root)")
	ErrNoRun = errors.New("ping server is not running")
	ErrEmpty = errors.New("empty list")
	server   serv
)

func Run(rate uint32) (err error) {
	server.Lock()
	defer server.Unlock()

	if os.Getuid() != 0 {
		return ErrAcces
	}

	server.tx_fd, err = syscall.Socket(syscall.AF_INET,
		syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return err
	}

	server.rx_fd, err = syscall.Socket(syscall.AF_INET,
		syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return err
	}
	server.rx_fp = os.NewFile(uintptr(server.rx_fd),
		fmt.Sprintf("ping/server/rx_fp"))

	server.rate = rate
	server.run = true
	server.running = make(chan struct{})
	server.ticker = time.NewTicker(time.Second / time.Duration(rate))
	server.ips = make(map[[4]byte]*list.ListHead)
	server.t_list.Init()

	go tx_loop(server.running, server.ticker.C)
	go rx_loop(server.rx_fp, server.running)

	return nil
}

func Kill() error {
	server.Lock()
	defer server.Unlock()

	if !server.run {
		return ErrNoRun
	}

	close(server.running)
	server.ticker.Stop()
	syscall.Close(server.tx_fd)
	server.rx_fp.Close()
	server.run = false

	return nil
}

func Go(ips [][4]byte, timeout, retry int, done chan *Task) *Task {
	var (
		s_list *list.ListHead
		ok     bool
	)
	t := &Task{
		Ips:        ips[:],
		Ret:        make([]bool, len(ips)),
		Timeout:    int64(timeout),
		RetryLimit: retry,
	}
	t.e_list.Init()

	if done == nil {
		done = make(chan *Task, 10) // buffered.
	} else {
		if cap(done) == 0 {
			log.Panic("ping: ping channel is unbuffered")
		}
	}
	t.Done = done

	if !server.run {
		t.Error = ErrNoRun
		t.Done <- t
		return t
	}

	server.Lock()
	t.Lock()

	for i, ip := range t.Ips {
		if s_list, ok = server.ips[ip]; !ok {
			s_list = &list.ListHead{}
			s_list.Init()
			server.ips[ip] = s_list
		}

		te := task_entry{}
		te.task = t
		te.recv = &t.Ret[i]
		t.e_list.AddTail(&te.t_list)
		s_list.AddTail(&te.s_list)
	}
	server.t_list.AddTail(&t.t_list)
	glog.V(4).Infof("task[%p] add to tail", t)

	t.Unlock()
	server.Unlock()

	return t
}

func Call(ips [][4]byte, timeout int, retry int) ([]bool, error) {
	task := <-Go(ips, timeout, retry, make(chan *Task, 1)).Done
	return task.Ret, task.Error
}

func list_to_task(list *list.ListHead) *Task {
	return (*Task)(unsafe.Pointer((uintptr(unsafe.Pointer(list)) -
		unsafe.Offsetof(((*Task)(nil)).t_list))))
}

/* for task_entry */
func list_to_entry(list *list.ListHead) *task_entry {
	return (*task_entry)(unsafe.Pointer((uintptr(unsafe.Pointer(list)) -
		unsafe.Offsetof(((*task_entry)(nil)).t_list))))
}

func list_to_entry_s(list *list.ListHead) *task_entry {
	return (*task_entry)(unsafe.Pointer((uintptr(unsafe.Pointer(list)) -
		unsafe.Offsetof(((*task_entry)(nil)).s_list))))
}

func (p *Task) next() (*Task, error) {
	now := time.Now().Unix()

	if p != nil {
		p.Lock()
		p.retry++
		p.lastTs = now
		p.t_list.Del()
		server.Lock()
		server.t_list.AddTail(&p.t_list)
		server.Unlock()
		p.Unlock()
		glog.V(4).Infof("task[%p] move to tail", p)
	}

	for n, _n := server.t_list.Next, server.t_list.Next.Next; n != &server.t_list; n, _n = _n, _n.Next {
		task := list_to_task(n)
		//glog.V(5).Infof("task[%p] check", task)

		if task.retry < task.RetryLimit {
			if now > task.startTs+task.Timeout {
				task.startTs = now
				glog.V(4).Infof("task[%p] get", task)
				return task, nil
			}
			glog.V(4).Infof("task[%p] skip wait until %d now %d",
				task, task.startTs+task.Timeout, now)
			continue
		}

		if now > task.lastTs+task.Timeout {
			/* return task -> done */
			server.Lock()
			task.Lock()

			task.t_list.Del()

			for pos, _pos := task.e_list.Next, task.e_list.Next.Next; pos != &task.e_list; pos, _pos = _pos, _pos.Next {

				/*
				 * Ugly Hack !
				 * rx loop may use this ptr, so ...
				 * remove this entry from e_list
				 * and keep next,prev ptr
				 */
				pos.Next.Prev = pos.Prev
				pos.Prev.Next = pos.Next

				/* TODO: check server.ips[ip], remove empty node */
				/* remove task_entry from server.ips list */
				list_to_entry(pos).s_list.Del()

			}
			task.Unlock()
			server.Unlock()
			task.Done <- task
			glog.V(4).Infof("task[%p] done", task)
		}
	}
	return nil, ErrEmpty
}

func tx_loop(running chan struct{}, C <-chan time.Time) {
	var (
		sa  syscall.SockaddrInet4
		idx int
		cur *Task
		err error
		ip  [4]byte
	)

	for {
		select {
		case <-running:
			glog.V(3).Infof("tx routine exit")
			return
		case <-C:
			if cur == nil {
				if cur, err = cur.next(); err != nil {
					//glog.V(5).Infof("tx loop: %s ", err)
					continue
				}
				idx = 0
			}

			for idx < len(cur.Ips) {
				if !cur.Ret[idx] {
					break
				}
				idx++
			}
			if idx >= len(cur.Ips) {
				if cur, err = cur.next(); err != nil {
					glog.V(4).Infof("tx loop: %s ", err)
					continue
				}
				idx = 0
			}
			ip = cur.Ips[idx]

			sa.Addr = ip
			glog.V(4).Infof("syscall.Sendto %d.%d.%d.%d",
				uint8(ip[0]), uint8(ip[1]),
				uint8(ip[2]), uint8(ip[3]))
			if err = syscall.Sendto(server.tx_fd, pkt([4]byte(ip)),
				0, &sa); err != nil {
				glog.V(2).Infof("syscall.Sendto error: %s ", err)

			}
			idx++
		}
	}
}

func rx_loop(fp *os.File, running chan struct{}) {
	var (
		buf     []byte = make([]byte, 1500)
		numRead int
		err     error
		src     [4]byte
	)
	for {
		select {
		case <-running:
			glog.V(3).Infof("rx routine exit")
			return
		default:
			/* recv icmp */
			if numRead, err = fp.Read(buf); err != nil {
				glog.Error(err)
			}
			if buf[0x14] == 0x00 {
				src[0], src[1], src[2], src[3] =
					buf[0x0c], buf[0x0d], buf[0x0e], buf[0x0f]
				glog.V(3).Infof("rx_loop read[%d] % X\n",
					numRead, buf[:numRead])
				glog.V(3).Infof("%d.%d.%d.%d",
					src[0], src[1], src[2], src[3])
				if h, ok := server.ips[src]; ok {
					for p := h.Next; p != h; p = p.Next {
						*(list_to_entry_s(p).recv) = true
					}
				}
			}
		}
	}
}
