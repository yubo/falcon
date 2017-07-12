/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/utils"
)

type Kv struct {
	Key     string
	Section string
	Value   string
}

var (
	// for Description
	ModuleMap = map[string]map[string]string{
		"ctrl": map[string]string{
			//utils.C_RUN_MODE:                "dev/prod",
			//utils.C_ENABLE_DOCS:             "ture/false",
			utils.C_MASTER_MODE:             "bool",
			utils.C_MI_MODE:                 "bool",
			utils.C_DEV_MODE:                "bool",
			utils.C_SESSION_GC_MAX_LIFETIME: "int",
			utils.C_SESSION_COOKIE_LIFETIME: "int",
			utils.C_AUTH_MODULE:             "ldap/misso/github/google",
			utils.C_CACHE_MODULE:            "string",
			utils.C_LDAP_ADDR:               "string",
			utils.C_LDAP_BASE_DN:            "string",
			utils.C_LDAP_BIND_DN:            "string",
			utils.C_LDAP_BIND_PWD:           "string",
			utils.C_LDAP_FILTER:             "string",
			utils.C_MISSO_REDIRECT_URL:      "string",
			utils.C_GITHUB_CLIENT_ID:        "string",
			utils.C_GITHUB_CLIENT_SECRET:    "string",
			utils.C_GITHUB_REDIRECT_URL:     "string",
			utils.C_GOOGLE_CLIENT_ID:        "string",
			utils.C_GOOGLE_CLIENT_SECRET:    "string",
			utils.C_GOOGLE_REDIRECT_URL:     "string",
		},
		"graph": map[string]string{
			utils.C_DEBUG:                "bool",
			utils.C_THRESHOLD:            "int",
			utils.C_HTTP_ENABLE:          "bool",
			utils.C_HTTP_LISTEN:          "string",
			utils.C_RPC_ENABLE:           "bool",
			utils.C_RPC_LISTEN:           "string",
			utils.C_GRPC_ENABLE:          "bool",
			utils.C_GRPC_LISTEN:          "string",
			utils.C_RRD_STORAGE:          "string",
			utils.C_DB_DSN:               "string",
			utils.C_DB_MAXIDLE:           "int",
			utils.C_GRAPH_FLUSHBEFORESCP: "bool",
			utils.C_GRAPH_MIGRATING:      "bool",
			utils.C_GRAPH_CONCURRENCY:    "int",
			utils.C_GRAPH_MINSCPINTERVAL: "int",
			utils.C_GRAPH_REPLICAS:       "int",
			utils.C_GRAPH_CONNTIMEOUT:    "int",
			utils.C_GRAPH_CALLTIMEOUT:    "int",
			utils.C_GRAPH_MAXCONNS:       "int",
			utils.C_GRAPH_MAXIDLECONNS:   "int",
			utils.C_GRAPH_CLUSTER:        "string",
			utils.C_CRON_INTERVAL:        "int",
			utils.C_CRON_PUSHADDR:        "string",
			utils.C_CRON_EXPPERIOD:       "int",
			utils.C_LEASETTL:             "int",
		},
		"transfer": map[string]string{
			utils.C_DEBUG:                  "bool",
			utils.C_SAFEGUARD_MAXMEMMBYTE:  "int",
			utils.C_SAFEGUARD_THRESHOLD:    "int",
			utils.C_SAFEGUARD_INTERVAL:     "int",
			utils.C_HTTP_ENABLE:            "bool",
			utils.C_HTTP_LISTEN:            "string",
			utils.C_RPC_ENABLE:             "bool",
			utils.C_RPC_LISTEN:             "string",
			utils.C_GRPC_ENABLE:            "bool",
			utils.C_GRPC_LISTEN:            "string",
			utils.C_SOCKET_ENABLE:          "bool",
			utils.C_SOCKET_LISTEN:          "string",
			utils.C_SOCKET_TIMEOUT:         "int",
			utils.C_MYSQL_ENABLE:           "bool",
			utils.C_MYSQL_ADDR:             "string",
			utils.C_MYSQL_IDLE:             "int",
			utils.C_MYSQL_MAX:              "int",
			utils.C_JUDGE_ENABLE:           "bool",
			utils.C_JUDGE_BATCH:            "int",
			utils.C_JUDGE_CONNTIMEOUT:      "int",
			utils.C_JUDGE_CALLTIMEOUT:      "int",
			utils.C_JUDGE_MAXCONNS:         "int",
			utils.C_JUDGE_MAXIDLE:          "int",
			utils.C_JUDGE_REPLICAS:         "int",
			utils.C_JUDGE_ISGRPCBACKEND:    "bool",
			utils.C_JUDGE_CLUSTER:          "string",
			utils.C_GRAPH_ENABLE:           "bool",
			utils.C_GRAPH_BATCH:            "int",
			utils.C_GRAPH_CONNTIMEOUT:      "int",
			utils.C_GRAPH_CALLTIMEOUT:      "int",
			utils.C_GRAPH_MAXCONNS:         "int",
			utils.C_GRAPH_MAXIDLE:          "int",
			utils.C_GRAPH_REPLICAS:         "int",
			utils.C_GRAPH_MIGRATING:        "bool",
			utils.C_GRAPH_ISGRPCBACKEND:    "bool",
			utils.C_GRAPH_CLUSTER:          "string",
			utils.C_GRAPH_CLUSTERMIGRATING: "string",
			utils.C_TRANSFER_ENABLE:        "bool",
			utils.C_TRANSFER_BATCH:         "int",
			utils.C_TRANSFER_CONNTIMEOUT:   "int",
			utils.C_TRANSFER_CALLTIMEOUT:   "int",
			utils.C_TRANSFER_MAXCONNS:      "int",
			utils.C_TRANSFER_MAXIDLE:       "int",
			utils.C_TRANSFER_CLUSTER:       "string",
			utils.C_NORNS_HOSTNAME:         "string",
			utils.C_NORNS_PDL:              "string",
			utils.C_NORNS_INTERVAL:         "int",
			utils.C_ALARM_UNRECOVERY:       "string",
			utils.C_ALARM_SOLVE:            "string",
			utils.C_LEASETTL:               "int",
		},
	}
	EtcdMap = map[string]map[string]string{
		"graph": map[string]string{
			utils.C_DEBUG:                "/yubo/falcon/graph/config/debug",
			utils.C_THRESHOLD:            "/yubo/falcon/graph/config/threshold",
			utils.C_HTTP_ENABLE:          "/yubo/falcon/graph/config/http/enable",
			utils.C_HTTP_LISTEN:          "/yubo/falcon/graph/config/http/listen",
			utils.C_RPC_ENABLE:           "/yubo/falcon/graph/config/rpc/enable",
			utils.C_RPC_LISTEN:           "/yubo/falcon/graph/config/rpc/listen",
			utils.C_GRPC_ENABLE:          "/yubo/falcon/graph/config/grpc/enable",
			utils.C_GRPC_LISTEN:          "/yubo/falcon/graph/config/grpc/listen",
			utils.C_RRD_STORAGE:          "/yubo/falcon/graph/config/rrd/storage",
			utils.C_DB_DSN:               "/yubo/falcon/graph/config/db/dsn",
			utils.C_DB_MAXIDLE:           "/yubo/falcon/graph/config/db/maxIdle",
			utils.C_GRAPH_FLUSHBEFORESCP: "/yubo/falcon/graph/config/graph/flushBeforeScp",
			utils.C_GRAPH_MIGRATING:      "/yubo/falcon/graph/config/graph/migrating",
			utils.C_GRAPH_CONCURRENCY:    "/yubo/falcon/graph/config/graph/concurrency",
			utils.C_GRAPH_MINSCPINTERVAL: "/yubo/falcon/graph/config/graph/minScpInterval",
			utils.C_GRAPH_REPLICAS:       "/yubo/falcon/graph/config/graph/replicas",
			utils.C_GRAPH_CONNTIMEOUT:    "/yubo/falcon/graph/config/graph/connTimeout",
			utils.C_GRAPH_CALLTIMEOUT:    "/yubo/falcon/graph/config/graph/callTimeout",
			utils.C_GRAPH_MAXCONNS:       "/yubo/falcon/graph/config/graph/maxConns",
			utils.C_GRAPH_MAXIDLECONNS:   "/yubo/falcon/graph/config/graph/maxIdleConns",
			utils.C_GRAPH_CLUSTER:        "/yubo/falcon/graph/config/graph/cluster",
			utils.C_CRON_INTERVAL:        "/yubo/falcon/graph/config/cron/interval",
			utils.C_CRON_PUSHADDR:        "/yubo/falcon/graph/config/cron/pushAddr",
			utils.C_CRON_EXPPERIOD:       "/yubo/falcon/graph/config/cron/expPeriod",
			utils.C_LEASETTL:             "/yubo/falcon/graph/config/leasettl",
		},
		"transfer": map[string]string{
			utils.C_DEBUG:                  "/yubo/falcon/transfer/config/debug",
			utils.C_SAFEGUARD_MAXMEMMBYTE:  "/yubo/falcon/transfer/config/safeguard/maxMemMByte",
			utils.C_SAFEGUARD_THRESHOLD:    "/yubo/falcon/transfer/config/safeguard/threshold",
			utils.C_SAFEGUARD_INTERVAL:     "/yubo/falcon/transfer/config/safeguard/interval",
			utils.C_HTTP_ENABLE:            "/yubo/falcon/transfer/config/http/enable",
			utils.C_HTTP_LISTEN:            "/yubo/falcon/transfer/config/http/listen",
			utils.C_RPC_ENABLE:             "/yubo/falcon/transfer/config/rpc/enable",
			utils.C_RPC_LISTEN:             "/yubo/falcon/transfer/config/rpc/listen",
			utils.C_GRPC_ENABLE:            "/yubo/falcon/transfer/config/grpc/enable",
			utils.C_GRPC_LISTEN:            "/yubo/falcon/transfer/config/grpc/listen",
			utils.C_SOCKET_ENABLE:          "/yubo/falcon/transfer/config/socket/enable",
			utils.C_SOCKET_LISTEN:          "/yubo/falcon/transfer/config/socket/listen",
			utils.C_SOCKET_TIMEOUT:         "/yubo/falcon/transfer/config/socket/timeout",
			utils.C_MYSQL_ENABLE:           "/yubo/falcon/transfer/config/mysql/enable",
			utils.C_MYSQL_ADDR:             "/yubo/falcon/transfer/config/mysql/addr",
			utils.C_MYSQL_IDLE:             "/yubo/falcon/transfer/config/mysql/idle",
			utils.C_MYSQL_MAX:              "/yubo/falcon/transfer/config/mysql/max",
			utils.C_JUDGE_ENABLE:           "/yubo/falcon/transfer/config/judge/enable",
			utils.C_JUDGE_BATCH:            "/yubo/falcon/transfer/config/judge/batch",
			utils.C_JUDGE_CONNTIMEOUT:      "/yubo/falcon/transfer/config/judge/connTimeout",
			utils.C_JUDGE_CALLTIMEOUT:      "/yubo/falcon/transfer/config/judge/callTimeout",
			utils.C_JUDGE_MAXCONNS:         "/yubo/falcon/transfer/config/judge/maxConns",
			utils.C_JUDGE_MAXIDLE:          "/yubo/falcon/transfer/config/judge/maxIdle",
			utils.C_JUDGE_REPLICAS:         "/yubo/falcon/transfer/config/judge/replicas",
			utils.C_JUDGE_ISGRPCBACKEND:    "/yubo/falcon/transfer/config/judge/isGrpcBackend",
			utils.C_JUDGE_CLUSTER:          "/yubo/falcon/transfer/config/judge/cluster",
			utils.C_GRAPH_ENABLE:           "/yubo/falcon/transfer/config/graph/enable",
			utils.C_GRAPH_BATCH:            "/yubo/falcon/transfer/config/graph/batch",
			utils.C_GRAPH_CONNTIMEOUT:      "/yubo/falcon/transfer/config/graph/connTimeout",
			utils.C_GRAPH_CALLTIMEOUT:      "/yubo/falcon/transfer/config/graph/callTimeout",
			utils.C_GRAPH_MAXCONNS:         "/yubo/falcon/transfer/config/graph/maxConns",
			utils.C_GRAPH_MAXIDLE:          "/yubo/falcon/transfer/config/graph/maxIdle",
			utils.C_GRAPH_REPLICAS:         "/yubo/falcon/transfer/config/graph/replicas",
			utils.C_GRAPH_MIGRATING:        "/yubo/falcon/transfer/config/graph/migrating",
			utils.C_GRAPH_ISGRPCBACKEND:    "/yubo/falcon/transfer/config/graph/isGrpcBackend",
			utils.C_GRAPH_CLUSTER:          "/yubo/falcon/transfer/config/graph/cluster",
			utils.C_GRAPH_CLUSTERMIGRATING: "/yubo/falcon/transfer/config/graph/clusterMigrating",
			utils.C_TRANSFER_ENABLE:        "/yubo/falcon/transfer/config/transfer/enable",
			utils.C_TRANSFER_BATCH:         "/yubo/falcon/transfer/config/transfer/batch",
			utils.C_TRANSFER_CONNTIMEOUT:   "/yubo/falcon/transfer/config/transfer/connTimeout",
			utils.C_TRANSFER_CALLTIMEOUT:   "/yubo/falcon/transfer/config/transfer/callTimeout",
			utils.C_TRANSFER_MAXCONNS:      "/yubo/falcon/transfer/config/transfer/maxConns",
			utils.C_TRANSFER_MAXIDLE:       "/yubo/falcon/transfer/config/transfer/maxIdle",
			utils.C_TRANSFER_CLUSTER:       "/yubo/falcon/transfer/config/transfer/cluster",
			utils.C_NORNS_HOSTNAME:         "/yubo/falcon/transfer/config/norns/hostname",
			utils.C_NORNS_PDL:              "/yubo/falcon/transfer/config/norns/pdl",
			utils.C_NORNS_INTERVAL:         "/yubo/falcon/transfer/config/norns/interval",
			utils.C_ALARM_UNRECOVERY:       "/yubo/falcon/transfer/config/alarm/unrecovery",
			utils.C_ALARM_SOLVE:            "/yubo/falcon/transfer/config/alarm/solve",
			utils.C_LEASETTL:               "/yubo/falcon/transfer/config/leasettl",
		},
	}
)

