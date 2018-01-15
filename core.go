/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

const (
	VERSION      = "0.0.3"
	DEFAULT_STEP = 60 //s
	IndentSize   = 4
	//MIN_STEP     = 30 //s
	REPLICAS = 500
	//GAUGE       = "GAUGE"
	//DERIVE      = "DERIVE"
	//COUNTER     = "COUNTER"
	MODULE_NAME = "\x1B[30m[FALCON]\x1B[0m"

	C_ETCD_ENDPOINTS = "etcdendpoints"
	C_ETCD_USERNAME  = "etcdusername"
	C_ETCD_PASSWORD  = "etcdpassword"
	C_ETCD_CERTFILE  = "certfile"
	C_ETCD_KEYFILE   = "keyfile"
	C_ETCD_CAFILE    = "cafile"
	C_LEASE_KEY      = "leasekey"
	C_LEASE_VALUE    = "leasevalue"
	C_LEASE_TTL      = "leasettl"
)

type CmdOpts struct {
	ConfigFile string
	Module     string
}

var (
	ModuleTpls map[string]Module
	Modules    map[string]Module
)
