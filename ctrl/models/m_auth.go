/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

var (
	AuthMap map[string]AuthInterface
	Auths   []AuthInterface
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

func (p *AuthModule) CallBack(c interface{}) {
}

func (p *AuthModule) Verify(c interface{}) (bool, string, error) {
	return false, "", EPERM
}

func (p *AuthModule) PreStart() error {
	if p.Prestarted {
		return ErrRePreStart
	}
	p.Prestarted = true
	return nil
}

type AuthInterface interface {
	GetName() string
	Verify(c interface{}) (success bool, uuid string, err error)
	CallBack(c interface{})
	PreStart() error
}

func RegisterAuth(p AuthInterface) error {
	if _, ok := AuthMap[p.GetName()]; ok {
		return ErrExist
	} else {
		AuthMap[p.GetName()] = p
		return nil
	}
}
