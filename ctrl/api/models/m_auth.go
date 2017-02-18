/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "github.com/yubo/falcon"

var (
	allAuths map[string]AuthInterface
	Auths    map[string]AuthInterface
)

type Auth struct {
	Method string
	Arg1   string
	Arg2   string
}

type AuthModule struct {
	Name       string
	Prestarted bool
}

func (p *AuthModule) GetName() string {
	return p.Name
}

func (p *AuthModule) AuthorizeUrl(ctx interface{}) string {
	return ""
}

func (p *AuthModule) CallBack(ctx interface{}) (uuid string, err error) {
	err = ErrNoModule
	return
}

func (p *AuthModule) Verify(c interface{}) (bool, string, error) {
	return false, "", EPERM
}

func (p *AuthModule) PreStart(conf falcon.ConfCtrl) error {
	if p.Prestarted {
		return ErrRePreStart
	}
	p.Prestarted = true
	return nil
}

type AuthInterface interface {
	GetName() string
	Verify(c interface{}) (success bool, uuid string, err error)
	CallBack(ctx interface{}) (uuid string, err error)
	PreStart(conf falcon.ConfCtrl) error
	AuthorizeUrl(ctx interface{}) string
}

func RegisterAuth(p AuthInterface) error {
	if _, ok := allAuths[p.GetName()]; ok {
		return ErrExist
	} else {
		allAuths[p.GetName()] = p
		return nil
	}
}
