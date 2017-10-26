/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/api/models"
)

const (
	WX_HEADER_CODE           = "X-WX-Code"
	WX_HEADER_ENCRYPTED_DATA = "X-WX-Encrypted-Data"
	WX_HEADER_IV             = "X-WX-IV"
	WX_SESSION_MAGIC_ID      = "F2C224D4-2BCE-4C64-AF9F-A6D872000D1A"
)

// Operations about weixin app
type WeappController struct {
	BaseController
}

// @Title login
// @Description wexin app login api
// @Success 200 string success
// @Failure 400 string error
// @router /login [get]
func (c *WeappController) Login() {
	code := c.Ctx.Input.Header(WX_HEADER_CODE)
	encrypt_data := c.Ctx.Input.Header(WX_HEADER_ENCRYPTED_DATA)
	iv := c.Ctx.Input.Header(WX_HEADER_IV)

	glog.V(4).Infof("code %s encrypt_data %s iv %s", code, encrypt_data, iv)

	ret, err := models.WeappLogin(code, encrypt_data, iv)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		ret[WX_SESSION_MAGIC_ID] = "1"
		c.SendMsg(200, ret)
	}

}

// @Title testRequest
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /testRequest [get]
func (c *WeappController) TestRequest() {
	c.SendMsg(200, "")
}

// @Title openid
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /openid [get]
func (c *WeappController) Openid() {
	ret, err := models.WeappOpenid(c.GetString("code"))
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title tunnel
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /tunnel [get]
func (c *WeappController) Tunnel() {
	c.SendMsg(200, "")
}

// @Title templateMessage
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /templateMessage [get]
func (c *WeappController) TemplateMessage() {
	c.SendMsg(200, "")
}

// @Title upload
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /upload [post]
func (c *WeappController) Upload() {
	c.SendMsg(200, "")
}
