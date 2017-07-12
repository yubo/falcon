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
	"github.com/yubo/falcon/utils"
)

const (
	MODULE_NAME     = "\x1B[32m[AGENT]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	modules []module
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
	appUpdateChan chan *[]*utils.MetaData
}

func (p *Agent) New(conf interface{}) falcon.Module {
	return &Agent{
		Conf:          conf.(*config.ConfAgent),
		appUpdateChan: make(chan *[]*utils.MetaData, 16),
	}
}

func (p *Agent) Name() string {
	return p.Conf.Name
}

func (p *Agent) Parse(text []byte, filename string, lino int, debug bool) interface{} {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfAgent)
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
	for i, n := 0, len(modules); i < n; i++ {
		if e := modules[n-i].stop(p); e != nil {
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
