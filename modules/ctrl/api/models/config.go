/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"github.com/yubo/falcon/lib/core"
	//"github.com/yubo/falcon/graph"
	//"github.com/yubo/falcon/transfer"
)

var (
	// for Description
	ModuleMap = map[string]map[string]string{
	/*
		"ctrl": ctrl.ConfDesc,
			"graph": map[string]string{
				graph.C_DEBUG:                "bool",
				graph.C_THRESHOLD:            "int",
				graph.C_HTTP_ENABLE:          "bool",
				graph.C_HTTP_LISTEN:          "string",
				graph.C_RPC_ENABLE:           "bool",
				graph.C_RPC_LISTEN:           "string",
				graph.C_GRPC_ENABLE:          "bool",
				graph.C_GRPC_LISTEN:          "string",
				graph.C_RRD_STORAGE:          "string",
				graph.C_DB_DSN:               "string",
				graph.C_DB_MAXIDLE:           "int",
				graph.C_GRAPH_FLUSHBEFORESCP: "bool",
				graph.C_GRAPH_MIGRATING:      "bool",
				graph.C_GRAPH_CONCURRENCY:    "int",
				graph.C_GRAPH_MINSCPINTERVAL: "int",
				graph.C_GRAPH_REPLICAS:       "int",
				graph.C_GRAPH_CALLTIMEOUT:    "int",
				graph.C_GRAPH_MAXCONNS:       "int",
				graph.C_GRAPH_MAXIDLECONNS:   "int",
				graph.C_GRAPH_CLUSTER:        "string",
				graph.C_CRON_INTERVAL:        "int",
				graph.C_CRON_PUSHADDR:        "string",
				graph.C_CRON_EXPPERIOD:       "int",
				graph.C_LEASETTL:             "int",
			},
			"transfer": map[string]string{
				transfer.C_DEBUG:                  "bool",
				transfer.C_SAFEGUARD_MAXMEMMBYTE:  "int",
				transfer.C_SAFEGUARD_THRESHOLD:    "int",
				transfer.C_SAFEGUARD_INTERVAL:     "int",
				transfer.C_HTTP_ENABLE:            "bool",
				transfer.C_HTTP_LISTEN:            "string",
				transfer.C_RPC_ENABLE:             "bool",
				transfer.C_RPC_LISTEN:             "string",
				transfer.C_GRPC_ENABLE:            "bool",
				transfer.C_GRPC_LISTEN:            "string",
				transfer.C_SOCKET_ENABLE:          "bool",
				transfer.C_SOCKET_LISTEN:          "string",
				transfer.C_SOCKET_TIMEOUT:         "int",
				transfer.C_MYSQL_ENABLE:           "bool",
				transfer.C_MYSQL_ADDR:             "string",
				transfer.C_MYSQL_IDLE:             "int",
				transfer.C_MYSQL_MAX:              "int",
				transfer.C_JUDGE_ENABLE:           "bool",
				transfer.C_JUDGE_BATCH:            "int",
				transfer.C_JUDGE_CALLTIMEOUT:      "int",
				transfer.C_JUDGE_MAXCONNS:         "int",
				transfer.C_JUDGE_MAXIDLE:          "int",
				transfer.C_JUDGE_REPLICAS:         "int",
				transfer.C_JUDGE_ISGRPCBACKEND:    "bool",
				transfer.C_JUDGE_CLUSTER:          "string",
				transfer.C_GRAPH_ENABLE:           "bool",
				transfer.C_GRAPH_BATCH:            "int",
				transfer.C_GRAPH_CALLTIMEOUT:      "int",
				transfer.C_GRAPH_MAXCONNS:         "int",
				transfer.C_GRAPH_MAXIDLE:          "int",
				transfer.C_GRAPH_REPLICAS:         "int",
				transfer.C_GRAPH_MIGRATING:        "bool",
				transfer.C_GRAPH_ISGRPCBACKEND:    "bool",
				transfer.C_GRAPH_CLUSTER:          "string",
				transfer.C_GRAPH_CLUSTERMIGRATING: "string",
				transfer.C_TRANSFER_ENABLE:        "bool",
				transfer.C_TRANSFER_BATCH:         "int",
				transfer.C_TRANSFER_CALLTIMEOUT:   "int",
				transfer.C_TRANSFER_MAXCONNS:      "int",
				transfer.C_TRANSFER_MAXIDLE:       "int",
				transfer.C_TRANSFER_CLUSTER:       "string",
				transfer.C_NORNS_HOSTNAME:         "string",
				transfer.C_NORNS_PDL:              "string",
				transfer.C_NORNS_INTERVAL:         "int",
				transfer.C_ALARM_UNRECOVERY:       "string",
				transfer.C_ALARM_SOLVE:            "string",
				transfer.C_LEASETTL:               "int",
			},
	*/
	}
	EtcdMap = map[string]map[string]string{
	/*
		"graph": map[string]string{
			graph.C_DEBUG:                "/yubo/falcon/graph/config/debug",
			graph.C_THRESHOLD:            "/yubo/falcon/graph/config/threshold",
			graph.C_HTTP_ENABLE:          "/yubo/falcon/graph/config/http/enable",
			graph.C_HTTP_LISTEN:          "/yubo/falcon/graph/config/http/listen",
			graph.C_RPC_ENABLE:           "/yubo/falcon/graph/config/rpc/enable",
			graph.C_RPC_LISTEN:           "/yubo/falcon/graph/config/rpc/listen",
			graph.C_GRPC_ENABLE:          "/yubo/falcon/graph/config/grpc/enable",
			graph.C_GRPC_LISTEN:          "/yubo/falcon/graph/config/grpc/listen",
			graph.C_RRD_STORAGE:          "/yubo/falcon/graph/config/rrd/storage",
			graph.C_DB_DSN:               "/yubo/falcon/graph/config/db/dsn",
			graph.C_DB_MAXIDLE:           "/yubo/falcon/graph/config/db/maxIdle",
			graph.C_GRAPH_FLUSHBEFORESCP: "/yubo/falcon/graph/config/graph/flushBeforeScp",
			graph.C_GRAPH_MIGRATING:      "/yubo/falcon/graph/config/graph/migrating",
			graph.C_GRAPH_CONCURRENCY:    "/yubo/falcon/graph/config/graph/concurrency",
			graph.C_GRAPH_MINSCPINTERVAL: "/yubo/falcon/graph/config/graph/minScpInterval",
			graph.C_GRAPH_REPLICAS:       "/yubo/falcon/graph/config/graph/replicas",
			graph.C_GRAPH_CALLTIMEOUT:    "/yubo/falcon/graph/config/graph/callTimeout",
			graph.C_GRAPH_MAXCONNS:       "/yubo/falcon/graph/config/graph/maxConns",
			graph.C_GRAPH_MAXIDLECONNS:   "/yubo/falcon/graph/config/graph/maxIdleConns",
			graph.C_GRAPH_CLUSTER:        "/yubo/falcon/graph/config/graph/cluster",
			graph.C_CRON_INTERVAL:        "/yubo/falcon/graph/config/cron/interval",
			graph.C_CRON_PUSHADDR:        "/yubo/falcon/graph/config/cron/pushAddr",
			graph.C_CRON_EXPPERIOD:       "/yubo/falcon/graph/config/cron/expPeriod",
			graph.C_LEASETTL:             "/yubo/falcon/graph/config/leasettl",
		},
		"transfer": map[string]string{
			transfer.C_DEBUG:                  "/yubo/falcon/transfer/config/debug",
			transfer.C_SAFEGUARD_MAXMEMMBYTE:  "/yubo/falcon/transfer/config/safeguard/maxMemMByte",
			transfer.C_SAFEGUARD_THRESHOLD:    "/yubo/falcon/transfer/config/safeguard/threshold",
			transfer.C_SAFEGUARD_INTERVAL:     "/yubo/falcon/transfer/config/safeguard/interval",
			transfer.C_HTTP_ENABLE:            "/yubo/falcon/transfer/config/http/enable",
			transfer.C_HTTP_LISTEN:            "/yubo/falcon/transfer/config/http/listen",
			transfer.C_RPC_ENABLE:             "/yubo/falcon/transfer/config/rpc/enable",
			transfer.C_RPC_LISTEN:             "/yubo/falcon/transfer/config/rpc/listen",
			transfer.C_GRPC_ENABLE:            "/yubo/falcon/transfer/config/grpc/enable",
			transfer.C_GRPC_LISTEN:            "/yubo/falcon/transfer/config/grpc/listen",
			transfer.C_SOCKET_ENABLE:          "/yubo/falcon/transfer/config/socket/enable",
			transfer.C_SOCKET_LISTEN:          "/yubo/falcon/transfer/config/socket/listen",
			transfer.C_SOCKET_TIMEOUT:         "/yubo/falcon/transfer/config/socket/timeout",
			transfer.C_MYSQL_ENABLE:           "/yubo/falcon/transfer/config/mysql/enable",
			transfer.C_MYSQL_ADDR:             "/yubo/falcon/transfer/config/mysql/addr",
			transfer.C_MYSQL_IDLE:             "/yubo/falcon/transfer/config/mysql/idle",
			transfer.C_MYSQL_MAX:              "/yubo/falcon/transfer/config/mysql/max",
			transfer.C_JUDGE_ENABLE:           "/yubo/falcon/transfer/config/judge/enable",
			transfer.C_JUDGE_BATCH:            "/yubo/falcon/transfer/config/judge/batch",
			transfer.C_JUDGE_CALLTIMEOUT:      "/yubo/falcon/transfer/config/judge/callTimeout",
			transfer.C_JUDGE_MAXCONNS:         "/yubo/falcon/transfer/config/judge/maxConns",
			transfer.C_JUDGE_MAXIDLE:          "/yubo/falcon/transfer/config/judge/maxIdle",
			transfer.C_JUDGE_REPLICAS:         "/yubo/falcon/transfer/config/judge/replicas",
			transfer.C_JUDGE_ISGRPCBACKEND:    "/yubo/falcon/transfer/config/judge/isGrpcBackend",
			transfer.C_JUDGE_CLUSTER:          "/yubo/falcon/transfer/config/judge/cluster",
			transfer.C_GRAPH_ENABLE:           "/yubo/falcon/transfer/config/graph/enable",
			transfer.C_GRAPH_BATCH:            "/yubo/falcon/transfer/config/graph/batch",
			transfer.C_GRAPH_CALLTIMEOUT:      "/yubo/falcon/transfer/config/graph/callTimeout",
			transfer.C_GRAPH_MAXCONNS:         "/yubo/falcon/transfer/config/graph/maxConns",
			transfer.C_GRAPH_MAXIDLE:          "/yubo/falcon/transfer/config/graph/maxIdle",
			transfer.C_GRAPH_REPLICAS:         "/yubo/falcon/transfer/config/graph/replicas",
			transfer.C_GRAPH_MIGRATING:        "/yubo/falcon/transfer/config/graph/migrating",
			transfer.C_GRAPH_ISGRPCBACKEND:    "/yubo/falcon/transfer/config/graph/isGrpcBackend",
			transfer.C_GRAPH_CLUSTER:          "/yubo/falcon/transfer/config/graph/cluster",
			transfer.C_GRAPH_CLUSTERMIGRATING: "/yubo/falcon/transfer/config/graph/clusterMigrating",
			transfer.C_TRANSFER_ENABLE:        "/yubo/falcon/transfer/config/transfer/enable",
			transfer.C_TRANSFER_BATCH:         "/yubo/falcon/transfer/config/transfer/batch",
			transfer.C_TRANSFER_CALLTIMEOUT:   "/yubo/falcon/transfer/config/transfer/callTimeout",
			transfer.C_TRANSFER_MAXCONNS:      "/yubo/falcon/transfer/config/transfer/maxConns",
			transfer.C_TRANSFER_MAXIDLE:       "/yubo/falcon/transfer/config/transfer/maxIdle",
			transfer.C_TRANSFER_CLUSTER:       "/yubo/falcon/transfer/config/transfer/cluster",
			transfer.C_NORNS_HOSTNAME:         "/yubo/falcon/transfer/config/norns/hostname",
			transfer.C_NORNS_PDL:              "/yubo/falcon/transfer/config/norns/pdl",
			transfer.C_NORNS_INTERVAL:         "/yubo/falcon/transfer/config/norns/interval",
			transfer.C_ALARM_UNRECOVERY:       "/yubo/falcon/transfer/config/alarm/unrecovery",
			transfer.C_ALARM_SOLVE:            "/yubo/falcon/transfer/config/alarm/solve",
			transfer.C_LEASETTL:               "/yubo/falcon/transfer/config/leasettl",
		},
	*/
	}
)

