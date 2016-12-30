/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"

	"github.com/yubo/falcon/ctrl/models"
)

func (c *MainController) GetProfile() {
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["User"] = c.Data["Me"]
	c.Data["Links"] = settingsLinks
	c.Data["CurLink"] = "Profile"
	c.Data["Method"] = "put"
	c.Data["H1"] = "edit Profile"

	c.TplName = "settings/profile.tpl"
}

func (c *MainController) GetAboutMe() {
	c.PrepareEnv()
	c.Data["Links"] = settingsLinks
	c.Data["CurLink"] = "AboutMe"
	c.Data["H1"] = "About Me"
	c.TplName = "settings/about.tpl"
}

func (c *MainController) GetConfigGlobal() {
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	c.PrepareEnv()
	c.Data["Links"] = settingsLinks
	c.Data["CurLink"] = "Global"
	c.Data["Moudle"] = "global"
	c.Data["Config"], _ = me.ConfigGet("global")
	c.TplName = "settings/config.tpl"
}

func (c *MainController) PostConfigGlobal() {
	conf := make(map[string]string)

	json.Unmarshal(c.Ctx.Input.RequestBody, &conf)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if err := me.ConfigSet("global", conf); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, "")
	}
}
