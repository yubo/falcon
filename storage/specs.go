/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yubo/falcon/specs"
)

const (
	_ = iota
	IO_TASK_M_FILE_READ
	IO_TASK_M_FILE_WRITE
	IO_TASK_M_RRD_UPDATE
	IO_TASK_M_RRD_FETCH
)

const (
	_               = iota
	NET_TASK_M_SEND // no used
	NET_TASK_M_QUERY
	NET_TASK_M_FETCH
)

const (
	RRD_F_MISS uint32 = 1 << iota
	RRD_F_ERR
	RRd_F_SENDING
	RRD_F_FETCHING
)

/* cache queue */
type cacheq struct {
	sync.RWMutex
	first *cacheEntry
	last  *cacheEntry
}

/* cache */
type cacheEntry struct {
	sync.RWMutex
	flag      uint32
	prev      *cacheEntry
	next      *cacheEntry
	idxPrev   *cacheEntry
	idxNext   *cacheEntry
	key       string
	idxTs     int64
	commitTs  int64
	createTs  int64
	ts        int64
	host      string
	k         string
	tags      string
	typ       string
	step      int
	heartbeat int
	min       string
	max       string
	cache     []*specs.RRDData
	history   []*specs.RRDData
}

/* rrdtool */
type rrdCheckout struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*specs.RRDData
}

type ioTask struct {
	method int
	args   interface{}
	done   chan error
}

type netTask struct {
	Method int
	e      *cacheEntry
	Done   chan error
	Args   interface{}
	Reply  interface{}
}

type MigrateOpts struct {
	Enable      bool              `hcl:"enable"`
	Concurrency int               `hcl:"concurrency"`
	Replicas    int               `hcl:"replicas"`
	CallTimeout int               `hcl:"call_timeout"`
	ConnTimeout int               `hcl:"conn_timeout"`
	Upstream    map[string]string `hcl:"upstream"`
}

func (o MigrateOpts) String() string {
	indent := strings.Repeat(" ", specs.IndentSize)

	ret := fmt.Sprintf(`enable      %v
concurrency %d
replicas    %d
calltimeout %d
conntimeout %d
upstream (
`, o.Enable, o.Concurrency, o.Replicas,
		o.CallTimeout, o.ConnTimeout)
	for k, v := range o.Upstream {
		ret += fmt.Sprintf("%s%10s = %s\n", indent, k, v)
	}
	ret += ")"
	return ret
}

/* config */
type StorageOpts struct {
	Debug           bool        `hcl:"debug"`
	PidFile         string      `hcl:"pid_file"`
	Http            bool        `hcl:"http"`
	HttpAddr        string      `hcl:"http_addr"`
	Rpc             bool        `hcl:"rpc"`
	RpcAddr         string      `hcl:"rpc_addr"`
	RrdStorage      string      `hcl:"rrd_storage"`
	Idx             bool        `hcl:"idx"`
	IdxInterval     int         `hcl:"idx_interval"`
	IdxFullInterval int         `hcl:"idx_full_interval"`
	Dsn             string      `hcl:"dsn"`
	DbMaxIdle       int         `hcl:"db_max_idle"`
	Migrate         MigrateOpts `hcl:"migrate"`
}

func (o *StorageOpts) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
debug             %v
http              %v
httpaddr          %s
rpc               %v
rpcaddr           %s
rrdstorage        %s
idx               %v
idx_interval      %d
idx_full_interval %d
dsn               %s
dbmaxidle         %d
migrage (
%s
)`,
		o.Debug, o.Http, o.HttpAddr, o.Rpc,
		o.RpcAddr, o.RrdStorage, o.Idx,
		o.IdxInterval, o.IdxFullInterval, o.Dsn,
		o.DbMaxIdle, specs.IndentLines(1, o.Migrate.String())))
}
