/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/config"
)

type weappAuth struct {
}

const (
	WEAPP_NAME = "weapp"
)

func init() {
	models.RegisterAuth(WEAPP_NAME, &weappAuth{})
}

func (p *weappAuth) Init(conf *config.Ctrl) error {
	return nil
}

func (p *weappAuth) Verify(_c interface{}) (bool, string, error) {
	return false, "", nil
}

func (p *weappAuth) AuthorizeUrl(c interface{}) string {
	return ""
}

func (p *weappAuth) LoginCb(c interface{}) (uuid string, err error) {
	return "", falcon.EPERM
}

func (p *weappAuth) LogoutCb(c interface{}) {
}