func prepareEtcdConfig() error {
	put := make(map[string]string)
	for module, kv := range EtcdMap {
		ks := make(map[string]bool)
		for _, k := range kv {
			ks[k] = false
		}

		prefix := fmt.Sprintf("/open-falcon/%s/config/", module)
		resp, err := ctrl.EtcdCli.GetPrefix(prefix)
		if err != nil {
			return err
		}
		for _, v := range resp.Kvs {
			if _, ok := ks[string(v.Key)]; ok {
				ks[string(v.Key)] = true
			}
		}
		for k, exist := range ks {
			if !exist {
				put[k] = ""
			}
		}
	}
	return ctrl.EtcdCli.Puts(put)
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

func (op *Operator) SetEtcdConfig(module string, conf map[string]string) error {
	ks, ok := EtcdMap[module]
	if !ok {
		//skip miss hit
		return nil
	}
	ekv := make(map[string]string)
	for k, ek := range ks {
		glog.V(4).Info(MODULE_NAME, k, "->", ek, "=", conf[k])
		if conf[k] != "" {
			ekv[ek] = conf[k]
		}
	}
	return ctrl.EtcdCli.Puts(ekv)
}

func (op *Operator) SetDbConfig(module string, conf map[string]string) error {
	kv, _ := GetDbConfig(op.O, module)
	for k, v := range conf {
		if v != "" {
			kv[k] = v
		}
	}
	v, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	s := string(v)
	_, err = op.O.Raw("INSERT INTO `kv`(`section`, `key`, `value`)"+
		" VALUES (?,'config',?) ON DUPLICATE KEY UPDATE `value`=?",
		module, s, s).Exec()

	return err
}

func (op *Operator) ConfigGet(module string) ([utils.APP_CONF_SIZE]map[string]string, error) {
	var c *utils.Configer

	switch module {
	case "ctrl":
		c = &ctrl.Configure.Ctrl
	case "agent":
		c = &ctrl.Configure.Agent
	case "loadbalance":
		c = &ctrl.Configure.Loadbalance
	case "backend":
		c = &ctrl.Configure.Backend
	case "graph": // for falcon-plus
		c = &ctrl.Configure.Graph
	case "transfer": // for falcon-plus
		c = &ctrl.Configure.Transfer
	default:
		return [utils.APP_CONF_SIZE]map[string]string{}, utils.ErrNoModule
	}

	conf, err := GetDbConfig(op.O, module)
	if err == nil {
		c.Set(utils.APP_CONF_DB, conf)
	}

	return c.Get(), nil
}

func (op *Operator) ConfigerGet(module string) (*utils.Configer, error) {
	var c *utils.Configer

	switch module {
	case "ctrl":
		c = &ctrl.Configure.Ctrl
	case "agent":
		c = &ctrl.Configure.Agent
	case "loadbalance":
		c = &ctrl.Configure.Loadbalance
	case "backend":
		c = &ctrl.Configure.Backend
	case "graph": // for falcon-plus
		c = &ctrl.Configure.Graph
	case "transfer": // for falcon-plus
		c = &ctrl.Configure.Transfer
	default:
		return nil, utils.ErrNoModule
	}

	conf, err := GetDbConfig(op.O, module)
	if err == nil {
		c.Set(utils.APP_CONF_DB, conf)
	}

	return c, nil
}

func (op *Operator) ConfigSet(module string, conf map[string]string) error {
	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		if module == "graph" {
			// disable set expansion from this interface
			delete(conf, utils.C_MIGRATE_ENABLE)
			delete(conf, utils.C_MIGRATE_NEW_ENDPOINT)
		}
		err := op.SetEtcdConfig(module, conf)
		if err != nil {
			return err
		}
		return op.SetDbConfig(module, conf)
	default:
		return utils.ErrNoModule
	}
}

