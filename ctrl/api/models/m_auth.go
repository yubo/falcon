/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "github.com/yubo/falcon"

var (
	allAuths = make(map[string]AuthInterface)
	Auths    = make(map[string]AuthInterface)
)

type Auth struct {
	Method string
	Arg1   string
	Arg2   string
}

type AuthInterface interface {
	Init(conf *falcon.ConfCtrl) error
	Verify(c interface{}) (success bool, uuid string, err error)
	AuthorizeUrl(ctx interface{}) string
	CallBack(ctx interface{}) (uuid string, err error)
}

func RegisterAuth(name string, p AuthInterface) error {
	if _, ok := allAuths[name]; ok {
		return ErrExist
	} else {
		allAuths[name] = p
		return nil
	}
}
