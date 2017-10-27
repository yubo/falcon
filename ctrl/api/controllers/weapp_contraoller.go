/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package controllers

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/api/models"
)

// Operations about weixin app
type WeappController struct {
	BaseController
}

// @Title login (acl: pub)
// @Description wexin app login api
// @Success 200 string success
// @Failure 400 string error
// @router /login [get]
func (c *WeappController) Login() {
	code := c.Ctx.Input.Header(models.WX_HEADER_CODE)
	encrypt_data := c.Ctx.Input.Header(models.WX_HEADER_ENCRYPTED_DATA)
	iv := c.Ctx.Input.Header(models.WX_HEADER_IV)

	glog.V(4).Infof("code %s encrypt_data %s iv %s", code, encrypt_data, iv)

	sess, err := models.WeappLogin(code, encrypt_data, iv)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, map[string]interface{}{
			models.WX_SESSION_MAGIC_ID: "1",
			"session": map[string]interface{}{
				"id":     sess.Wxopenid,
				"openid": sess.Wxopenid,
				"skey":   sess.Key,
				"user":   sess.User,
			},
		})
	}

}

// @Title openid (acl: pub)
// @Description wexin app api
// @Param	code	query   string     true       "code from weapp wx.login()"
// @Success 200 string success
// @Failure 400 string error
// @router /openid [get]
func (c *WeappController) Openid() {
	ret, err := models.WeappOpenid(c.GetString("code"))
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, map[string]string{"openid": ret})
	}
}

// @Title create get auth task qrcode
// @Description get auth task qrcode
// @Success 200 {object} models.QrTask qr code image(encode by base64)
// @Failure 400 string error
// @router /authqr [get]
func (c *WeappController) Authqr() {
	ret, err := models.WeappAuthQr()
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title ack auth to falcon request(acl: weapp login && binded)
// @Description bind weapp to cur user
// @Param	key	query   string     true       "task key"
// @Success 200 {object} models.User bind to falcon user
// @Failure 400 string error
// @router /taskack [get]
func (c *WeappController) Taskack() {

	sess, err := models.WeappGetSession(c.Ctx.Input.Header(models.WX_HEADER_SKEY))
	if err != nil {
		c.SendMsg(401, map[string]string{models.WX_SESSION_MAGIC_ID: "1"})
		return
	}

	ret, err := models.WeappTaskAck(c.GetString("key"), sess)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create bind weapp to falcon user task(acl: falcon session login)
// @Description bind weapp to cur user
// @Success 200 {object} models.QrTask qr code image(encode by base64)
// @Failure 400 string error
// @router /bindqr [get]
func (c *WeappController) Bindqr() {
	op, ok := c.Ctx.Input.GetData("op").(*models.Operator)
	if !ok || op.User == nil {
		c.SendMsg(401, "Unauthorized")
		return
	}

	if op.User.Muid != 0 {
		c.SendMsg(400, fmt.Sprintf("%s already bind to uid(%d)",
			op.User.Name, op.User.Muid))
		return
	}

	ret, err := models.WeappBindQr(op.User.Id)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title get task states
// @Description get task states
// @Param	key	query   string     true       "task key"
// @Success 200 {string} states
// @Failure 400 string error
// @router /task [get]
func (c *WeappController) Task() {
	ret, err := models.GetTask(c.GetString("key"), c)
	if err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title testRequest (acl: weapp skey)
// @Description wexin app api
// @Success 200 string success
// @Failure 400 string error
// @router /testRequest [get]
func (c *WeappController) TestRequest() {
	c.SendMsg(200, "")
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
