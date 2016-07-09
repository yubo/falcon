/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

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
	CACHE_SIZE              = 1 << 5   // must pow(2,n)
	CACHE_SIZE_MASK         = 1<<5 - 1 //
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
		RrdStorage:      "/home/work/data/7020",
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
type BackendOpts struct {
	Debug           int         `hcl:"debug"`
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

func (o *BackendOpts) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
debug             %d
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
	appConfigfile = filename
}
