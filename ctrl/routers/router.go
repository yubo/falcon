// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routers

import (
	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/controllers"
)

func init() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	}

	ns := beego.NewNamespace("/v1.0",
		beego.NSNamespace("/auth", beego.NSInclude(&controllers.AuthController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/host", beego.NSInclude(&controllers.HostController{})),
	)
	beego.AddNamespace(ns)

}
