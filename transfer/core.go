/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/transfer/config"
)

const (
	MODULE_NAME     = "\x1B[35m[TRANSFER]\x1B[0m"
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
	CTRL_STEP       = 360
	C_API_ADDR      = "apiaddr"
	C_HTTP_ADDR     = "httpaddr"
	C_CALL_TIMEOUT  = "calltimeout"
	C_BURST_SIZE    = "burstsize"
)

var (
	modules     []module
	ConfDefault = map[string]string{
		C_CALL_TIMEOUT: "5000",
		C_BURST_SIZE:   "16",
	}
)

type module interface {
	prestart(*Transfer) error // alloc public data
	start(*Transfer) error    // alloc private data, run private goroutine
	stop(*Transfer) error     // free private data, private goroutine exit
	reload(*Transfer) error   // try to keep the data, refresh configure
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type reqPayload struct {
	action  int
	shardId int
	data    interface{}
	done    chan interface{}
}

type Transfer struct {
	Conf    *config.Transfer
	oldConf *config.Transfer
	// runtime
	status uint32
	//appPutChan chan *service.Item // upstreams
	shardmap []chan *reqPayload
}

func (p *Transfer) New(conf interface{}) falcon.Module {
	return &Transfer{
		Conf: conf.(*config.Transfer),
	}
}

func (p *Transfer) Name() string {
	return p.Conf.Name
}

func (p *Transfer) Parse(text []byte, filename string, lino int) falcon.ModuleConf {
	p.Conf = config.Parse(text, filename, lino).(*config.Transfer)
	p.Conf.Configer.Set(falcon.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Transfer) String() string {
	return p.Conf.String()
}

func (p *Transfer) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = falcon.APP_STATUS_INIT
	p.shardmap = make([]chan *reqPayload, falcon.SHARD_NUM)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].prestart(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Transfer) Start() (err error) {
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

func (p *Transfer) Stop() (err error) {
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

func (p *Transfer) Reload(c interface{}) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	p.oldConf = p.Conf
	p.Conf = c.(*config.Transfer)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err

}

func (p *Transfer) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}

type BaseModule struct{}

func (p *BaseModule) prestart(transfer *Transfer) error { return nil }
func (p *BaseModule) start(transfer *Transfer) error    { return nil }
func (p *BaseModule) stop(transfer *Transfer) error     { return nil }
func (p *BaseModule) reload(transfer *Transfer) error   { return nil }
