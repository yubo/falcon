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

	// must include
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yubo/falcon/ctrl/api/models/auth"
	_ "github.com/yubo/falcon/ctrl/api/models/session"

	// plugin modules
	_ "github.com/yubo/falcon/ctrl/plugin/demo"
	_ "github.com/yubo/falcon/ctrl/plugin/rate_limits"
	//_ "github.com/yubo/falcon/ctrl/plugin/mi"
)

func init() {

	falcon.RegisterModule(&ctrl.Ctrl{}, "ctrl",
		falcon.GetType(config.Ctrl{}))

	ctrl.RegisterModule(&ctrl.ClientModule{})
	ctrl.RegisterModule(&ctrl.EtcdCliModule{})
	ctrl.RegisterModule(&ctrl.OrmModule{})
	ctrl.RegisterModule(&module.ModelsModule{})
	ctrl.RegisterModule(&module.ApiModule{})
}
