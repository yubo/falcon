/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/transfer"
	"github.com/yubo/falcon/transfer/config"
)

func init() {
	falcon.RegisterModule(&transfer.Transfer{}, falcon.GetType(config.ConfTransfer{}), "transfer")
	transfer.RegisterModule(&transfer.StatsModule{})
	transfer.RegisterModule(&transfer.GrpcModule{})
	transfer.RegisterModule(&transfer.HttpModule{})
	transfer.RegisterModule(&transfer.BackendModule{})

}
