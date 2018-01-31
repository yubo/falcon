/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm"
	"github.com/yubo/falcon/service/config"
)

const (
	CACHE_DATA_SIZE       = 1 << 4
	CACHE_DATA_SIZE_MASK  = CACHE_DATA_SIZE - 1
	INDEX_QPS             = 100
	INDEX_UPDATE_INTERVAL = 3600 * 24
	CACHE_EXPIRE_TIME     = 3600 * 26 // 超时时间 (s)
	CACHE_CLEAN_INTERVAL  = 3600      // 回收检查间隔时间 (s)
	MODULE_NAME           = "\x1B[36m[SERVICE]\x1B[0m"

	C_CALL_TIMEOUT   = "calltimeout"
	C_API_ADDR       = "apiaddr"
	C_HTTP_ADDR      = "httpaddr"
	C_IDX            = "idx"
	C_DSN            = "dsn"
	C_IDX_DSN        = "idxdsn"
	C_DB_MAX_IDLE    = "dbmaxidle"
	C_DB_MAX_CONN    = "dbmaxconn"
	C_CONF_INTERVAL  = "confinterval"
	C_SHARD_IDS      = "shardids"
	C_JUDGE_INTERVAL = "judgeinterval"
	C_JUDGE_NUM      = "judgenum"
	C_ALARM_ADDR     = "alarmaddr"
	C_UPSTREAM       = "upstream"
	C_BURST_SIZE     = "burstsize"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CALL_TIMEOUT:  "5000",
		C_IDX:           "true",
		C_DB_MAX_IDLE:   "4",
		C_DB_MAX_CONN:   "4",
		C_JUDGE_NUM:     "8",
		C_BURST_SIZE:    "32",
		C_CONF_INTERVAL: "600",
	}
)

type module interface {
	prestart(*Service) error // alloc public data
	start(*Service) error    // alloc private data, run private goroutine
	stop(*Service) error     // free private data, private goroutine exit
	reload(*Service) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type Service struct {
	Conf    *config.Service
	oldConf *config.Service
	// runtime
	status uint32

	//shardModule
	cache *CacheModule

	// tsdb
	tsdb *TsdbModule

	// event_trigger
	eventChan chan *alarm.Event

	//storageModule
	//hdisk []string
	//storageNetTaskCh map[string]chan *netTask
	//storageIoTaskCh  []chan *ioTask
}

func (p *Service) New(conf interface{}) falcon.Module {
	return &Service{
		Conf:      conf.(*config.Service),
		eventChan: make(chan *alarm.Event, 1024),
	}
}

func (p *Service) Name() string {
	return p.Conf.Name
}

func (p *Service) Parse(text []byte, filename string, lino int) falcon.ModuleConf {
	p.Conf = config.Parse(text, filename, lino).(*config.Service)
	p.Conf.Configer.Set(falcon.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Service) String() string {
	return p.Conf.String()
}

func (p *Service) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			panic(err)
			//glog.Error(err)
		}
	}
	return err
}

func (p *Service) Start() (err error) {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = falcon.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].start(p); err != nil {
			panic(err)
		}
	}

	p.status = falcon.APP_STATUS_RUNNING
	return err
}

func (p *Service) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = falcon.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, falcon.GetType(modules[n-i-1]))
		if err = modules[n-i-1].stop(p); err != nil {
			//panic(err)
			glog.Error(err)
		}
	}
	return err
}

func (p *Service) Reload(c interface{}) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Service)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].reload(p); err != nil {
			glog.Error(err)
		}
	}
	return err
}

func (p *Service) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
