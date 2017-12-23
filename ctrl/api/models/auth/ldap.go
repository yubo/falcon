/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"crypto/tls"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/controllers"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/config"
	"gopkg.in/ldap.v2"
)

type ldapAuth struct {
	addr    string
	baseDN  string
	bindDN  string
	bindPwd string
	filter  string
	tls     bool
}

const (
	LDAP_NAME = "ldap"
)

func init() {
	models.RegisterAuth(LDAP_NAME, &ldapAuth{})
}

func (p *ldapAuth) Init(conf *config.ConfCtrl) error {
	p.addr = conf.Ctrl.Str(ctrl.C_LDAP_ADDR)
	p.baseDN = conf.Ctrl.Str(ctrl.C_LDAP_BASE_DN)
	p.bindDN = conf.Ctrl.Str(ctrl.C_LDAP_BIND_DN)
	p.bindPwd = conf.Ctrl.Str(ctrl.C_LDAP_BIND_PWD)
	p.filter = conf.Ctrl.Str(ctrl.C_LDAP_FILTER)
	return nil
}

func (p *ldapAuth) Verify(_c interface{}) (bool, string, error) {
	c := _c.(*controllers.AuthController)
	username := c.GetString("username")
	password := c.GetString("password")

	if beego.BConfig.RunMode == "dev" && username[:4] == "test" {
		return true, username, nil
	}

	success, uuid, err := ldapUserAuthentication(p.addr, p.baseDN, p.filter,
		username, password,
		p.bindDN, p.bindPwd, p.tls)
	if success {
		uuid = fmt.Sprintf("%s@%s", uuid, LDAP_NAME)
	}

	return success, uuid, err
}

func (p *ldapAuth) AuthorizeUrl(c interface{}) string {
	return ""
}

func (p *ldapAuth) LoginCb(c interface{}) (uuid string, err error) {
	return "", falcon.EPERM
}

func (p *ldapAuth) LogoutCb(c interface{}) {
}

func ldapUserAuthentication(addr, baseDN, filter, username, password, bindusername, bindpassword string, TLS bool) (success bool, userDN string, err error) {
	var (
		sr *ldap.SearchResult
	)

	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()

	// Reconnect with TLS
	if TLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			beego.Warning(err)
		}
	}

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		return
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filter, username),
		[]string{"dn"},
		nil,
	)

	sr, err = l.Search(searchRequest)
	if err != nil {
		return
	}

	if len(sr.Entries) != 1 {
		return
	}

	userDN = sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return
	}
	return true, userDN, err

}
