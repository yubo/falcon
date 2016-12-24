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
		AuthModule: models.AuthModule{
			Name: "ldap",
			Html: `
      <form class="form-signin" action="/login" method="post">
        <label for="inputUsername" class="sr-only">Username</label>
        <input name="username" type="text" class="form-control" placeholder="Username " required autofocus>
        <label for="inputPassword" class="sr-only">Password</label>
        <input name="password" type="password" class="form-control" placeholder="Password" required>
        <div class="checkbox">
          <label>
            <input name="remember-me" type="checkbox" value="remember-me"> Remember me
          </label>
        </div>
	<input type="hidden" name="method" value="ldap">
        <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
      </form>
`,
		},
		addr:    "localhost:389",
		baseDN:  "dc=yubo,dc=org",
		bindDN:  "cn=admin,dc=yubo,dc=org",
		bindPwd: "12341234",
		filter:  "(&(objectClass=posixAccount)(cn=%s))",
		tls:     false,
	}
)

func init() {
	models.RegisterAuth(_ldap)
}

func (p *ldapAuth) Verify(_c interface{}) (bool, string, error) {
	c := _c.(*controllers.AuthController)

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
