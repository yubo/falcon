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
	"github.com/yubo/falcon/alarm"
	"github.com/yubo/falcon/lib/core"
)

const (
	CACHE_DATA_SIZE       = 1 << 4
	CACHE_DATA_SIZE_MASK  = CACHE_DATA_SIZE - 1
	INDEX_QPS             = 100
	INDEX_UPDATE_INTERVAL = 3600 * 24
	CACHE_EXPIRE_TIME     = 3600 * 26 // 超时时间 (s)
	CACHE_CLEAN_INTERVAL  = 3600      // 回收检查间隔时间 (s)
	MODULE_NAME           = "service"
	//MODULE_NAME           = "\x1B[36m[SERVICE]\x1B[0m"
)

var (
	modules []module
)

type module interface {
	prestart(*Service) error // alloc public data
	start(*Service) error    // alloc private data, run private goroutine
	stop(*Service) error     // free private data, private goroutine exit
	reload(*Service) error   // try to keep the data, refresh configure
}

func init() {
	registerModule(&timerModule{})
	registerModule(&tsdbModule{})
	registerModule(&cacheModule{})
	registerModule(&triggerModule{})
	registerModule(&clientModule{})
	registerModule(&apiModule{})
	registerModule(&apiGwModule{})
}

func registerModule(m module) {
	modules = append(modules, m)
}

type ServiceConfig struct {
	Disable        bool                `json:"disable"`
	Debug          bool                `json:"debug"`
	ApiAddr        string              `json:"api_addr"`
	HttpAddr       string              `json:"http_addr"`
	AlarmAddr      string              `json:"alarm_addr"`
	Dsn            string              `json:"dsn"`
	IdxDsn         string              `json:"index_dsn"`
	DbMaxIdle      int                 `json:"db_max_idle"`
	DbMaxConn      int                 `json:"db_max_conn"`
	CallTimeout    int                 `json:"call_timeout"`
	ConfInterval   int                 `json:"conf_interval"`
	JudgeInterval  int                 `json:"judge_interval"`
	JudgeNum       int                 `json:"judge_num"`
	ShardIds       []int               `json:"shard_ids"`
	CacheTimeout   int                 `json:"cache_timeout"`
	RrdTimeout     int                 `json:"rrd_timeout"`
	TsdbBucketNum  int                 `json:"tsdb_bucket_num"`
	TsdbBucketSize int                 `json:"tsdb_bucket_size"`
	TsdbTimeout    int                 `json:"tsdb_timeout"`
	TsdbDir        string              `json:"tsdb_dir"`
	EtcdClient     *core.EtcdCliConfig `json:"etcd_client"`
}

type Service struct {
	Conf    *ServiceConfig
	oldConf *ServiceConfig
	// runtime
	status uint32

	//shardModule
	cache *cacheModule

	// tsdb
	tsdb *tsdbModule

	// event_trigger
	eventChan chan *alarm.Event

	//storageModule
	//hdisk []string
	//storageNetTaskCh map[string]chan *netTask
	//storageIoTaskCh  []chan *ioTask
}

func (p *Service) ReadConfig(conf *core.Configer) error {
	p.Conf = &ServiceConfig{Disable: true}

	err := conf.Read(MODULE_NAME, p.Conf)
	if err == core.ErrNoExits {
		return nil
	}

	return err
}

func (p *Service) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = core.APP_STATUS_INIT

	p.eventChan = make(chan *alarm.Event, 1024)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			glog.Fatal(err)
			//glog.Error(err)
		}
	}
	return err
}

func (p *Service) Start() (err error) {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = core.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].start(p); err != nil {
			glog.Fatal(err)
		}
	}

	p.status = core.APP_STATUS_RUNNING
	return err
}

func (p *Service) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = core.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, core.GetType(modules[n-i-1]))
		if err = modules[n-i-1].stop(p); err != nil {
			//panic(err)
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Service) Reload(conf *core.Configer) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	newConfig := &ServiceConfig{}
	err = conf.Read(MODULE_NAME, newConfig)
	if err != nil {
		return err
	}

	p.oldConf = p.Conf
	p.Conf = newConfig

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].reload(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Service) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
