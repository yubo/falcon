/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

const (
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	DefaultLb = Lb{
		Debug:       0,
		Http:        true,
		HttpAddr:    "0.0.0.0:6060",
		Rpc:         true,
		RpcAddr:     "0.0.0.0:8433",
		Replicas:    500,
		Concurrency: 2,
		Backends:    make([]Backend, 0),
	}
)
