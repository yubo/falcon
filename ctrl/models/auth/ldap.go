/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"crypto/tls"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/controllers"
	"github.com/yubo/falcon/ctrl/models"
	"gopkg.in/ldap.v2"
)

type ldapAuth struct {
	models.AuthModule
	addr    string
	baseDN  string
	bindDN  string
	bindPwd string
	filter  string
	tls     bool
}

var (
	_ldap = &ldapAuth{
		AuthModule: models.AuthModule{Name: "ldap"},
		addr:       "localhost:389",
		baseDN:     "dc=yubo,dc=org",
		bindDN:     "cn=admin,dc=yubo,dc=org",
		bindPwd:    "12341234",
		filter:     "(&(objectClass=posixAccount)(cn=%s))",
		tls:        false,
	}
)

func init() {
	models.RegisterAuth(_ldap)
}

func (p *ldapAuth) PreStart() error {
	if p.AuthModule.Prestarted {
		return models.ErrRePreStart
	}
	p.AuthModule.Prestarted = true
	p.addr = beego.AppConfig.String("ldapaddr")
	p.baseDN = beego.AppConfig.String("ldapbasedn")
	p.bindDN = beego.AppConfig.String("ldapbinddn")
	p.bindPwd = beego.AppConfig.String("ldapbindpwd")
	p.filter = beego.AppConfig.String("ldapfilter")
	p.tls, _ = beego.AppConfig.Bool("ldaptls")
	return nil
}

func (p *ldapAuth) Verify(_c interface{}) (bool, string, error) {
	c := _c.(*controllers.AuthController)

	beego.Debug("ldap verify user", c.GetString("username"), "password", c.GetString("password"))

	if beego.BConfig.RunMode == "dev" && c.GetString("username") == "test" {
		return true, "test", nil
	}

	success, uuid, err := ldapUserAuthentication(p.addr, p.baseDN, p.filter,
		c.GetString("username"), c.GetString("password"),
		p.bindDN, p.bindPwd, p.tls)
	if success {
		uuid = fmt.Sprintf("%s@%s", uuid, p.Name)
	}

	return success, uuid, err

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
