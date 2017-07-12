/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package config

import (
	"fmt"

	"github.com/yubo/falcon/utils"
)

type ConfCtrl struct {
	// only in falcon.conf
	Debug       int
	Disabled    bool
	Name        string
	Host        string
	Metrics     []string
	Ctrl        utils.Configer
	Agent       utils.Configer
	Loadbalance utils.Configer
	Backend     utils.Configer
	Graph       utils.Configer
	Transfer    utils.Configer
	// 1: default, 2: db, 3: ConfCtrl.Container
	// height will cover low
}

func (c ConfCtrl) String() string {
	var s string
	for k, v := range c.Metrics {
		s += fmt.Sprintf("%s ", v)
		if k%5 == 4 {
			s += "\n"
		}
	}
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"Metrics", utils.IndentLines(1, s),
		c.Ctrl.String(),
	)
}

var (
	ConfDefault = map[string]string{
		//C_RUN_MODE:                "pub",
		utils.C_MASTER_MODE:             "true",
		utils.C_MI_MODE:                 "false",
		utils.C_DEV_MODE:                "false",
		utils.C_HTTP_ADDR:               "8001",
		utils.C_SESSION_GC_MAX_LIFETIME: "86400",
		utils.C_SESSION_COOKIE_LIFETIME: "86400",
		utils.C_AUTH_MODULE:             "ldap",
		utils.C_CACHE_MODULE:            "host,role,system,tag,user",
		utils.C_DB_MAX_CONN:             "30",
		utils.C_DB_MAX_IDLE:             "30",
		utils.C_MI_NORNS_URL:            "http://norns.dev/api/v1/tagstring/cop.xiaomi/hostinfos",
		utils.C_MI_NORNS_INTERVAL:       "5",
	}
)
