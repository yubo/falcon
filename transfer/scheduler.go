/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import "github.com/yubo/falcon"

var (
	shardNum int
)

type SchedulerModule struct {
	BaseModule
}

func (p *SchedulerModule) start(transfer *Transfer) (err error) {
	shardNum = falcon.SHARD_NUM
	return nil
}

func scheduler(key []byte) int {
	return int(falcon.Sum64(key) % uint64(shardNum))
}
