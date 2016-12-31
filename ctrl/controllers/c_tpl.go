/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import "github.com/yubo/falcon/ctrl/models"

func (c *MainController) GetTplAcl() {
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["H1"] = "Access"
	c.Data["H1P"] = "在节点上配置用户(user)、权限模板(role)、权限(token)之间的关系"

	c.TplName = "tpl/acl.tpl"
}

func (c *MainController) GetTplRule() {
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["H1"] = "Rule"
	c.Data["H1P"] = "在节点上配置机器(host)、报警模板、报警(trigger)之间的关系"

	c.TplName = "tpl/rule.tpl"
}
