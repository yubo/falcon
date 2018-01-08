/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package modules

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm"
	"github.com/yubo/falcon/alarm/config"
)

func init() {
	falcon.RegisterModule(&alarm.Alarm{}, "alarm",
		falcon.GetType(config.Alarm{}))
	alarm.RegisterModule(&alarm.StatsModule{})
	alarm.RegisterModule(&alarm.ApiModule{})
	alarm.RegisterModule(&alarm.ApiGwModule{})
	alarm.RegisterModule(&alarm.ClientModule{})

}