func (op *Operator) SetEtcdConfig(module string, conf map[string]string) error {
	/* TODO
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
	return ctrl.EtcdPuts(ekv)
	*/
	return nil
}

func (op *Operator) SetDbConfig(module string, conf map[string]string) error {
	/* TODO
	kv, _ := ctrl.GetDbConfig(op.O, module)
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
	*/

	return nil
}

func (op *Operator) ConfigGet(module string) (map[string]interface{}, error) {
	/* TODO
	var c *core.Configer

	switch module {
	case "ctrl":
		c = &ctrl.Configure.Ctrl
	case "agent":
		c = &ctrl.Configure.Agent
	//case "loadbalance":
	//	c = &ctrl.Configure.Loadbalance
	//case "backend":
	//	c = &ctrl.Configure.Backend
	//case "graph": // for falcon-plus
	//	c = &ctrl.Configure.Graph
	//case "transfer": // for falcon-plus
	//	c = &ctrl.Configure.Transfer
	default:
		return [falcon.APP_CONF_SIZE]map[string]string{}, falcon.ErrNoModule
	}

	conf, err := ctrl.GetDbConfig(op.O, module)
	if err == nil {
		c.Set(falcon.APP_CONF_DB, conf)
	}

	return c.Get(), nil
	*/
	return nil, nil
}

