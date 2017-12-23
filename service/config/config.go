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

type Migrate struct {
	Disabled bool
	Upstream map[string]string
}

func (p Migrate) String() string {
	var s string

	for k, v := range p.Upstream {
		s += fmt.Sprintf("%-17s %s\n", k, v)
	}
	if s != "" {
		s = fmt.Sprintf("\n%s\n", config.IndentLines(1, s))
	}

	return fmt.Sprintf("%-17s %v\n"+
		"%s (%s)",
		"disable", p.Disabled,
		"cluster", s)
}

type Service struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Migrate  Migrate
	Configer config.Configer
}

func (c Service) GetName() string {
	return c.Name
}

func (c Service) String() string {
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
		"migrate", config.IndentLines(1, c.Migrate.String()),
		c.Configer.String(),
	)
}
