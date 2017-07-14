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

type TransferBackend struct {
	Disabled bool
	Name     string
	Type     string
	Upstream map[string]string
}

func (p TransferBackend) String() string {
	var s1, s2 string

	s1 = fmt.Sprintf("%s %s", p.Type, p.Name)
	if p.Disabled {
		s1 += "(Disable)"
	}

	for k, v := range p.Upstream {
		s2 += fmt.Sprintf("%-17s %s\n", k, v)
	}
	return fmt.Sprintf("%s cluster (\n%s\n)", s1, config.IndentLines(1, s2))
}

type ConfTransfer struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Backend  []TransferBackend
	Configer config.Configer
}

func (p ConfTransfer) GetName() string {
	return p.Name
}

func (p ConfTransfer) String() string {
	var s1 string
	for _, v := range p.Backend {
		s1 += fmt.Sprintf("%s\n", v.String())
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
		"backend", config.IndentLines(1, s1),
		p.Configer.String(),
	)
}
