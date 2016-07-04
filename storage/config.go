/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/hcl"
)

const (
	CACHE_TIME           = 1800 //s
	FIRST_FLUSH_DISK     = 10   //s
	FLUSH_DISK_STEP      = 1    //s
	DEFAULT_HISTORY_SIZE = 3
	CONN_RETRY           = 2 /* history */
)

var (
	defaultOptions = StorageOpts{
		Debug:           false,
		Http:            true,
		HttpAddr:        "0.0.0.0:6071",
		Rpc:             true,
		RpcAddr:         "0.0.0.0:6070",
		RrdStorage:      "/home/work/data/6070",
		Idx:             true,
		IdxInterval:     30,
		IdxFullInterval: 86400,
		Dsn:             "root:@tcp(127.0.0.1:3306)/storage?loc=Local&parseTime=true",
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

/*
func config() *StorageOpts {
	return (*StorageOpts)(atomic.LoadPointer(&configPtr))
}
*/

func applyConfigFile(opts *StorageOpts, filePath string) error {
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

func parse(config *StorageOpts, filename string) {

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
