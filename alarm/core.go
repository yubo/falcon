/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm/config"
	"github.com/yubo/falcon/alarm/parse"
	fconfig "github.com/yubo/falcon/config"
)

const (
	MODULE_NAME     = "\x1B[31m[ALARM]\x1B[0m"
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
	CTRL_STEP       = 360
	EVENT_TIMEOUT   = 180

	C_API_ADDR         = "apiaddr"
	C_HTTP_ADDR        = "httpaddr"
	C_CALL_TIMEOUT     = "calltimeout"
	C_BURST_SIZE       = "burstsize"
	C_DB_MAX_IDLE      = "dbmaxidle"
	C_DB_MAX_CONN      = "dbmaxconn"
	C_SYNC_INTERVAL    = "syncinterval"
	C_SYNC_DSN         = "syncdsn"
	C_WORKER_PROCESSES = "worker_processes"
)

const (
	ACTION_F_EMAIL = 1 << iota
	ACTION_F_SMS
	ACTION_F_SCRIPT
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CALL_TIMEOUT:     "5000",
		C_BURST_SIZE:       "16",
		C_DB_MAX_IDLE:      "4",
		C_DB_MAX_CONN:      "4",
		C_WORKER_PROCESSES: "4",
		C_SYNC_INTERVAL:    "600",
	}
)

type module interface {
	prestart(*Alarm) error // alloc public data
	start(*Alarm) error    // alloc private data, run private goroutine
	stop(*Alarm) error     // free private data, private goroutine exit
	reload(*Alarm) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type Alarm struct {
	Conf    *config.Alarm
	oldConf *config.Alarm
	// runtime
	status            uint32
	putEventChan      chan *Event      // api put
	actionChan        chan *Action     // upstreams
	delEventEntryChan chan *eventEntry // time out entry
	lru               queue            //lru queue

}

func (p *Alarm) New(conf interface{}) falcon.Module {
	return &Alarm{
		Conf:              conf.(*config.Alarm),
		putEventChan:      make(chan *Event, 1024),
		actionChan:        make(chan *Action, 144),
		delEventEntryChan: make(chan *eventEntry, 144),
	}
}

func (p *Alarm) Name() string {
	return p.Conf.Name
}

func (p *Alarm) Parse(text []byte, filename string, lino int) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino).(*config.Alarm)
	p.Conf.Configer.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Alarm) String() string {
	return p.Conf.String()
}

func (p *Alarm) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = falcon.APP_STATUS_INIT
	p.lru.init()

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].prestart(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Alarm) Start() (err error) {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = falcon.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].start(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	p.status = falcon.APP_STATUS_RUNNING
	return err
}

func (p *Alarm) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = falcon.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, falcon.GetType(modules[n-i-1]))
		if e := modules[n-i-1].stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	return err
}

func (p *Alarm) Reload(c interface{}) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Alarm)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err

}

func (p *Alarm) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