func (op *Operator) OnlineGet(module string) ([]KeyValue, error) {

	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		prefix := fmt.Sprintf("/open-falcon/%s/online/", module)
		resp, err := ctrl.EtcdCli.GetPrefix(prefix)
		if err != nil {
			return nil, err
		}

		ret := make([]KeyValue, len(resp.Kvs))

		for i, kv := range resp.Kvs {
			ret[i] = KeyValue{
				Key:   string(kv.Key[len(prefix):]),
				Value: string(kv.Value),
			}
		}

		return ret, nil
	default:
		return nil, utils.ErrNoModule
	}

}

type ExpansionStatus struct {
	Migrating    bool   `json:"migrating"`
	GraphCluster string `json:"graph_cluster"`
	NewEndpoint  string `json:"new_endpoint"`
}

// expansion
func (op *Operator) ExpansionGet(module string) (ret *ExpansionStatus, err error) {
	var enable string
	ret = &ExpansionStatus{}

	if module != "graph" {
		return nil, utils.ErrUnsupported
	}

	if enable, err = ctrl.EtcdCli.Get(EtcdMap["graph"][utils.C_MIGRATE_ENABLE]); err == nil && enable == "true" {
		ret.Migrating = true
	}

	ret.GraphCluster, _ = ctrl.EtcdCli.Get(EtcdMap["transfer"][utils.C_GRAPH_CLUSTER])
	ret.NewEndpoint, _ = ctrl.EtcdCli.Get(EtcdMap["graph"][utils.C_MIGRATE_NEW_ENDPOINT])

	return ret, nil
}

