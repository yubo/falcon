/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package config

import (
	"fmt"

	"github.com/yubo/falcon"
)

type Alarm struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Upstream map[int]string
	Configer falcon.Configer
}

func (p Alarm) GetName() string {
	return p.Name
}

func (p Alarm) String() string {
	var s1 string
	for k, v := range p.Upstream {
		s1 += fmt.Sprintf("%d %s\n", k, v)
	}
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s",
		"debug", p.Debug,
		"disabled", p.Disabled,
		"Name", p.Name,
		"Host", p.Host,
		"Upstream", falcon.IndentLines(1, s1),
		p.Configer.String(),
	)
}
