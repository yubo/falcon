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
)

const (
	MODULE_NAME     = "\x1B[32m[TRANSFER]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
	CTRL_STEP       = 360
)

var (
	modules []module
)

func init() {
	falcon.RegisterModule(falcon.GetType(falcon.ConfTransfer{}), &Transfer{})
	registerModule(&statsModule{})
	registerModule(&httpModule{})
	registerModule(&rpcModule{})
	registerModule(&backendModule{})
}

// module {{{
type module interface {
	prestart(*Transfer) error // alloc public data
	start(*Transfer) error    // alloc private data, run private goroutine
	stop(*Transfer) error     // free private data, private goroutine exit
	reload(*Transfer) error   // try to keep the data, refresh configure
}

func registerModule(m module) {
	modules = append(modules, m)
}

// }}}

// Transfer {{{
type Transfer struct {
	Conf    *falcon.ConfTransfer
	oldConf *falcon.ConfTransfer
	// runtime
	status        uint32
	appUpdateChan chan *[]*falcon.MetaData // upstreams
}

func (p *Transfer) New(conf interface{}) falcon.Module {
	return &Transfer{
		Conf:          conf.(*falcon.ConfTransfer),
		appUpdateChan: make(chan *[]*falcon.MetaData, 16),
	}
}

func (p *Transfer) Name() string {
	return p.Conf.Name
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

	for i, n := 0, len(modules); i < n; i++ {
		if e := modules[n-i].stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	return err
}

func (p *Transfer) Reload(config interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = config.(*falcon.ConfTransfer)

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
