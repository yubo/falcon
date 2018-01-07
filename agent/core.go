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
	"github.com/yubo/falcon/agent/parse"
	fconfig "github.com/yubo/falcon/config"
)

const (
	MODULE_NAME     = "\x1B[32m[AGENT]\x1B[0m "
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

// module {{{

type module interface {
	prestart(*Agent) error // alloc public data
	start(*Agent) error    // alloc private data, run private goroutine
	stop(*Agent) error     // free private data, private goroutine exit
	reload(*Agent) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	glog.Infof("%s RegisterModule %s", MODULE_NAME, falcon.GetType(m))
	modules = append(modules, m)
}

// }}}

// Agent {{{
type Agent struct {
	Conf    *config.Agent
	oldConf *config.Agent
	// runtime
	status     uint32
	appPutChan chan []*falcon.Item
}

func (p *Agent) New(conf interface{}) falcon.Module {
	return &Agent{
		Conf:       conf.(*config.Agent),
		appPutChan: make(chan []*falcon.Item, 64),
	}
}

func (p *Agent) Name() string {
	return p.Conf.Name
}

func (p *Agent) Parse(text []byte, filename string, lino int) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino).(*config.Agent)
	p.Conf.Configer.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Agent) String() string {
	return p.Conf.String()
}

func (p *Agent) Prestart() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Prestart()", p.Conf.Name)
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		if err = modules[i].prestart(p); err != nil {
			panic(err)
			//glog.Error(err)
		}
	}
	return err
}

func (p *Agent) Start() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Name)
	p.status = falcon.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		if e := modules[i].start(p); e != nil {
			panic(err)
			err = e
			glog.Error(err)
		}
	}
	p.status = falcon.APP_STATUS_RUNNING

	return err
}

func (p *Agent) Stop() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT
	for n, i := len(modules), 0; i < n; i++ {
		if e := modules[n-i-1].stop(p); e != nil {
			panic(err)
			err = e
			glog.Error(err)
		}
	}

	return err
}

func (p *Agent) Reload(c interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Agent)

	for i := 0; i < len(modules); i++ {
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	return err
}

func (p *Agent) Signal(sig os.Signal) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return err
}

// }}}
