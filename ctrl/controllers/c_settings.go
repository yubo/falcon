/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/ctrl/models"
)

func (c *MainController) GetProfile() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_SETTINGS].SubLinks, "Profile")
	c.Data["User"] = c.Data["Me"]
	c.Data["Method"] = "put"
	c.Data["H1"] = "edit Profile"

	c.TplName = "user/edit.tpl"
}

func (c *MainController) GetAboutMe() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_SETTINGS].SubLinks, "About Me")
	c.Data["H1"] = "About Me"
	c.TplName = "settings/about.tpl"
}

func (c *MainController) GetConfig() {
	var err error

	module := c.GetString(":module")
	me, _ := c.Ctx.Input.GetData("me").(*models.User)

	conf, err := me.ConfigGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	c.Data["Config"] = conf.Value.([]models.ConfigEntry)
	c.Data["Module"] = conf.Key
	c.Data["Note"] = conf.Note
	beego.Debug(c.Data["Config"])
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_SETTINGS].SubLinks, module)
	c.TplName = "settings/config.tpl"
}

func (c *MainController) PostConfig() {
	module := c.GetString(":module")
	conf := make(map[string]string)

	json.Unmarshal(c.Ctx.Input.RequestBody, &conf)
	me, _ := c.Ctx.Input.GetData("me").(*models.User)
	if err := me.ConfigSet(module, conf); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, "")
	}
}

func (c *MainController) GetDebug() {
	c.PrepareEnv(headLinks[HEAD_LINK_IDX_SETTINGS].SubLinks, "Debug")
	c.Data["H1"] = "Debug"
	c.TplName = "settings/debug.tpl"
}

func (c *MainController) GetDebugAction() {
	var err error
	var obj interface{}
	action := c.GetString(":action")

	switch action {
	case "populate":
		obj, err = populate()
	case "reset_db":
		obj, err = resetDb()
	case "msg":
		obj, err = msg()
	default:
		err = fmt.Errorf("%s %s", models.ErrUnsupported.Error(), action)
	}

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendObj(200, obj)
	}
}

func msg() (interface{}, error) {
	return "hello world", nil
}

func populate() (interface{}, error) {
	var (
		ret       string
		err       error
		items     []string
		tag_idx   = make(map[string]int64)
		user_idx  = make(map[string]int64)
		role_idx  = make(map[string]int64)
		token_idx = make(map[string]int64)
		host_idx  = make(map[string]int64)
	)
	admin, _ := models.GetUser(1)
	tag_idx["/"] = 1

	// user
	items = []string{
		"test",
	}
	for _, item := range items {
		if user_idx[item], err = admin.AddUser(&models.User{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd user(%s)\n", ret, item)
	}

	// tag
	items = []string{
		"cop=xiaomi",
		"cop=xiaomi,owt=inf",
		"cop=xiaomi,owt=miliao",
		"cop=xiaomi,owt=miliao,pdl=op",
		"cop=xiaomi,owt=miliao,pdl=micloud",
	}
	for _, item := range items {
		if tag_idx[item], err = admin.AddTag(&models.Tag{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd tag(%s)\n", ret, item)
	}

	// host
	items2 := [][2]string{
		{"cop=xiaomi", "mi1.bj"},
		{"cop=xiaomi", "mi2.bj"},
		{"cop=xiaomi", "mi3.bj"},
		{"cop=xiaomi,owt=inf", "inf1.bj"},
		{"cop=xiaomi,owt=inf", "inf2.bj"},
		{"cop=xiaomi,owt=inf", "inf3.bj"},
		{"cop=xiaomi,owt=miliao", "miliao1.bj"},
		{"cop=xiaomi,owt=miliao", "miliao2.bj"},
		{"cop=xiaomi,owt=miliao", "miliao3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud3.bj"},
	}
	for _, item2 := range items2 {
		if host_idx[item2[1]], err = admin.AddHost(&models.Host{Name: item2[1]}); err != nil {
			return nil, err
		}

		if err = admin.CreateTagHost(models.RelTagHost{Tag_id: tag_idx[item2[0]],
			Host_id: host_idx[item2[1]]}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd host(%s, %s)\n", ret, item2[1], item2[0])
	}

	// role
	items = []string{
		"adm",
		"sre",
		"dev",
		"usr",
	}
	for _, item := range items {
		if role_idx[item], err = admin.AddRole(&models.Role{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd role(%s)\n", ret, item)
	}

	// token
	items = []string{
		models.SYS_W_SCOPE,
		models.SYS_R_SCOPE,
		models.SYS_B_SCOPE,
		models.SYS_O_SCOPE,
		models.SYS_A_SCOPE,
	}
	for _, item := range items {
		if token_idx[item], err = admin.AddToken(&models.Token{Name: item}); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sadd token(%s)\n", ret, item)
	}

	// bind user
	binds := [][3]string{
		{"cop=xiaomi,owt=miliao", "test", "adm"},
		{"cop=xiaomi,owt=miliao", "test", "sre"},
		{"cop=xiaomi,owt=miliao", "test", "dev"},
		{"cop=xiaomi,owt=miliao", "test", "usr"},
	}
	for _, s := range binds {
		if err := admin.BindAclUser(tag_idx[s[0]], role_idx[s[2]],
			user_idx[s[1]]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) user(%s) role(%s)\n",
			ret, s[0], s[1], s[2])
	}

	// bind token
	binds = [][3]string{
		{models.SYS_W_SCOPE, "adm", "/"},
		{models.SYS_R_SCOPE, "adm", "/"},
		{models.SYS_W_SCOPE, "sre", "/"},
		{models.SYS_R_SCOPE, "sre", "/"},
		{models.SYS_R_SCOPE, "dev", "/"},
		{models.SYS_R_SCOPE, "usr", "/"},
		{models.SYS_W_SCOPE, "adm", "cop=xiaomi,owt=miliao"},
		{models.SYS_R_SCOPE, "sre", "cop=xiaomi"},
		{models.SYS_B_SCOPE, "dev", "cop=xiaomi,owt=miliao,pdl=op"},
		{models.SYS_O_SCOPE, "usr", "cop=xiaomi"},
		{models.SYS_A_SCOPE, "usr", "cop=xiaomi,owt=miliao"},
	}
	for _, s := range binds {
		if err := admin.BindAclToken(tag_idx[s[2]], role_idx[s[1]],
			token_idx[s[0]]); err != nil {
			return nil, err
		}
		ret = fmt.Sprintf("%sbind tag(%s) token(%s) role(%s)\n",
			ret, s[1], s[2], s[0])
	}

	return ret, nil
}

func resetDb() (interface{}, error) {
	var err error
	tables := []string{
		"host",
		"kv",
		"log",
		"role",
		"tag",
		"tag_host",
		"tag_rel",
		"token",
		"tpl_rel",
		"user",
	}
	o := orm.NewOrm()

	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range tables {
		if _, err = o.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return nil, err
		}
	}
	o.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	o.Insert(&models.User{Name: "system"})

	// init root tree tag
	o.Insert(&models.Tag{Name: ""})

	// reset cache
	models.CacheInit()

	return "reset db done", nil
}
