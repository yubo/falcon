/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package modules

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend"
	"github.com/yubo/falcon/backend/config"
)

func init() {
	falcon.RegisterModule(&backend.Backend{}, "backend",
		falcon.GetType(config.ConfBackend{}))

	// cache should early register(init cache data)
	backend.RegisterModule(&backend.StorageModule{})
	backend.RegisterModule(&backend.CacheModule{})
	backend.RegisterModule(&backend.GrpcModule{})
	backend.RegisterModule(&backend.HttpModule{})
	backend.RegisterModule(&backend.IndexModule{})
	backend.RegisterModule(&backend.StatsModule{})
	backend.RegisterModule(&backend.TimerModule{})

}
