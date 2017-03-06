/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package loadbalance

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

const (
	MODULE_NAME     = "\x1B[32m[LB]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
	CTRL_STEP       = 360
)

var (
	modules []module
)

func init() {
	falcon.RegisterModule(falcon.GetType(falcon.ConfLoadbalance{}), &Loadbalance{})
	registerModule(&statsModule{})
	registerModule(&httpModule{})
	registerModule(&rpcModule{})
	registerModule(&backendModule{})
}

// module {{{
type module interface {
	prestart(*Loadbalance) error // alloc public data
	start(*Loadbalance) error    // alloc private data, run private goroutine
	stop(*Loadbalance) error     // free private data, private goroutine exit
	reload(*Loadbalance) error   // try to keep the data, refresh configure
}

func registerModule(m module) {
	modules = append(modules, m)
}

// }}}

// Loadbalance {{{
type Loadbalance struct {
	Conf    *falcon.ConfLoadbalance
	oldConf *falcon.ConfLoadbalance
	// runtime
	status        uint32
	appUpdateChan chan *[]*falcon.MetaData // upstreams
}

func (p *Loadbalance) New(conf interface{}) falcon.Module {
	return &Loadbalance{
		Conf:          conf.(*falcon.ConfLoadbalance),
		appUpdateChan: make(chan *[]*falcon.MetaData, 16),
	}
}

func (p *Loadbalance) Name() string {
	return p.Conf.Name
}

func (p *Loadbalance) String() string {
	return p.Conf.String()
}

func (p *Loadbalance) Prestart() (err error) {
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

func (p *Loadbalance) Start() (err error) {
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

func (p *Loadbalance) Stop() (err error) {
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

func (p *Loadbalance) Reload(config interface{}) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	p.oldConf = p.Conf
	p.Conf = config.(*falcon.ConfLoadbalance)

	for i := 0; i < len(modules); i++ {
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err

}

func (p *Loadbalance) Signal(sig os.Signal) (err error) {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return err
}

// }}}
