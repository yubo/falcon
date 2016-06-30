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
	"database/sql"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yubo/falcon/specs"
)

var (
	configFile string
	configPtr  unsafe.Pointer
	configOpts Options = defaultOptions
	exitChans  []chan chan error

	/* db */
	DB        *sql.DB
	dbLock    sync.RWMutex
	dbConnMap map[string]*sql.DB

	/* history */
	// mem:  front = = = back
	// time: latest ...  old
	//HistoryCache = tmap.NewSafeMap()

	/* linkedlist */
	Consistent  *consistent.Consistent
	Net_task_ch map[string]chan *Net_task_t
	clients     map[string][]*rpc.Client

	/* rpc */
	rpc_exit chan chan error
	connects conn_list

	/* http */
	http_exit chan chan error

	/* cache */
	cache cache_t

	/* rrdtool */
	sync_exit    chan chan error
	io_task_chan chan *io_task_t
)

func registerExitChans(e chan chan error) {
	exitChans = append(exitChans, e)
}

func init() {
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// core
	runtime.GOMAXPROCS(runtime.NumCPU())

	// rpc
	rpc_exit = make(chan chan error)
	connects = conn_list{list: list.New()}

	// http
	http_exit = make(chan chan error)
	httpRoutes()

	// rrdtool/sync_disk/migrate
	sync_exit = make(chan chan error)
	io_task_chan = make(chan *io_task_t, 16)
	Consistent = consistent.New()
	Net_task_ch = make(map[string]chan *Net_task_t)
	clients = make(map[string][]*rpc.Client)

	// store
	size := CACHE_TIME / FLUSH_DISK_STEP
	if size < 0 {
		log.Panicf("store.init, bad size %d\n", size)
	}

	// cache
	cache.hash = make(map[string]*cacheEntry)
}

func start_signal(pid int, cfg *Options) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			done := make(chan error)
			for i := len(exitChans); i > 0; i-- {
				exitChans[i-1] <- done
				if err := <-done; err != nil {
					log.Println(err)
				}
			}
			commitCaches(true)

			log.Printf("rrd data commit complete\n")
			log.Print("pid:%d exit\n", pid)
			os.Exit(0)
		}
	}
}

func Handle(arg interface{}) {
	parse(arg.(*specs.CmdOptions).ConfigFile)

	dbInit()
	rrdStart()
	rpcStart()
	//indexStart()
	httpStart()

	start_signal(os.Getpid(), config())
}
