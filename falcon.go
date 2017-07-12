/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"fmt"

	"github.com/yubo/falcon/utils"

	agent "github.com/yubo/falcon/agent/config"
	ctrl "github.com/yubo/falcon/ctrl/config"
)

const (
	VERSION      = "0.0.3"
	IndentSize   = 4
	DEFAULT_STEP = 60 //s
	//MIN_STEP     = 30 //s
	REPLICAS = 500
	//GAUGE       = "GAUGE"
	//DERIVE      = "DERIVE"
	//COUNTER     = "COUNTER"
	MODULE_NAME = "\x1B[32m[FALCON]\x1B[0m "
)

type CmdOpts struct {
	ConfigFile string
	Args       []string
}

var (
	ModuleTpls map[string]Module
)

type FalconConfig struct {
	ConfigFile string
	PidFile    string
	Log        string
	Logv       int
	Conf       []interface{}
}

func (p FalconConfig) String() string {
	ret := fmt.Sprintf("%-17s %s"+
		"\n%-17s %s"+
		"\n%-17s %d",
		"pidfile", p.PidFile,
		"log", p.Log,
		"logv", p.Logv,
	)
	for _, v := range p.Conf {
		switch utils.GetType(v) {
		case "ConfAgent":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*agent.ConfAgent).Name,
				utils.IndentLines(1, v.(*agent.ConfAgent).String()))
		case "ConfCtrl":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*ctrl.ConfCtrl).Name,
				utils.IndentLines(1, v.(*ctrl.ConfCtrl).String()))
			/*
				case "ConfLoadbalance":
					ret += fmt.Sprintf("\n%s (\n%s\n)",
						v.(*ConfLoadbalance).Name,
						utils.IndentLines(1, v.(*ConfLoadbalance).String()))
				case "ConfBackend":
					ret += fmt.Sprintf("\n%s (\n%s\n)",
						v.(*ConfBackend).Name,
						utils.IndentLines(1, v.(*ConfBackend).String()))
			*/
		}
	}
	return ret
}
