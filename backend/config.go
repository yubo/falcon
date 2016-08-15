/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

/*
#include "cache.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/yubo/falcon/hcl"
	"github.com/yubo/falcon/specs"
)

const (
	CACHE_TIME              = 1800 //s
	FIRST_FLUSH_DISK        = 1    //s
	FLUSH_DISK_STEP         = 1    //s
	DEFAULT_HISTORY_SIZE    = 3
	CONN_RETRY              = 2
	CACHE_SIZE              = C.CACHE_SIZE     // must pow(2,n)
	CACHE_SIZE_MASK         = C.CACHE_SIZE - 1 //
	DATA_TIMESTAMP_REGULATE = true
	INDEX_QPS               = 100
	INDEX_UPDATE_CYCLE_TIME = 86400
	INDEX_TIMEOUT           = 86400
	INDEX_TRASH_LOOPTIME    = 600
	INDEX_MAX_OPEN_CONNS    = 4
	DEBUG_MULTIPLES         = 20    // demo 时间倍数
	DEBUG_STEP              = 60    //
	DEBUG_SAMPLE_NB         = 18000 //单周期生成样本数量
	DEBUG_STAT_MODULE       = ST_M_CACHE | ST_M_INDEX
	DEBUG_STAT_STEP         = 60
)

var (
	defaultOptions = BackendOpts{
		Debug:           0,
		Http:            true,
		HttpAddr:        "0.0.0.0:7021",
		Rpc:             true,
		RpcAddr:         "0.0.0.0:7020",
		Idx:             true,
		IdxInterval:     30,
		IdxFullInterval: 86400,
		Dsn:             "falcon:1234@tcp(127.0.0.1:3306)/falcon?loc=Local&parseTime=true",
		DbMaxIdle:       4,
		Migrate: MigrateOpts{
			Enable:      false,
			Concurrency: 2,
			Replicas:    500,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Upstream:    map[string]string{},
		},
		Storage: StorageOpts{
			Type: "rrdlite",
		},
	}
)

type MigrateOpts struct {
	Enable      bool              `hcl:"enable"`
	Concurrency int               `hcl:"concurrency"`
	Replicas    int               `hcl:"replicas"`
	CallTimeout int               `hcl:"call_timeout"`
	ConnTimeout int               `hcl:"conn_timeout"`
	Upstream    map[string]string `hcl:"upstream"`
}

func (o MigrateOpts) String() string {
	var upstream string
	indent := strings.Repeat(" ", specs.IndentSize)

	for k, v := range o.Upstream {
		upstream += fmt.Sprintf("%s%-10s = %s\n", indent, k, v)
	}

	return fmt.Sprintf("%-12s %v\n%-12s %d\n%-12s %d\n"+
		"%-12s %d\n%-12s %d\n%s (\n%s\n)",
		"enable", o.Enable, "concurrency", o.Concurrency,
		"replicas", o.Replicas, "callTimeout", o.CallTimeout,
		"conntimeout", o.ConnTimeout, "upstream", strings.TrimRight(upstream, "\n"))
}

type StorageOpts struct {
	Type   string   `hcl:"type"`
	Hdisks []string `hcl:"hdisks"`
}

func (o StorageOpts) String() string {
	var ret string
	indent := strings.Repeat(" ", specs.IndentSize)

	for _, v := range o.Hdisks {
		ret += fmt.Sprintf("%s%s\n", indent, v)
	}

	return fmt.Sprintf("%-12s %s\n%s (\n%s\n)",
		"type", o.Type, "hdisks", strings.TrimRight(ret, "\n"))
}

type ShmOpts struct {
	Magic uint32 `hcl:"magic_code"`
	Key   int    `hcl:"key_start_id"`
	Size  int    `hcl:"segment_size"`
}

func (o ShmOpts) String() string {
	var ret string
	indent := strings.Repeat(" ", specs.IndentSize)

	ret += fmt.Sprintf("%s%-14s 0x%x\n", indent, "magic_code", o.Magic)
	ret += fmt.Sprintf("%s%-14s 0x%x\n", indent, "key_start_id", o.Key)
	ret += fmt.Sprintf("%s%-14s %v\n", indent, "segment_size", o.Size)

	return ret
}

/* config */
type BackendOpts struct {
	Debug           int         `hcl:"debug"`
	PidFile         string      `hcl:"pid_file"`
	Http            bool        `hcl:"http"`
	HttpAddr        string      `hcl:"http_addr"`
	Rpc             bool        `hcl:"rpc"`
	RpcAddr         string      `hcl:"rpc_addr"`
	Idx             bool        `hcl:"idx"`
	IdxInterval     int         `hcl:"idx_interval"`
	IdxFullInterval int         `hcl:"idx_full_interval"`
	Dsn             string      `hcl:"dsn"`
	DbMaxIdle       int         `hcl:"db_max_idle"`
	Migrate         MigrateOpts `hcl:"migrate"`
	Storage         StorageOpts `hcl:"storage"`
	Shm             ShmOpts     `hcl:"shm"`
}

func (o *BackendOpts) String() string {
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %s\n%-17s %v\n%-17s %s\n%-17s %v\n"+
		"%-17s %s\n%-17s %v\n%-17s %d\n%-17s %d\n"+
		"%-17s %s\n%-17s %d\n%s (\n%s\n)\n%s "+
		"(\n%s\n)\n%s (\n%s\n)",
		"debug", o.Debug,
		"pid file", o.PidFile,
		"http", o.Http,
		"httpaddr", o.HttpAddr,
		"rpc", o.Rpc,
		"rpcaddr", o.RpcAddr,
		"idx", o.Idx,
		"idx_interval", o.IdxInterval,
		"idx_full_interval", o.IdxFullInterval,
		"dsn", o.Dsn,
		"dbmaxidle", o.DbMaxIdle,
		"migrate", specs.IndentLines(1, o.Migrate.String()),
		"storage", specs.IndentLines(1, o.Storage.String()),
		"shm", specs.IndentLines(1, o.Shm.String()))
}

func applyConfigFile(opts *BackendOpts, filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	fileString := []byte{}
	glog.V(3).Infof("Loading config file at: %s", filePath)

	fileString, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := hcl.Decode(opts, string(fileString)); err != nil {
		return err
	}

	glog.V(3).Infof("config options:\n%s", opts)
	return nil
}

func parse(config *BackendOpts, filename string) {

	if err := applyConfigFile(config, filename); err != nil {
		glog.Errorln(err)
		os.Exit(2)
	}

	if config.Migrate.Enable && len(config.Migrate.Upstream) == 0 {
		config.Migrate.Enable = false
	}

	// set config
	//atomic.StorePointer(&configPtr, unsafe.Pointer(&configOpts))

	glog.V(3).Infof("ParseConfig ok, file %s", filename)
	glog.V(3).Infof("cache_size %d", CACHE_SIZE)
	appConfigfile = filename
}
