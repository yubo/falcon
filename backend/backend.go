/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"os"
	"sync/atomic"
	"time"

	"stathat.com/c/consistent"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

const (
	CACHE_TIME              = 1800 //s
	FIRST_FLUSH_DISK        = 1    //s
	FLUSH_DISK_STEP         = 1    //s
	DEFAULT_HISTORY_SIZE    = 3
	CONN_RETRY              = 2
	CACHE_SIZE              = 1 << 5
	CACHE_SIZE_MASK         = CACHE_SIZE - 1
	DATA_TIMESTAMP_REGULATE = true
	INDEX_QPS               = 100
	INDEX_UPDATE_CYCLE_TIME = 86400
	INDEX_TIMEOUT           = 86400
	INDEX_TRASH_LOOPTIME    = 600
	INDEX_MAX_OPEN_CONNS    = 4
	DEBUG_MULTIPLES         = 20    // demo 时间倍数
	DEBUG_STEP              = 60    //
	DEBUG_SAMPLE_NB         = 18000 //单周期生成样本数量
	DEBUG_STAT_STEP         = 60
	MODULE_NAME             = "\x1B[32m[BACKEND]\x1B[0m "
	CTRL_STEP               = 360
)

var (
	modules []module
)

func init() {
	falcon.RegisterModule(falcon.GetType(falcon.ConfBackend{}), &Backend{})
	registerModule(&storageModule{})
	// cache should early register(init cache data)
	registerModule(&cacheModule{})
	registerModule(&httpModule{})
	registerModule(&rpcModule{})
	registerModule(&indexModule{})
	registerModule(&statsModule{})
	registerModule(&timerModule{})
}

// module {{{
type module interface {
	prestart(*Backend) error // alloc public data
	start(*Backend) error    // alloc private data, run private goroutine
	stop(*Backend) error     // free private data, private goroutine exit
	reload(*Backend) error   // try to keep the data, refresh configure
}

func registerModule(m module) {
	modules = append(modules, m)
}

// }}}

// {{{ Backend
type Backend struct {
	Conf    *falcon.ConfBackend
	oldConf *falcon.ConfBackend
	// runtime
	status uint32

	//cacheModule
	cache *backendCache

	//storageModule
	hdisk                    []string
	storageNetTaskCh         map[string]chan *netTask
	storageIoTaskCh          []chan *ioTask
	storageMigrateConsistent *consistent.Consistent

	ts           int64
	statTicker   chan time.Time
	timeTicker   chan time.Time
	commitTicker chan time.Time
}

func (p *Backend) New(conf interface{}) falcon.Module {
	return &Backend{
		Conf: conf.(*falcon.ConfBackend),
	}
}

func (p *Backend) Name() string {
	return p.Conf.Name
}

func (p *Backend) String() string {
	return p.Conf.String()
}

func (p *Backend) Prestart() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Init()", p.Conf.Name)
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		if e := modules[i].prestart(p); e != nil {
			//panic(err)
			err = e
			glog.Error(err)
		}
	}
	return err
	// core

	/*
		//cache
		p.cacheInit()

		// rpc
		p.rpcConnects = connList{list: list.New()}

		// http
		p.httpMux = http.NewServeMux()
		p.httpRoutes()

		// rrdtool/sync_disk/migrate
		p.storageNetTaskCh = make(map[string]chan *netTask)
		p.storageMigrateClients = make(map[string][]*rpc.Client)
		p.storageMigrateConsistent = consistent.New()

		// store
		size := CACHE_TIME / FLUSH_DISK_STEP
		if size < 0 {
			glog.Fatalf(MODULE_NAME+"store.init, bad size %d\n", size)
		}
		p.status = falcon.APP_STATUS_INIT
		return nil
	*/

}

func (p *Backend) Start() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Name)
	p.status = falcon.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		if e := modules[i].start(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	p.status = falcon.APP_STATUS_RUNNING
	return err

	/*
		p.running = make(chan struct{}, 0)
		p.timeStart()
		p.rrdStart()
		p.rpcStart()
		p.indexStart()
		p.httpStart()
		p.statStart()
		p.cacheStart()
	*/
}

func (p *Backend) Stop() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT

	for i, n := 0, len(modules); i < n; i++ {
		if e := modules[n-i].stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
	/*
		close(p.running)
		p.cacheStop()
		p.statStop()
		p.httpStop()
		p.indexStop()
		p.rpcStop()
		p.rrdStop()
		p.timeStop()
	*/
}

func (p *Backend) Reload(config interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = config.(*falcon.ConfBackend)

	for i := 0; i < len(modules); i++ {
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Backend) Signal(sig os.Signal) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return err
}

func (p *Backend) timeNow() int64 {
	if p.ts != 0 {
		return atomic.LoadInt64(&p.ts)
	} else {
		return time.Now().Unix()
	}
}

//}}}
