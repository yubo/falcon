/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"os"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend/config"
	"github.com/yubo/falcon/backend/parse"
	fconfig "github.com/yubo/falcon/config"
)

const (
	CACHE_TIME              = 1800 //s
	FIRST_FLUSH_DISK        = 1    //s
	FLUSH_DISK_STEP         = 1    //s
	DEFAULT_HISTORY_SIZE    = 3
	CONN_RETRY              = 2
	CACHE_SIZE              = 1 << 5
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
	MODULE_NAME             = "\x1B[32m[BACKEND]\x1B[0m "
	CTRL_STEP               = 360

	C_CONN_TIMEOUT     = "conntimeout"
	C_CALL_TIMEOUT     = "calltimeout"
	C_WORKER_PROCESSES = "workerprocesses"
	C_HTTP_ENABLE      = "http_enable"
	C_HTTP_ADDR        = "httpaddr"
	C_RPC_ENABLE       = "rpc_enable"
	C_RPC_ADDR         = "rpcaddr"
	C_GRPC_ENABLE      = "grpc_enable"
	C_GRPC_ADDR        = "grpcaddr"
	C_IDX              = "idx"
	C_IDXINTERVAL      = "idxinterval"
	C_IDXFULLINTERVAL  = "idxfullinterval"
	C_DB_MAX_IDLE      = "dbmaxidle"
	C_DB_MAX_CONN      = "dbmaxconn"
	C_SHMMAGIC         = "shmmagic"
	C_SHMKEY           = "shmkey"
	C_SHMSIZE          = "shmsize"
	C_DSN              = "dsn"
	//C_HDISK            = "hdisk"
	C_PAYLOADSIZE = "payloadsize"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CONN_TIMEOUT:     "1000",
		C_CALL_TIMEOUT:     "5000",
		C_WORKER_PROCESSES: "2",
		C_HTTP_ENABLE:      "true",
		C_HTTP_ADDR:        "127.0.0.1:7021",
		C_RPC_ENABLE:       "true",
		C_RPC_ADDR:         "127.0.0.1:7020",
		C_GRPC_ENABLE:      "true",
		C_GRPC_ADDR:        "127.0.0.1:7022",
		C_IDX:              "true",
		C_IDXINTERVAL:      "30",
		C_IDXFULLINTERVAL:  "86400",
		C_DB_MAX_IDLE:      "4",
		C_SHMMAGIC:         "0x80386",
		C_SHMKEY:           "0x7020",
		C_SHMSIZE:          "0x10000000",
	}
)

// module {{{
type module interface {
	prestart(*Backend) error // alloc public data
	start(*Backend) error    // alloc private data, run private goroutine
	stop(*Backend) error     // free private data, private goroutine exit
	reload(*Backend) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

// }}}

// {{{ Backend
type Backend struct {
	Conf    *config.ConfBackend
	oldConf *config.ConfBackend
	// runtime
	status uint32

	//cacheModule
	cache *backendCache

	//storageModule
	//hdisk []string
	//storageNetTaskCh map[string]chan *netTask
	//storageIoTaskCh  []chan *ioTask

	ts           int64
	statTicker   chan time.Time
	timeTicker   chan time.Time
	commitTicker chan time.Time
}

func (p *Backend) New(conf interface{}) falcon.Module {
	return &Backend{
		Conf: conf.(*config.ConfBackend),
	}
}

func (p *Backend) Name() string {
	return p.Conf.Name
}

func (p *Backend) Parse(text []byte, filename string, lino int, debug bool) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfBackend)
	p.Conf.Configer.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
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
}

func (p *Backend) Stop() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT

	for i := len(modules) - 1; i >= 0; i-- {
		if e := modules[i].stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Backend) Reload(c interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = c.(*config.ConfBackend)

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
