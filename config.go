/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"fmt"
)

type ModuleConf interface {
	GetName() string
	String() string
}

type FalconConfig struct {
	ConfigFile string
	PidFile    string
	Log        string
	Logv       int
	Conf       []ModuleConf
}

func (p FalconConfig) String() string {
	ret := fmt.Sprintf("%-17s %s"+
		"\n%-17s %s"+
		"\n%-17s %d",
		"pidfile", p.PidFile,
		"log", p.Log,
		"logv", p.Logv,
	)
	for _, c := range p.Conf {
		ret += fmt.Sprintf("\n%s (\n%s\n)",
			c.GetName(),
			IndentLines(1, c.String()))
	}
	return ret
}
