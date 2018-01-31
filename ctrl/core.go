/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"encoding/json"
	"os"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/config"
)

const (
	MODULE_NAME = "\x1B[33m[CTRL]\x1B[0m"

	C_MASTER_MODE             = "mastermode"
	C_MI_MODE                 = "mimode"
	C_DEV_MODE                = "devmode"
	C_SESSION_GC_MAX_LIFETIME = "sessiongcmaxlifetime"
	C_SESSION_COOKIE_LIFETIME = "sessioncookielifetime"
	C_AUTH_MODULE             = "authmodule"
	C_CACHE_MODULE            = "cachemodule"
	C_LDAP_ADDR               = "ldapaddr"
	C_LDAP_BASE_DN            = "ldapbasedn"
	C_LDAP_BIND_DN            = "ldapbinddn"
	C_LDAP_BIND_PWD           = "ldapbindpwd"
	C_LDAP_FILTER             = "ldapfilter"
	C_MISSO_REDIRECT_URL      = "missoredirecturl"
	C_GITHUB_CLIENT_ID        = "githubclientid"
	C_GITHUB_CLIENT_SECRET    = "githubclientsecret"
	C_GITHUB_REDIRECT_URL     = "githubredirecturl"
	C_GOOGLE_CLIENT_ID        = "googleclientid"
	C_GOOGLE_CLIENT_SECRET    = "googleclientsecret"
	C_GOOGLE_REDIRECT_URL     = "googleredirecturl"
	C_HTTP_ADDR               = "httpaddr"
	C_DB_SCHEMA               = "dbschema"
	C_DB_MAX_CONN             = "dbmaxconn"
	C_DB_MAX_IDLE             = "dbmaxidle"
	C_ADMIN                   = "admin"
	C_DSN                     = "dsn"
	C_IDX_DSN                 = "idxdsn"
	C_ALARM_DSN               = "alarmdsn"
	C_TAG_SCHEMA              = "tagschema"
	C_WEIXIN_APP_ID           = "wxappid"
	C_WEIXIN_APP_SECRET       = "wxappsecret"
	C_TRANSFER_ADDR           = "transferaddr"
	C_CALL_TIMEOUT            = "calltimeout"
)

var (
	modules   []module
	Configure *config.Ctrl
	RunMode   uint32

	ConfDefault = map[string]string{
		C_MASTER_MODE:             "true",
		C_MI_MODE:                 "false",
		C_DEV_MODE:                "false",
		C_SESSION_GC_MAX_LIFETIME: "86400",
		C_SESSION_COOKIE_LIFETIME: "86400",
		C_AUTH_MODULE:             "ldap",
		C_CACHE_MODULE:            "host,role,system,tag,user",
		C_DB_MAX_CONN:             "30",
		C_DB_MAX_IDLE:             "30",
		C_CALL_TIMEOUT:            "5000",
	}

	ConfDesc = map[string]string{
		//ctrl.C_ENABLE_DOCS:             "ture/false",
		C_MASTER_MODE:             "bool",
		C_MI_MODE:                 "bool",
		C_DEV_MODE:                "bool",
		C_SESSION_GC_MAX_LIFETIME: "int",
		C_SESSION_COOKIE_LIFETIME: "int",
		C_AUTH_MODULE:             "ldap/misso/github/google",
		C_CACHE_MODULE:            "string",
		C_LDAP_ADDR:               "string",
		C_LDAP_BASE_DN:            "string",
		C_LDAP_BIND_DN:            "string",
		C_LDAP_BIND_PWD:           "string",
		C_LDAP_FILTER:             "string",
		C_MISSO_REDIRECT_URL:      "string",
		C_GITHUB_CLIENT_ID:        "string",
		C_GITHUB_CLIENT_SECRET:    "string",
		C_GITHUB_REDIRECT_URL:     "string",
		C_GOOGLE_CLIENT_ID:        "string",
		C_GOOGLE_CLIENT_SECRET:    "string",
		C_GOOGLE_REDIRECT_URL:     "string",
	}
)

// ctl runmode name
const (
	CTL_RUNMODE_MASTER = 1 << iota
	CTL_RUNMODE_DEV
	CTL_RUNMODE_MI
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

func RegisterModule(m module) {
	modules = append(modules, m)
}

type Ctrl struct {
	Conf    *config.Ctrl
	oldConf *config.Ctrl
	// runtime
	status uint32
	// rpcListener  *net.TCPListener
	// httpListener *net.TCPListener
	// httpMux      *http.ServeMux
}

func (p *Ctrl) New(conf interface{}) falcon.Module {
	return &Ctrl{
		Conf: conf.(*config.Ctrl),
	}
}

func (p *Ctrl) Name() string {
	return p.Conf.Name
}

func (p *Ctrl) Parse(text []byte, filename string, lino int) falcon.ModuleConf {
	p.Conf = config.Parse(text, filename, lino).(*config.Ctrl)
	/* TODO: fill agent, transfer, backend, graph */
	p.Conf.Ctrl.Set(falcon.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Ctrl) String() string {
	return p.Conf.String()
}

// ugly hack
// should called by main package
func (p *Ctrl) Prestart() (err error) {
	glog.V(3).Infof("%s Prestart()", MODULE_NAME)
	Configure = p.Conf
	p.status = falcon.APP_STATUS_INIT

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.prestart()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].PreStart(p); e != nil {
			err = e
			glog.Error(err)
		}
	}
	return err
}

func (p *Ctrl) Start() (err error) {
	glog.V(3).Infof("%s Start()", MODULE_NAME)
	p.status = falcon.APP_STATUS_PENDING

	conf := &p.Conf.Ctrl
	// ctrl config
	conf.Set(falcon.APP_CONF_DEFAULT, ConfDefault)

	// connect db, can not register db twice  :(
	// get ctrl config from db
	if conf.DefaultBool(C_MASTER_MODE, false) {
		RunMode |= CTL_RUNMODE_MASTER
	}
	if conf.DefaultBool(C_MI_MODE, false) {
		RunMode |= CTL_RUNMODE_MI
	}
	if conf.DefaultBool(C_DEV_MODE, false) {
		RunMode |= CTL_RUNMODE_DEV
	}

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.start()", MODULE_NAME, falcon.GetType(modules[i]))
		if e := modules[i].Start(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	p.status = falcon.APP_STATUS_RUNNING
	return err
}

func (p *Ctrl) Stop() (err error) {
	glog.V(3).Infof("%s Stop()", MODULE_NAME)
	p.status = falcon.APP_STATUS_EXIT

	for n, i := len(modules), 0; i < n; i++ {
		glog.V(4).Infof("%s %s.stop()", MODULE_NAME, falcon.GetType(modules[n-i-1]))
		if e := modules[n-i-1].Stop(p); e != nil {
			err = e
			glog.Error(err)
		}
	}

	return nil
}

// TODO: reload is not yet implemented
func (p *Ctrl) Reload(c interface{}) (err error) {

	return nil

	glog.V(3).Infof("%s Reload()", MODULE_NAME)
	p.Conf = c.(*config.Ctrl)

	for i := 0; i < len(modules); i++ {
		glog.V(4).Infof("%s %s.reload()", MODULE_NAME, falcon.GetType(modules[i]))
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
