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

type ConfAgent struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Configer utils.Configer
}

func (c ConfAgent) String() string {
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

var (
	ConfDefault = map[string]string{
		utils.C_CONN_TIMEOUT:     "1000",
		utils.C_CALL_TIMEOUT:     "5000",
		utils.C_WORKER_PROCESSES: "2",
		utils.C_HTTP_ENABLE:      "true",
		utils.C_HTTP_ADDR:        "127.0.0.1:1988",
		utils.C_RPC_ENABLE:       "true",
		utils.C_RPC_ADDR:         "127.0.0.1:1989",
		utils.C_GRPC_ENABLE:      "true",
		utils.C_GRPC_ADDR:        "127.0.0.1:1990",
		utils.C_INTERVAL:         "60",
		utils.C_PAYLOADSIZE:      "16",
		utils.C_IFACE_PREFIX:     "eth,em",
	}
)
