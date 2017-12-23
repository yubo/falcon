/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package modules

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/module"
	"github.com/yubo/falcon/ctrl/config"
)

func init() {
	falcon.RegisterModule(&ctrl.Ctrl{}, "ctrl",
		falcon.GetType(config.Ctrl{}))
	ctrl.RegisterModule(&module.ApiModule{})
	ctrl.RegisterModule(&ctrl.EtcdCliModule{})
	ctrl.RegisterModule(&ctrl.StatsModule{})
}
