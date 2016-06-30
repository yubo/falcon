/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"io/ioutil"
	"os"
	"sync/atomic"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon/hcl"
)

const (
	GAUGE           = "GAUGE"
	DERIVE          = "DERIVE"
	COUNTER         = "COUNTER"
	CACHE_TIME      = 1800 //s
	FLUSH_DISK_STEP = 1    //s
	DEFAULT_STEP    = 60   //s
	MIN_STEP        = 30   //s

	/* history */
	defaultHistorySize = 3
)

var (
	defaultOptions = Options{
		Debug:       false,
		Http:        true,
		HttpAddr:    "0.0.0.0:6071",
		Rpc:         true,
		RpcAddr:     "0.0.0.0:6070",
		RrdStorage:  "/home/work/data/6070",
		Dsn:         "root:@tcp(127.0.0.1:3306)/graph?loc=Local&parseTime=true",
		DbMaxIdle:   4,
		CallTimeout: 5000,
		Migrate: MigrateOpt{
			Enable:      false,
			Concurrency: 2,
			Replicas:    500,
			Cluster: map[string]string{
				"graph-00": "127.0.0.1:6070",
			},
		},
	}
)

func config() *Options {
	return (*Options)(atomic.LoadPointer(&configPtr))
}

func applyConfigFile(opts *Options, filePath string) error {
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

	return nil
}

func parse(config string) {

	_, err := os.Stat(config)
	if !os.IsNotExist(err) {
		if err := applyConfigFile(&configOpts, config); err != nil {
			glog.Errorln(err)
			os.Exit(2)
		}
	}

	if configOpts.Migrate.Enable && len(configOpts.Migrate.Cluster) == 0 {
		configOpts.Migrate.Enable = false
	}

	// set config
	atomic.StorePointer(&configPtr, unsafe.Pointer(&configOpts))

	glog.V(3).Infof("ParseConfig ok, file %s", config)
	configFile = config
}
