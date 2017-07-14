/*
 * Copyright 2016 falcon Author. All rights reserved.
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

	C_UPSTREAM         = "upstream"
	C_CONN_TIMEOUT     = "conntimeout"
	C_CALL_TIMEOUT     = "calltimeout"
	C_WORKER_PROCESSES = "workerprocesses"
	C_HTTP_ADDR        = "httpaddr"
	C_RPC_ADDR         = "rpcaddr"
	C_GRPC_ADDR        = "grpcaddr"
	C_INTERVAL         = "interval"
	C_PAYLOADSIZE      = "payloadsize"
	C_IFACE_PREFIX     = "ifaceprefix"
	C_HTTP_ENABLE      = "http_enable"
	C_GRPC_ENABLE      = "grpc_enable"
	C_RPC_ENABLE       = "rpc_enable"
	C_PLUGINS          = "plugins"
	C_EMU_ENABLE       = "emuenable"
	C_EMU_HOST         = "emuhost"
	C_EMU_HOSTNUM      = "emuhostnum"
	C_EMU_METRIC       = "emumetric"
	C_EMU_METRICNUM    = "emumetricnum"
	C_EMU_TPL          = "tpl"
	C_EMU_TPLNUM       = "tplnum"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CONN_TIMEOUT:     "1000",
		C_CALL_TIMEOUT:     "5000",
		C_WORKER_PROCESSES: "2",
		C_HTTP_ENABLE:      "true",
		C_HTTP_ADDR:        "127.0.0.1:1988",
		C_RPC_ENABLE:       "true",
		C_RPC_ADDR:         "127.0.0.1:1989",
		C_GRPC_ENABLE:      "true",
		C_GRPC_ADDR:        "127.0.0.1:1990",
		C_INTERVAL:         "60",
		C_PAYLOADSIZE:      "16",
		C_IFACE_PREFIX:     "eth,em",
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
	modules = append(modules, m)
}

// }}}

// Agent {{{
type Agent struct {
	Conf    *config.ConfAgent
	oldConf *config.ConfAgent
	// runtime
	status        uint32
	appUpdateChan chan *[]*falcon.MetaData
}

func (p *Agent) New(conf interface{}) falcon.Module {
	return &Agent{
		Conf:          conf.(*config.ConfAgent),
		appUpdateChan: make(chan *[]*falcon.MetaData, 16),
	}
}

func (p *Agent) Name() string {
	return p.Conf.Name
}

func (p *Agent) Parse(text []byte, filename string, lino int, debug bool) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfAgent)
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
		if e := modules[i].prestart(p); e != nil {
			//panic(err)
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Agent) Start() (err error) {
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

func (p *Agent) Stop() (err error) {
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

func (p *Agent) Reload(c interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = c.(*config.ConfAgent)

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
