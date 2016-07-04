/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package handoff

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/yubo/falcon/hcl"
)

const (
	CONN_RETRY = 2
)

var (
	defaultOptions = HandoffOpts{
		Debug:       false,
		Http:        true,
		HttpAddr:    "0.0.0.0:6060",
		Rpc:         true,
		RpcAddr:     "0.0.0.0:8433",
		Replicas:    500,
		Concurrency: 2,
		Backends:    make(map[string]BackendOpt),
	}
)

type BackendOpt struct {
	Enable      bool              `hcl:"enable"`
	Type        string            `hcl:"type"`
	Sched       string            `hcl:"sched"`
	Batch       int               `hcl:"batch"`
	ConnTimeout int               `hcl:"conn_timeout"`
	CallTimeout int               `hcl:"call_timeout"`
	Upstream    map[string]string `hcl:"upstream"`
}

type HandoffOpts struct {
	Debug       bool                  `hcl:"debug"`
	Http        bool                  `hcl:"http"`
	HttpAddr    string                `hcl:"http_addr"`
	Rpc         bool                  `hcl:"rpc"`
	RpcAddr     string                `hcl:"rpc_addr"`
	Replicas    int                   `hcl:"replicas"`
	Concurrency int                   `hcl:"concurrency"`
	Backends    map[string]BackendOpt `hcl:"backends"`
}

func applyConfigFile(opts *HandoffOpts, filePath string) error {
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

func parse(config *HandoffOpts, filename string) {

	if err := applyConfigFile(config, filename); err != nil {
		glog.Errorln(err)
		os.Exit(2)
	}

	glog.V(3).Infof("ParseConfig ok, file %s", filename)
	appConfigFile = filename
}
