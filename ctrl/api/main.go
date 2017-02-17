/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yubo/falcon/ctrl/api/models"
	_ "github.com/yubo/falcon/ctrl/api/models/auth"
	//_ "github.com/yubo/falcon/ctrl/api/models/plugin/demo"
	_ "github.com/yubo/falcon/ctrl/api/routers"
)

var (
	conf = map[string]string{
		"mysqldsn":     "root:123456@tcp(localhost:3306)/yubo_falcon?loc=Local&charset=utf8",
		"mysqlmaxidle": "30",
		"mysqlmaxconn": "30",
	}
)

func main() {
	models.ConfInit(conf)
	beego.Run()
}
