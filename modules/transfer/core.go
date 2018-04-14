/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/service"
)

const (
	//MODULE_NAME = "\x1B[35m[TRANSFER]\x1B[0m"
	MODULE_NAME = "transfer"
)

var (
	modules []module
)

type module interface {
	prestart(*Transfer) error // alloc public data
	start(*Transfer) error    // alloc private data, run private goroutine
	stop(*Transfer) error     // free private data, private goroutine exit
	reload(*Transfer) error   // try to keep the data, refresh configure
}

func init() {
	registerModule(&clientModule{})
	registerModule(&apiModule{})
	registerModule(&apiGwModule{})
}

func registerModule(m module) {
	modules = append(modules, m)
}

type TransferConfig struct {
	Disable         bool                `json:"disable"`
	Debug           bool                `json:"debug"`
	ApiAddr         string              `json:"api_addr"`
	HttpAddr        string              `json:"http_addr"`
	WorkerProcesses int                 `json:"worker_processes"`
	CallTimeout     int                 `json:"call_timeout"`
	ShardNum        int                 `json:"shard_num"`
	ShardMap        []string            `json:"shard_map"`
	BurstSize       int                 `json:"burst_size"`
	EtcdClient      *core.EtcdCliConfig `json:"etcd_client"`
}

type reqPayload struct {
	action  int
	shardId int
	data    interface{}
	done    chan interface{}
}

type Transfer struct {
	Conf    *TransferConfig
	oldConf *TransferConfig
	// runtime
	status uint32
	//appPutChan chan *service.Item // upstreams
	shardmap []chan *reqPayload
}

func (p *Transfer) ReadConfig(conf *core.Configer) error {
	p.Conf = &TransferConfig{Disable: true}

	err := conf.Read(MODULE_NAME, p.Conf)
	if err == core.ErrNoExits {
		return nil
	}

	return err
}

func (p *Transfer) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = core.APP_STATUS_INIT

	p.shardmap = make([]chan *reqPayload, service.SHARD_NUM)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Transfer) Start() (err error) {
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

func (p *Transfer) Stop() (err error) {
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

func (p *Transfer) Reload(conf *core.Configer) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	newConfig := &TransferConfig{}
	err = conf.Read(MODULE_NAME, newConfig)
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

func (p *Transfer) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