func (op *Operator) ConfigerGet(module string) (*core.Configer, error) {
	/* TODO
	var c *falcon.Configer

	switch module {
	case "ctrl":
		c = &ctrl.Configure.Ctrl
	case "agent":
		c = &ctrl.Configure.Agent
	//case "loadbalance":
	//	c = &ctrl.Configure.Loadbalance
	//case "backend":
	//	c = &ctrl.Configure.Backend
	//case "graph": // for falcon-plus
	//	c = &ctrl.Configure.Graph
	//case "transfer": // for falcon-plus
	//	c = &ctrl.Configure.Transfer
	default:
		return nil, falcon.ErrNoModule
	}

	conf, err := ctrl.GetDbConfig(op.O, module)
	if err == nil {
		c.Set(falcon.APP_CONF_DB, conf)
	}

	return c, nil
	*/
	return nil, nil
}

func (op *Operator) ConfigSet(module string, conf map[string]string) error {
	/* TODO
	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		//if module == "graph" {
		//	// disable set expansion from this interface
		//	delete(conf, graph.C_MIGRATE_ENABLE)
		//	delete(conf, graph.C_MIGRATE_NEW_ENDPOINT)
		//}
		err := op.SetEtcdConfig(module, conf)
		if err != nil {
			return err
		}
		return op.SetDbConfig(module, conf)
	default:
		return falcon.ErrNoModule
	}
	*/
	return nil
}

