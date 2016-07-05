/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/yudai/hcl"
)

const (
	CONN_RETRY = 2
)

var (
	DefaultOptions = AgentOpts{
		Debug:    false,
		PidFile:  "/var/run/agent.pid",
		Host:     "",
		Http:     true,
		HttpAddr: "0.0.0.0:1988",
		IfPre:    []string{"eth", "em"},
		Interval: 60,
		Handoff: HandoffOpts{
			Enable:      true,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Upstreams:   []string{},
		},
	}
)

type HandoffOpts struct {
	Enable      bool     `hcl:"enable"`
	ConnTimeout int      `hcl:"conn_timeout"`
	CallTimeout int      `hcl:"call_timeout"`
	Upstreams   []string `hcl:"upstreams"`
}

type AgentOpts struct {
	Debug    bool        `hcl:"debug"`
	PidFile  string      `hcl:"pid_file"`
	Host     string      `hcl:"host"`
	Http     bool        `hcl:"http"`
	HttpAddr string      `hcl:"http_addr"`
	IfPre    []string    `hcl:"iface_prefix"`
	Interval int         `hcl:"interval"`
	Handoff  HandoffOpts `hcl:"handoff"`
}

func applyConfigFile(opts *AgentOpts, filePath string) error {
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

func parse(config *AgentOpts, filename string) {

	if err := applyConfigFile(config, filename); err != nil {
		glog.Errorln(err)
		os.Exit(2)
	}
	if config.Host == "" {
		config.Host, _ = os.Hostname()
	}

	glog.V(3).Infof("ParseConfig ok, file %s", filename)
	appConfigFile = filename
}
