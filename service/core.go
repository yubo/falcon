/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/service/config"
	"github.com/yubo/falcon/service/parse"
	"github.com/yubo/gotool/list"
)

const (
	CACHE_TIME              = 1800 //s
	FIRST_FLUSH_DISK        = 1    //s
	FLUSH_DISK_STEP         = 1    //s
	DEFAULT_HISTORY_SIZE    = 3
	CONN_RETRY              = 2
	CACHE_SIZE              = 1 << 4
	CACHE_SIZE_MASK         = CACHE_SIZE - 1
	DATA_TIMESTAMP_REGULATE = false
	INDEX_QPS               = 100
	INDEX_UPDATE_CYCLE_TIME = 86400
	INDEX_TIMEOUT           = 86400
	INDEX_TRASH_LOOPTIME    = 600
	INDEX_MAX_OPEN_CONNS    = 4
	DEBUG_MULTIPLES         = 20    // demo 时间倍数
	DEBUG_STEP              = 60    //
	DEBUG_SAMPLE_NB         = 18000 //单周期生成样本数量
	DEBUG_STAT_STEP         = 60
	MODULE_NAME             = "\x1B[32m[SERVICE]\x1B[0m "
	CTRL_STEP               = 360

	C_CONN_TIMEOUT    = "conntimeout"
	C_CALL_TIMEOUT    = "calltimeout"
	C_API_ADDR        = "apiaddr"
	C_HTTP_ADDR       = "httpaddr"
	C_IDX             = "idx"
	C_IDXINTERVAL     = "idxinterval"
	C_IDXFULLINTERVAL = "idxfullinterval"
	C_DB_MAX_IDLE     = "dbmaxidle"
	C_DB_MAX_CONN     = "dbmaxconn"
	C_DSN             = "dsn"
	//C_WORKER_PROCESSES = "workerprocesses"
	//C_SHMMAGIC         = "shmmagic"
	//C_SHMKEY           = "shmkey"
	//C_SHMSIZE          = "shmsize"
	//C_HDISK            = "hdisk"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CONN_TIMEOUT:    "1000",
		C_CALL_TIMEOUT:    "5000",
		C_IDX:             "true",
		C_IDXINTERVAL:     "30",
		C_IDXFULLINTERVAL: "86400",
		C_DB_MAX_IDLE:     "4",
		//C_HTTP_ADDR:       "127.0.0.1:7021",
		//C_API_ADDR:        "127.0.0.1:7020",
		//C_WORKER_PROCESSES: "2",
		//C_SHMMAGIC:         "0x80386",
		//C_SHMKEY:           "0x7020",
		//C_SHMSIZE:          "0x10000000",
	}
)

// module {{{
type module interface {
	prestart(*Service) error // alloc public data
	start(*Service) error    // alloc private data, run private goroutine
	stop(*Service) error     // free private data, private goroutine exit
	reload(*Service) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

// }}}

// {{{ Service
type Service struct {
	Conf    *config.Service
	oldConf *config.Service
	// runtime
	status uint32

	//cacheModule
	cache *serviceCache

	//storageModule
	//hdisk []string
	//storageNetTaskCh map[string]chan *netTask
	//storageIoTaskCh  []chan *ioTask

	ts int64
}

func (p *Service) New(conf interface{}) falcon.Module {
	return &Service{
		Conf: conf.(*config.Service),
	}
}

func (p *Service) Name() string {
	return p.Conf.Name
}

func (p *Service) Parse(text []byte, filename string, lino int, debug bool) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.Service)
	p.Conf.Configer.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Service) String() string {
	return p.Conf.String()
}

func (p *Service) Prestart() (err error) {
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
}

func (p *Service) Start() (err error) {
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
}

func (p *Service) Stop() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		if e := modules[n-i-1].stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Service) Reload(c interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Service)

	for i := 0; i < len(modules); i++ {
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Service) Signal(sig os.Signal) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return err
}

func (p *Service) timeNow() int64 {
	if p.ts != 0 {
		return atomic.LoadInt64(&p.ts)
	} else {
		return time.Now().Unix()
	}
}

//}}}

// called by rpc
func (p *Service) createEntry(key string, item *falcon.Item) (*cacheEntry, error) {
	var (
		e     *cacheEntry
		ok    bool
		cache *serviceCache
	)

	cache = p.cache

	statsInc(ST_CACHE_CREATE, 1)
	if e, ok = cache.data[key]; ok {
		return e, falcon.ErrExist
	}

	e = &cacheEntry{
		createTs:  p.timeNow(),
		endpoint:  item.Endpoint,
		metric:    item.Metric,
		tags:      item.Tags,
		typ:       item.Type,
		hashkey:   key,
		timestamp: make([]int64, CACHE_SIZE),
		value:     make([]float64, CACHE_SIZE),
		//heartbeat: item.Heartbeat,
		//min:       []byte(item.Min)[0],
		//max:       []byte(item.Max)[0],
	}

	cache.Lock()
	cache.data[key] = e
	cache.Unlock()

	cache.idx0q.enqueue(&e.list_idx)

	return e, nil
}

func (p *Service) getItems(key string) (ret []*falcon.Item) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItems(CACHE_SIZE)
}

func (p *Service) getLastItem(key string) (ret *falcon.Item) {
	e := p.cache.get(key)
	if e == nil {
		return
	}
	return e.getItem()
}

func list_idx_entry(l *list.ListHead) *cacheEntry {
	return (*cacheEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*cacheEntry)(nil)).list_idx))))
}
