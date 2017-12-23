/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package config

import (
	"fmt"

	"github.com/yubo/falcon/config"
)

type Ctrl struct {
	// only in config.conf
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Metrics  []string
	Ctrl     config.Configer
	Agent    config.Configer
	Transfer config.Configer
	Backend  config.Configer
	// 1: default, 2: db, 3: ConfCtrl.Container
	// height will cover low
}

func (c Ctrl) GetName() string {
	return c.Name
}

func (c Ctrl) String() string {
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
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)\n"+
		"%s (\n%s\n)\n",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"Metrics", config.IndentLines(1, s),
		"ctrl", config.IndentLines(1, c.Ctrl.String()),
		"agent", config.IndentLines(1, c.Agent.String()),
		"transfer", config.IndentLines(1, c.Transfer.String()),
		"backend", config.IndentLines(1, c.Backend.String()),
	)
}
