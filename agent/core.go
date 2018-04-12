/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
)

const (
	//MODULE_NAME = "\x1B[32m[AGENT]\x1B[0m"
	MODULE_NAME      = "agent"
	PUT_REQUEST_SIZE = 144
)

var (
	modules []module
)

type module interface {
	prestart(*Agent) error // alloc public data
	start(*Agent) error    // alloc private data, run private goroutine
	stop(*Agent) error     // free private data, private goroutine exit
	reload(*Agent) error   // try to keep the data, refresh configure
}

func init() {
	RegisterModule(&CollectModule{})
	RegisterModule(&ClientModule{})
	RegisterModule(&ApiModule{})
	RegisterModule(&ApiGwModule{})
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type AgentConfig struct {
	Disable         bool                   `json:"disable"`
	Debug           bool                   `json:"debug"`
	ApiAddr         string                 `json:"api_addr"`
	HttpAddr        string                 `json:"http_addr"`
	Burstsize       int                    `json:"burst_size"`
	WorkerProcesses int                    `json:"worker_processes"`
	CallTimeout     int                    `json:"call_timeout"`
	Host            string                 `json:"hostname"`
	Interval        int                    `json:"interval"`
	IfacePrefix     []string               `json:"iface_prefix"`
	Upstream        string                 `json:"upstream"`
	EmuTplDir       string                 `json:"emu_tpl_dir"`
	Plugins         []string               `json:"plugins"`
	IpcEnable       bool                   `json:"ipc_enable"`
	PluginsConfig   map[string]interface{} `json:"plugins_config"`
	EtcdClient      *core.EtcdCliConfig    `json:"etcd_client"`
}

type Agent struct {
	Conf    *AgentConfig
	oldConf *AgentConfig
	// runtime
	status  uint32
	PutChan chan *PutRequest
}

func (p *Agent) ReadConfig(conf *core.Configer) error {
	p.Conf = &AgentConfig{Disable: true}

	err := conf.Read(MODULE_NAME, p.Conf)
	if err == core.ErrNoExits {
		return nil
	}

	return err
}

func (p *Agent) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	p.status = core.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].prestart(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Agent) Start() (err error) {
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

func (p *Agent) Stop() (err error) {
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

func (p *Agent) Reload(conf *core.Configer) (err error) {
	glog.V(3).Infof("%s Reload()", MODULE_NAME)

	newConfig := &AgentConfig{}
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

func (p *Agent) Signal(sig os.Signal) (err error) {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return err
}
