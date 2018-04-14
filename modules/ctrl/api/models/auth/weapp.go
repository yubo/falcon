/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/ctrl/api/models"
)

type weappAuth struct {
}

const (
	WEAPP_NAME = "weapp"
)

func init() {
	models.RegisterAuth(WEAPP_NAME, &weappAuth{})
}

func (p *weappAuth) Init(conf *core.Configer) error {
	return nil
}

func (p *weappAuth) Verify(_c interface{}) (bool, string, error) {
	return false, "", nil
}

func (p *weappAuth) AuthorizeUrl(c interface{}) string {
	return ""
}

func (p *weappAuth) LoginCb(c interface{}) (uuid string, err error) {
	return "", core.EPERM
}

func (p *weappAuth) LogoutCb(c interface{}) {
}
