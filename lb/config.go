/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

import "github.com/yubo/falcon/specs"

const (
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
	CTRL_STEP       = 360
)

var (
	DefaultLb = Lb{
		Params: specs.ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Disabled:    false,
			Http:        true,
			HttpAddr:    "0.0.0.0:6060",
			Rpc:         true,
			RpcAddr:     "0.0.0.0:8433",
			CtrlAddr:    "",
		},
		Batch:    16,
		Backends: make([]specs.Backend, 0),
	}
)
