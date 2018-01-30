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

	service.RegisterModule(&service.TsdbModule{})
	service.RegisterModule(&service.ShardModule{})
	service.RegisterModule(&service.ApiModule{})
	service.RegisterModule(&service.ApiGwModule{})
	service.RegisterModule(&service.TimerModule{})
	service.RegisterModule(&service.TriggerModule{})
	service.RegisterModule(&service.ClientModule{})

}
