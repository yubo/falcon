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

type TransferShareMap struct {
	ShareCount int
	ShareMap   map[int]string
}

type ConfTransfer struct {
	Debug      int
	Disabled   bool
	Name       string
	Host       string
	ShareMap   map[int]string
	ShareCount int
	Configer   config.Configer
}

func (p ConfTransfer) GetName() string {
	return p.Name
}

func (p ConfTransfer) String() string {
	var s1 string
	for k, v := range p.ShareMap {
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
		"shareMap", config.IndentLines(1, s1),
		p.Configer.String(),
	)
}
