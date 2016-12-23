/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"fmt"
	"testing"
)

func TestVerify(t *testing.T) {

	addr := "localhost:389"
	baseDN := "dc=yubo,dc=org"
	username := "yubo"
	password := "12341234"
	bindusername := "cn=admin,dc=yubo,dc=org"
	bindpassword := "12341234"
	filter := "(&(objectClass=posixAccount)(cn=%s))"

	success, userDN, err := ldapUserAuthentication(addr, baseDN, filter, username, password, bindusername, bindpassword, false)
	fmt.Println("success:", success)
	fmt.Println("userDN:", userDN)
	fmt.Println("err:", err)
}
