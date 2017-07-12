/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/agent/config"
	"github.com/yubo/falcon/utils"
	//_ "github.com/yubo/falcon/ctrl/api/models/plugin/demo"
	//_ "github.com/yubo/falcon/ctrl/api/models/plugin/alarm"
)

func init() {
	falcon.RegisterModule(&agent.Agent{}, utils.GetType(config.ConfAgent{}), "agent")
	agent.RegisterModule(&agent.StatsModule{})
	agent.RegisterModule(&agent.CollectModule{})
	agent.RegisterModule(&agent.UpstreamModule{})
	agent.RegisterModule(&agent.HttpModule{})
}
