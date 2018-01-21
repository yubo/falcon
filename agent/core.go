/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent/config"
)

const (
	MODULE_NAME     = "\x1B[32m[AGENT]\x1B[0m"
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60

	C_UPSTREAM     = "upstream"
	C_CALL_TIMEOUT = "calltimeout"
	C_API_ADDR     = "apiaddr"
	C_HTTP_ADDR    = "httpaddr"
	C_INTERVAL     = "interval"
	C_BURST_SIZE   = "burstsize"
	C_IFACE_PREFIX = "ifaceprefix"
	C_PLUGINS      = "plugins"
	C_EMU_ENABLE   = "emuenable"
	C_EMU_TPL_DIR  = "emutpldir"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CALL_TIMEOUT: "5000",
		C_INTERVAL:     "60",
		C_BURST_SIZE:   "16",
		C_IFACE_PREFIX: "eth,em",
	}
)

type module interface {
	prestart(*Agent) error // alloc public data
	start(*Agent) error    // alloc private data, run private goroutine
	stop(*Agent) error     // free private data, private goroutine exit
	reload(*Agent) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type Agent struct {
	Conf    *config.Agent
	oldConf *config.Agent
	// runtime
	status  uint32
	putChan chan *putContext
}

func (p *Agent) New(conf interface{}) falcon.Module {
	return &Agent{
		Conf:    conf.(*config.Agent),
		putChan: make(chan *putContext, 144),
	}
}

func (p *Agent) Name() string {
	return p.Conf.Name
}

func (p *Agent) Parse(text []byte, filename string, lino int) falcon.ModuleConf {
	p.Conf = config.Parse(text, filename, lino).(*config.Agent)
	p.Conf.Configer.Set(falcon.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Agent) String() string {
	return p.Conf.String()
}

func (p *Agent) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Agent) Start() (err error) {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = falcon.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].start(p); err != nil {
			glog.Fatal(err)
		}
	}
	p.status = falcon.APP_STATUS_RUNNING

	return err
}

func (p *Agent) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = falcon.APP_STATUS_EXIT
	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, falcon.GetType(modules[n-i-1]))
		if err = modules[n-i-1].stop(p); err != nil {
			glog.Fatal(err)
		}
	}

	return err
}

func (p *Agent) Reload(c interface{}) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Agent)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, falcon.GetType(modules[i]))
		if err = modules[i].reload(p); err != nil {
			glog.Error(err)
		}
	}

	return err
}

func (p *Agent) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
