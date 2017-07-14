/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/ctrl/config"
	"github.com/yubo/falcon/ctrl/parse"
)

type hookfunc func(conf *config.ConfCtrl) error

type Ctrl struct {
	Conf    *config.ConfCtrl
	oldConf *config.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
}

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60

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
	C_MI_NORNS_INTERVAL       = "minornsinterval"
	C_HTTP_ADDR               = "httpaddr"
	C_DB_MAX_CONN             = "dbmaxconn"
	C_DB_MAX_IDLE             = "dbmaxidle"
	C_MI_NORNS_URL            = "minornsurl"
	C_RL_GC_INTERVAL          = "rlgcinterval"
	C_RL_GC_TIMEOUT           = "rlgctimeout"
	C_RL_LIMIT                = "rllimit"
	C_RL_ACCURACY             = "rlaccuracy"
	C_ADMIN                   = "admin"
	C_DSN                     = "dsn"
	C_IDX_DSN                 = "idxdsn"
	C_ALARM_DSN               = "alarmdsn"
	C_TAG_SCHEMA              = "tagschema"
	C_TRANSFER_URL            = "transferurl"
)

var (
	prestartHooks = make([]hookfunc, 0)
	reloadHooks   = make([]hookfunc, 0)
	Configure     *config.ConfCtrl
	EtcdCli       *falcon.EtcdCli

	ConfDefault = map[string]string{
		//C_RUN_MODE:                "pub",
		C_MASTER_MODE:             "true",
		C_MI_MODE:                 "false",
		C_DEV_MODE:                "false",
		C_HTTP_ADDR:               "8001",
		C_SESSION_GC_MAX_LIFETIME: "86400",
		C_SESSION_COOKIE_LIFETIME: "86400",
		C_AUTH_MODULE:             "ldap",
		C_CACHE_MODULE:            "host,role,system,tag,user",
		C_DB_MAX_CONN:             "30",
		C_DB_MAX_IDLE:             "30",
		C_MI_NORNS_URL:            "http://norns.dev/api/v1/tagstring/cop.xiaomi/hostinfos",
		C_MI_NORNS_INTERVAL:       "5",
	}

	ConfDesc = map[string]string{
		//ctrl.C_RUN_MODE:                "dev/prod",
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

func RegisterPrestart(fn hookfunc) {
	prestartHooks = append(prestartHooks, fn)
}

func RegisterReload(fn hookfunc) {
	reloadHooks = append(reloadHooks, fn)
}

func (p *Ctrl) New(conf interface{}) falcon.Module {
	return &Ctrl{Conf: conf.(*config.ConfCtrl)}
}

func (p *Ctrl) Name() string {
	return fmt.Sprintf("%s", p.Conf.Name)
}

func (p *Ctrl) Parse(text []byte, filename string, lino int, debug bool) fconfig.ModuleConf {
	p.Conf = parse.Parse(text, filename, lino, debug).(*config.ConfCtrl)
	/* TODO: fill agent, transfer, backend, graph */
	p.Conf.Ctrl.Set(fconfig.APP_CONF_DEFAULT, ConfDefault)
	return p.Conf
}

func (p *Ctrl) String() string {
	return p.Conf.String()
}

// ugly hack
// should called by main package
func (p *Ctrl) Prestart() error {
	glog.V(3).Infof(MODULE_NAME + "Prestart() entering")
	Configure = p.Conf

	EtcdCli = falcon.NewEtcdCli(Configure.Ctrl)

	EtcdCli.Prestart()
	for _, fn := range prestartHooks {
		glog.V(3).Infof(MODULE_NAME+"%s() called\n", runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}
	glog.V(3).Infof(MODULE_NAME + "Prestart() leaving")
	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME + "Start() entering")

	p.status = falcon.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	EtcdCli.Start()
	// p.rpcStart()
	// p.httpStart()
	p.statStart()
	go beego.Run()
	glog.V(3).Infof(MODULE_NAME + "Start() leaving")
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	EtcdCli.Stop()
	// p.httpStop()
	// p.rpcStop()
	return nil
}

// TODO: reload is not yet implemented
func (p *Ctrl) Reload(conf interface{}) error {
	p.Conf = conf.(*config.ConfCtrl)
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	Configure = p.Conf

	EtcdCli.Reload(Configure.Ctrl)
	for _, fn := range prestartHooks {
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}

	return nil
}

func (p *Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return nil
}