func (op *Operator) OnlineGet(module string) ([]KeyValue, error) {

	/* TODO
	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		prefix := fmt.Sprintf("/open-falcon/%s/online/", module)
		resp, err := ctrl.EtcdGetPrefix(prefix)
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
		return nil, falcon.ErrNoModule
	}
	*/
	return nil, nil

}

type ExpansionStatus struct {
	Migrating    bool   `json:"migrating"`
	GraphCluster string `json:"graph_cluster"`
	NewEndpoint  string `json:"new_endpoint"`
}

// expansion
func (op *Operator) ExpansionGet(module string) (ret *ExpansionStatus, err error) {
	//var enable string
	ret = &ExpansionStatus{}

	if module != "graph" {
		return nil, core.ErrUnsupported
	}

	//if enable, err = ctrl.EtcdGet(EtcdMap["graph"][ctrl.C_MIGRATE_ENABLE]); err == nil && enable == "true" {
	//	ret.Migrating = true
	//}

	//ret.GraphCluster, _ = ctrl.EtcdGet(EtcdMap["transfer"][transfer.C_GRAPH_CLUSTER])
	//ret.NewEndpoint, _ = ctrl.EtcdGet(EtcdMap["graph"][graph.C_MIGRATE_NEW_ENDPOINT])

	return ret, nil
}

func (op *Operator) ExpansionBegin(module string, newEndpoint string) error {
	/*
		if module != "graph" {
			return falcon.ErrUnsupported
		}

		// get
		transfer, err := ctrl.GetDbConfig(op.O, "transfer")
		if err != nil {
			return err
		}

		replicas, ok := transfer[transfer.C_GRAPH_REPLICAS]
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
		_old_cluster, ok := transfer[transfer.C_GRAPH_CLUSTER]
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
			graph.C_MIGRATE_NEW_ENDPOINT: newEndpoint,
			graph.C_MIGRATE_ENABLE:       "true",
			graph.C_MIGRATE_REPLICAS:     replicas,
			graph.C_MIGRATE_CLUSTER:      _old_cluster,
		})
		op.SetDbConfig("transfer", map[string]string{
			transfer.C_GRAPH_CLUSTER: _new_cluster,
		})

		ekv := make(map[string]string)
		ekv[EtcdMap["graph"][graph.C_MIGRATE_NEW_ENDPOINT]] = newEndpoint
		ekv[EtcdMap["graph"][graph.C_MIGRATE_ENABLE]] = "true"
		ekv[EtcdMap["graph"][graph.C_MIGRATE_REPLICAS]] = replicas
		ekv[EtcdMap["graph"][graph.C_MIGRATE_CLUSTER]] = _old_cluster
		ekv[EtcdMap["transfer"][transfer.C_GRAPH_CLUSTER]] = _new_cluster

		return ctrl.EtcdPuts(ekv)
	*/
	return nil

}

func (op *Operator) ExpansionFinish(module string) error {
	/*
		if module != "graph" {
			return falcon.ErrUnsupported
		}

		op.SetDbConfig("graph", map[string]string{
			graph.C_MIGRATE_ENABLE:       "false",
			graph.C_MIGRATE_NEW_ENDPOINT: " ",
		})

		ekv := make(map[string]string)
		ekv[EtcdMap["graph"][graph.C_MIGRATE_ENABLE]] = "false"
		ekv[EtcdMap["graph"][graph.C_MIGRATE_NEW_ENDPOINT]] = " "
		return ctrl.EtcdPuts(ekv)
	*/
	return nil
}
