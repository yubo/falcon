/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
)

const (
	//MODULE_NAME          = "\x1B[31m[ALARM]\x1B[0m"
	MODULE_NAME          = "alarm"
	EVENT_CLEAN_INTERVAL = 2 // second
)

const (
	ACTION_F_EMAIL = 1 << iota
	ACTION_F_SMS
	ACTION_F_SCRIPT
)

var (
	modules []module
)

type module interface {
	prestart(*Alarm) error // alloc public data
	start(*Alarm) error    // alloc private data, run private goroutine
	stop(*Alarm) error     // free private data, private goroutine exit
	reload(*Alarm) error   // try to keep the data, refresh configure
}

func init() {
	RegisterModule(&ApiModule{})
	RegisterModule(&ApiGwModule{})
	RegisterModule(&ClientModule{})
	RegisterModule(&TriggerModule{})
	RegisterModule(&TaskModule{})
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type AlarmConfig struct {
	Disable         bool   `json:"disable"`
	Debug           bool   `json:"debug"`
	ApiAddr         string `json:"api_addr"`
	HttpAddr        string `json:"http_addr"`
	Burstsize       int    `json:"burst_size"`
	DbMaxIdle       int    `json:"db_max_idle"`
	DbMaxConn       int    `json:"db_max_conn"`
	WorkerProcesses int    `json:"worker_processes"`
	CallTimeout     int    `json:"call_timeout"`
	//Upstream        []string            `json:"upstream"`
	SyncDsn         string              `json:"sync_dsn"`
	SyncInterval    int                 `json:"sync_interval"`
	EventExpireTime int                 `json:"event_expire_time"`
	EtcdClient      *core.EtcdCliConfig `json:"etcd_client"`
}

type Alarm struct {
	Conf    *AlarmConfig
	oldConf *AlarmConfig
	// runtime
	status            uint32
	putEventChan      chan *Event      // api put
	actionChan        chan *Action     // upstreams
	delEventEntryChan chan *eventEntry // time out entry
	lru               queue            //lru queue

}

func (p *Alarm) ReadConfig(conf *core.Configer) error {
	p.Conf = &AlarmConfig{Disable: true}

	err := conf.Read(MODULE_NAME, p.Conf)
	if err == core.ErrNoExits {
		return nil
	}

	return err
}

func (p *Alarm) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = core.APP_STATUS_INIT
	p.lru.init()

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Alarm) Start() (err error) {
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

func (p *Alarm) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = core.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, core.GetType(modules[n-i-1]))
		if err = modules[n-i-1].stop(p); err != nil {
			glog.Fatal(err)
		}
	}

	return err
}

func (p *Alarm) Reload(conf *core.Configer) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	newConfig := &AlarmConfig{}
	err = conf.Read(MODULE_NAME, &newConfig)
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

func (p *Alarm) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
