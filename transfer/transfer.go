/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/transfer/config"
	"github.com/yubo/falcon/transfer/parse"
)

const (
	MODULE_NAME        = "\x1B[32m[TRANSFER]\x1B[0m "
	CONN_RETRY         = 2
	DEBUG_STAT_STEP    = 60
	CTRL_STEP          = 360
	C_HTTP_ENABLE      = "http_enable"
	C_HTTP_ADDR        = "httpaddr"
	C_RPC_ENABLE       = "rpc_enable"
	C_RPC_ADDR         = "rpcaddr"
	C_WORKER_PROCESSES = "workerprocesses"
	C_CONN_TIMEOUT     = "conntimeout"
	C_CALL_TIMEOUT     = "calltimeout"
	C_PAYLOADSIZE      = "payloadsize"
	C_GRPC_ENABLE      = "grpc_enable"
	C_GRPC_ADDR        = "grpcaddr"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CONN_TIMEOUT:     "1000",
		C_CALL_TIMEOUT:     "5000",
		C_WORKER_PROCESSES: "2",
		C_HTTP_ENABLE:      "true",
		C_HTTP_ADDR:        "127.0.0.1:6060",
		C_RPC_ENABLE:       "true",
		C_RPC_ADDR:         "127.0.0.1:8433",
		C_GRPC_ENABLE:      "true",
		C_GRPC_ADDR:        "127.0.0.1:8434",
		C_PAYLOADSIZE:      "16",
	}
)

func init() {
	//falcon.RegisterModule(falcon.GetType(config.ConfTransfer{}), &Transfer{})
}

// module {{{
type module interface {
	prestart(*Transfer) error // alloc public data
	start(*Transfer) error    // alloc private data, run private goroutine
	stop(*Transfer) error     // free private data, private goroutine exit
	reload(*Transfer) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)

}

// }}}

// Transfer {{{
type Transfer struct {
	Conf    *config.ConfTransfer
	oldConf *config.ConfTransfer
	// runtime
	status        uint32
	appUpdateChan chan *[]*falcon.MetaData // upstreams
}

func (p *Transfer) New(conf interface{}) falcon.Module {
	return &Transfer{
		Conf:          conf.(*config.ConfTransfer),
		appUpdateChan: make(chan *[]*falcon.MetaData, 16),
	}
}

func (p *Transfer) Name() string {
	return p.Conf.Name
}

func (p *Transfer) Parse(text []byte, filename string, lino int, debug bool) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfTransfer)
	p.Conf.Configer.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Transfer) String() string {
	return p.Conf.String()
}

func (p *Transfer) Prestart() (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Prestart()", p.Conf.Name)
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		if e := modules[i].prestart(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Transfer) Start() (err error) {
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

func (p *Transfer) Stop() (err error) {
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

func (p *Transfer) Reload(c interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = c.(*config.ConfTransfer)

	for i := 0; i < len(modules); i++ {
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err

}

func (p *Transfer) Signal(sig os.Signal) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return err
}

// }}}
