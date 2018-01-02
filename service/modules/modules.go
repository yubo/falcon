/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package modules

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
	"github.com/yubo/falcon/service/config"
)

func init() {
	falcon.RegisterModule(&service.Service{}, "service",
		falcon.GetType(config.Service{}))

	// cache should early register(init cache data)
	// service.RegisterModule(&service.StorageModule{})
	service.RegisterModule(&service.ShardModule{})
	service.RegisterModule(&service.ApiModule{})
	// service.RegisterModule(&service.HttpModule{})
	service.RegisterModule(&service.IndexModule{})
	service.RegisterModule(&service.StatsModule{})
	service.RegisterModule(&service.TimerModule{})

}
