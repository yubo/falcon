/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yubo/falcon/ctrl/models/auth"
	_ "github.com/yubo/falcon/ctrl/routers"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	dsn := beego.AppConfig.String("mysqldsn")
	maxIdle, _ := beego.AppConfig.Int("mysqlmaxidle")
	maxConn, _ := beego.AppConfig.Int("mysqlmaxconn")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn, maxIdle, maxConn)

	beego.Run()
}
