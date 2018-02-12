/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package modules

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/agent/config"
	_ "github.com/yubo/falcon/agent/plugin"
)

func init() {
	falcon.RegisterModule(&agent.Agent{}, "agent",
		falcon.GetType(config.Agent{}))

	agent.RegisterModule(&agent.CollectModule{})
	agent.RegisterModule(&agent.ClientModule{})
	agent.RegisterModule(&agent.ApiModule{})
	agent.RegisterModule(&agent.ApiGwModule{})
	agent.RegisterModule(&agent.MsgModule{})
}
