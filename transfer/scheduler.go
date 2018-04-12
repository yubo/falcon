/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/core"
)

func scheduler(key []byte) int {
	return int(core.Sum64(key) % uint64(falcon.SHARD_NUM))
}
