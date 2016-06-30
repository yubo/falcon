/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"errors"
	"sync"

	"github.com/yubo/falcon/specs"
)

var (
	ErrUnsupported = errors.New("unsupported")
	ErrExist       = errors.New("entry exists")
	ErrNoent       = errors.New("entry not exists")
	ErrParam       = errors.New("param error")
)

const (
	_ = iota
	IO_TASK_M_READ
	IO_TASK_M_WRITE
	IO_TASK_M_COMMIT
	IO_TASK_M_CHECKOUT
)

const (
	_               = iota
	NET_TASK_M_SEND // no used
	NET_TASK_M_QUERY
	NET_TASK_M_FETCH
)

const (
	GRAPH_F_MISS uint32 = 1 << iota
	GRAPH_F_ERR
	GRAPH_F_SENDING
	GRAPH_F_FETCHING
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
	idx_prev  *cacheEntry
	idx_next  *cacheEntry
	data_prev *cacheEntry
	data_next *cacheEntry
	key       string
	idxTs     int64
	commitTs  int64
	createTs  int64
	putTs     int64
	endpoint  string
	metric    string
	tags      map[string]string
	dsType    string
	step      int
	heartbeat int
	min       string
	max       string
	cache     []*specs.RRDData
	history   []*specs.RRDData
}

/* config */
type Options struct {
	Debug       bool       `hcl:"debug"`
	Http        bool       `hcl:"http"`
	HttpAddr    string     `hcl:"http_addr"`
	Rpc         bool       `hcl:"rpc"`
	RpcAddr     string     `hcl:"rpc_addr"`
	RrdStorage  string     `hcl:"rrd_storage"`
	Dsn         string     `hcl:"dsn"`
	DbMaxIdle   int        `hcl:"db_max_idle"`
	CallTimeout int        `hcl:"call_timeout"`
	Migrate     MigrateOpt `hcl:"migrate"`
}

/* http */
type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

/* rrdtool */
type io_task_t struct {
	method int
	args   interface{}
	done   chan error
}

type rrdCheckout_t struct {
	filename string
	cf       string
	start    int64
	end      int64
	step     int
	data     []*specs.RRDData
}

type flushfile_t struct {
	filename string
	items    []*specs.GraphItem
}

type readfile_t struct {
	filename string
	data     []byte
}

type Net_task_t struct {
	Method int
	e      *cacheEntry
	Done   chan error
	Args   interface{}
	Reply  interface{}
}

type MigrateOpt struct {
	Enable      bool              `hcl:"enable"`
	Concurrency int               `hcl:"concurrency"`
	Replicas    int               `hcl:"replicas"`
	Cluster     map[string]string `hcl:"cluster"`
}