func (op *Operator) ExpansionBegin(module string, newEndpoint string) error {
	if module != "graph" {
		return utils.ErrUnsupported
	}

	// get
	transfer, err := GetDbConfig(op.O, "transfer")
	if err != nil {
		return err
	}

	replicas, ok := transfer[utils.C_GRAPH_REPLICAS]
	if !ok {
		return errors.New("can not get transfer->graph->replicas config")
	}

	online := make(map[string]string)
	_online, _ := op.OnlineGet("graph")
	for _, kv := range _online {
		online[kv.Key] = kv.Value
	}

	old_cluster := make(map[string]string)
	new_cluster := make(map[string]string)
	_old_cluster, ok := transfer[utils.C_GRAPH_CLUSTER]
	if !ok {
		return errors.New("can not get transfer->graph->cluster config")
	}
	for _, v := range strings.Split(_old_cluster, ";") {
		// alias=hostname
		// s[0]: alias
		// s[1]: hostname
		s := strings.Split(v, "=")
		old_cluster[s[1]] = s[0]
		new_cluster[s[0]] = s[1]
	}

	for _, v := range strings.Split(newEndpoint, ";") {
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return errors.New("endpoint format error " + v)
		}
		h := strings.Split(s[1], ":")
		if _, ok := online[h[0]]; !ok {
			return errors.New(fmt.Sprintf("endpoint(%s) is not online", h[0]))
		}
		new_cluster[s[0]] = s[1]
	}

	_new_cluster := ""
	for k, v := range new_cluster {
		_new_cluster += fmt.Sprintf("%s=%s;", k, v)
	}
	if len(_new_cluster) > 0 {
		_new_cluster = _new_cluster[0 : len(_new_cluster)-1]
	}

	// set
	op.SetDbConfig("graph", map[string]string{
		utils.C_MIGRATE_NEW_ENDPOINT: newEndpoint,
		utils.C_MIGRATE_ENABLE:       "true",
		utils.C_MIGRATE_REPLICAS:     replicas,
		utils.C_MIGRATE_CLUSTER:      _old_cluster,
	})
	op.SetDbConfig("transfer", map[string]string{
		utils.C_GRAPH_CLUSTER: _new_cluster,
	})

	ekv := make(map[string]string)
	ekv[EtcdMap["graph"][utils.C_MIGRATE_NEW_ENDPOINT]] = newEndpoint
	ekv[EtcdMap["graph"][utils.C_MIGRATE_ENABLE]] = "true"
	ekv[EtcdMap["graph"][utils.C_MIGRATE_REPLICAS]] = replicas
	ekv[EtcdMap["graph"][utils.C_MIGRATE_CLUSTER]] = _old_cluster
	ekv[EtcdMap["transfer"][utils.C_GRAPH_CLUSTER]] = _new_cluster

	return ctrl.EtcdCli.Puts(ekv)

}

func (op *Operator) ExpansionFinish(module string) error {
	if module != "graph" {
		return utils.ErrUnsupported
	}

	op.SetDbConfig("graph", map[string]string{
		utils.C_MIGRATE_ENABLE:       "false",
		utils.C_MIGRATE_NEW_ENDPOINT: " ",
	})

	ekv := make(map[string]string)
	ekv[EtcdMap["graph"][utils.C_MIGRATE_ENABLE]] = "false"
	ekv[EtcdMap["graph"][utils.C_MIGRATE_NEW_ENDPOINT]] = " "
	return ctrl.EtcdCli.Puts(ekv)
}
