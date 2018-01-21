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

type Agent struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Configer falcon.Configer
}

func (c Agent) GetName() string {
	return c.Name
}

func (c Agent) String() string {
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		c.Configer.String(),
	)
}
