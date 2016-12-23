/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import "github.com/yubo/falcon/ctrl/models"

func (c *MainController) GetProfile() {
	c.Data["Me"], _ = c.Ctx.Input.GetData("me").(*models.User)
	c.Data["User"] = c.Data["Me"]
	c.Data["Links"] = settingsLinks
	c.Data["CurLink"] = "Profile"
	c.Data["Method"] = "put"
	c.Data["H1"] = "edit Profile"

	c.TplName = "settings/profile.tpl"
}

func (c *MainController) PutProfile() {
	var user models.User

	//settings/profile
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if me.Name != "" && user.Name != "" {
		user.Name = ""
	}

	if u, err := models.UpdateUser(me.Id, &user); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendObj(200, u)
	}
}

func (c *MainController) GetAbout() {
	c.PrepareEnv()
	c.Data["Links"] = settingsLinks
	c.Data["CurLink"] = "About"
	c.Data["H1"] = "About"
	c.TplName = "settings/about.tpl"
}
