/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

var (
	rels = map[string]bool{
		"tag_host":         true,
		"tag_role_user":    true,
		"tag_role_token":   true,
		"tag_rule_trigger": true,
	}
)

func (c *MainController) GetTagHost() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Host")
	c.TplName = "rel/tag_host.tpl"
}

func (c *MainController) GetTagRoleUser() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Role User")
	c.TplName = "rel/tag_role_user.tpl"
}

func (c *MainController) GetTagRoleToken() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Role Token")
	c.TplName = "rel/tag_role_token.tpl"
}

func (c *MainController) GetTagRuleTrigger() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_REL].SubLinks, "Tag Rule Trigger")
	c.TplName = "rel/tag_rule_trigger.tpl"
}
