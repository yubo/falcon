/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/stats"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/transfer"
)

const (
	MODULE_NAME = "ctrl"
)

// ctl runmode name
const (
	CTL_RUNMODE_MASTER = 1 << iota
	CTL_RUNMODE_DEV
	CTL_RUNMODE_BEEGODEV
	CTL_RUNMODE_MI
)

var (
	modules []module
	//Configure *config.Ctrl
)

type Kv struct {
	Key     string
	Section string
	Value   string
}

type module interface {
	PreStart(*Ctrl) error // alloc public data
	Start(*Ctrl) error    // alloc private data, run private goroutine
	Stop(*Ctrl) error     // free private data, private goroutine exit
	Reload(*Ctrl) error   // try to keep the data, refresh configure
}

func init() {
	RegisterModule(&dbModule{})
	RegisterModule(&etcdCliModule{})
	RegisterModule(&clientModule{})
	RegisterModule(&modelsModule{})
	RegisterModule(&rateLimitsModule{})
	RegisterModule(&apiModule{})
}

func RegisterModule(m module) {
	modules = append(modules, m)
}

type CtrlConfig struct {
	Configer              *core.Configer      `json:"-"`
	Disable               bool                `json:"disable"`
	EtcdClient            *core.EtcdCliConfig `json:"etcd_client"`
	MasterMode            bool                `json:"master_mode"`
	MiMode                bool                `json:"mi_mode"`
	DevMode               bool                `json:"dev_mode"`
	BeeMode               bool                `json:"bee_mode"`
	BeegoDevMode          bool                `json:"beego_dev_mode"`
	Debug                 int                 `json:"debug"`
	Dsn                   string              `json:"dsn"`
	IdxDsn                string              `json:"idx_dsn"`
	AlarmDsn              string              `json:"alarm_dsn"`
	EtcdEndpoints         string              `json:"etcd_endpoints"`
	TransferAddr          string              `json:"transfer_addr"`
	HttpAddr              string              `json:"http_addr"`
	SessionGcMaxLifetime  int                 `json:"session_gc_max_lifetime"`
	SessionCookieLifetime int                 `json:"session_cookie_lifetime"`
	CallTimeout           int                 `json:"call_timeout"`
	DbMaxIdle             int                 `json:"db_max_idle"`
	DbMaxConn             int                 `json:"db_max_conn"`
	DbSchema              string              `json:"db_schema"`
	EnableDocs            bool                `json:"enable_docs"`
	PluginAlarm           bool                `json:"plugin_alarm"`
	MiNornsUrl            string              `json:"mi_norns_url"`
	MiNornsInterval       int                 `json:"mi_norns_interval"`
	WxAppId               string              `json:"wx_app_id"`
	WxAppSecret           string              `json:"wx_app_secret"`
	HttpRateLimit         struct {
		Enable     bool `json:"enable"`
		Limit      int  `json:"limit"`
		Accuracy   int  `json:"accuracy"`
		GcTimeout  int  `json:"gc_timeout"`
		GcInterval int  `json:"gc_interval"`
	} `json:"http_rate_limit"`
	/*
		Auth       struct {
			Ldap struct {
				Addr    string `json:"addr"`
				baseDn  string `json:"base_dn"`
				bindDn  string `json:"bind_dn"`
				bindPwd string `json:"bind_pwd"`
				filter  string `json:"filter"`
			} `json:"ldap"`
			Misso struct {
				RedirectUrl string `json:"redirect_url"`
			} `json:"misso"`
			Google struct {
				ClientId     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				RedirectUrl  string `json:"redirect_url"`
			} `json:"google"`
			Github struct {
				ClientId     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				RedirectUrl  string `json:"redirect_url"`
			} `json:"github"`
			Fw struct {
				ClientId     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
				RedirectUrl  string `json:"redirect_url"`
			} `json:"fw"`
		} `json:"auth"`
	*/
}

type Ctrl struct {
	// TODO Conf -> conf
	Conf    *CtrlConfig
	oldConf *CtrlConfig
	status  uint32 // runtime

	// for moduels
	db          *models.CtrlDb // orm module
	transferCli transfer.TransferClient
	etcdCli     *core.EtcdCli
	// rpcListener  *net.TCPListener
	// httpListener *net.TCPListener
	// httpMux      *http.ServeMux
}

func (p *Ctrl) ReadConfig(conf *core.Configer) error {
	p.Conf = &CtrlConfig{Disable: true}

	err := conf.Read(MODULE_NAME, p.Conf)
	if err == core.ErrNoExits {
		return nil
	}
	p.Conf.Configer = core.ToConfiger(conf.GetRaw(MODULE_NAME))

	return err
}

// ugly hack
// should called by main package
func (p *Ctrl) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)

	//Configure = p.Conf
	p.status = core.APP_STATUS_INIT
	//core.Sort(modules)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, core.GetType(modules[i]))
		if err = modules[i].PreStart(p); err != nil {
			glog.Fatal(err)
		}
	}
	return err
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = core.APP_STATUS_PENDING

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, core.GetType(modules[i]))
		if err := modules[i].Start(p); err != nil {
			glog.Fatal(err)
		}
	}

	p.status = core.APP_STATUS_RUNNING
	return nil
}

func (p *Ctrl) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = core.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, core.GetType(modules[n-i-1]))
		if err = modules[n-i-1].Stop(p); err != nil {
			glog.Fatal(err)
		}
	}

	return nil
}

// TODO: reload is not yet implemented
func (p *Ctrl) Reload(conf *core.Configer) (err error) {

	return nil

	glog.V(3).Infof("%s Reload()", MODULE_NAME)
	newConfig := &CtrlConfig{}
	err = conf.Read(MODULE_NAME, newConfig)
	if err != nil {
		return err
	}
	newConfig.Configer = core.ToConfiger(conf.GetRaw(MODULE_NAME))

	p.oldConf = p.Conf
	p.Conf = newConfig

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, core.GetType(modules[i]))
		if e := modules[i].Reload(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.Infof("%s recv signal %#v", MODULE_NAME, sig)
	return nil
}

func (p *Ctrl) Stats(conf *core.Configer) (s string, err error) {
	// http api
	var counter []uint64

	url := conf.GetStr(fmt.Sprintf("%s.http_addr", MODULE_NAME))
	if strings.HasPrefix(url, ":") {
		url = "http://localhost" + url
	}
	url += "/v1.0/pub/stats"

	if err := core.GetJson(url, &counter, time.Duration(200)*time.Millisecond); err != nil {
		return "", err
	}

	for i := 0; i < stats.ST_ARRAY_SIZE; i++ {
		s += fmt.Sprintf("%-30s %d\n",
			stats.CounterName[i], counter[i])
	}
	return s, nil
}

func GetDbConfig(o orm.Ormer, module string) (ret map[string]string, err error) {
	var row Kv
	ret = make(map[string]string)

	err = o.Raw("SELECT `section`, `key`, `value` FROM `kv` where "+
		"`section` = ? and `key` = 'config'", module).QueryRow(&row)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(row.Value), &ret)
	return
}
