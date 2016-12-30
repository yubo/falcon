/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	scope_test_init bool
)

func init() {

	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	user := env("MYSQL_TEST_USER", "root")
	pass := env("MYSQL_TEST_PASS", "12341234")
	prot := env("MYSQL_TEST_PROT", "tcp")
	addr := env("MYSQL_TEST_ADDR", "localhost:3306")
	dbname := env("MYSQL_TEST_DBNAME", "falcon_test")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	if err := orm.RegisterDataBase("default", "mysql", dsn, 7, 7); err != nil {
		return
	}
	//orm.Debug = true
	scope_test_init = true
}

func testScopeInitDb(t *testing.T, o orm.Ormer) (err error) {
	tables := []string{
		"tag_host",
		"tag_rel",
		"tag_role_scope",
		"tag_role_user",
		"host",
		"scope",
		"system",
		"tag",
		"role",
		"user",
	}
	t.Log("enter testScopeInitDb")
	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range tables {
		if _, err = o.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return
		}
	}
	o.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	admin := User{
		Name: "admin",
	}
	o.Insert(&admin)

	// init root tree tag
	o.Insert(&Tag{})

	return nil
}

func TestScope(t *testing.T) {
	if !scope_test_init {
		t.Logf("scope orm not inited, skip testing\n")
		return
	}
	o := orm.NewOrm()
	err := testScopeInitDb(t, o)
	if err != nil {
		t.Error("init db failed", err)
	}

	schema, _ := NewTagSchema("a,b,c,d,")
	admin, _ := GetUser(1)
	admin.addTag(&Tag{Name: ""}, schema)

	// tag
	tag_idx := make(map[string]int64)
	tags := []string{
		"a=1",
		"a=1,b=1",
		"a=1,b=2",
		"a=1,b=2,c=1",
		"a=1,b=2,c=2",
	}
	for _, tag := range tags {
		if tag_idx[tag], err = admin.addTag(&Tag{Name: tag}, schema); err != nil {
			t.Error(err)
		}
	}

	// user
	user_idx := make(map[string]int64)
	users := []string{
		"u1",
	}
	for _, user := range users {
		if user_idx[user], err = admin.AddUser(&User{Name: user}); err != nil {
			t.Error(err)
		}
	}

	// role
	role_idx := make(map[string]int64)
	roles := []string{
		"r1",
		"r2",
		"r3",
		"r4",
	}
	for _, role := range roles {
		if role_idx[role], err = admin.AddRole(&Role{Name: role}); err != nil {
			t.Error(err)
		}
	}

	// system
	system_idx := make(map[string]int64)
	systems := []string{
		"s1",
		"s2",
		"s3",
	}
	for _, system := range systems {
		if system_idx[system], err = admin.AddSystem(&System{Name: system}); err != nil {
			t.Error(err)
		}
	}

	// scope
	scope_idx := make(map[string]int64)
	scopes := []string{
		"scope1",
		"scope2",
		"scope3",
		"scope41",
		"scope42",
	}
	for _, scope := range scopes {
		if scope_idx[scope], err = admin.AddScope(&Scope{Name: scope, System_id: system_idx["s1"]}); err != nil {
			t.Error(err)
		}
	}

	// bind user
	binds := [][3]int64{
		{user_idx["u1"], role_idx["r1"], tag_idx["a=1,b=2"]},
		{user_idx["u1"], role_idx["r2"], tag_idx["a=1,b=2"]},
		{user_idx["u1"], role_idx["r3"], tag_idx["a=1,b=2"]},
		{user_idx["u1"], role_idx["r4"], tag_idx["a=1,b=2"]},
	}
	for _, n := range binds {
		if err := admin.BindUserRole(n[0], n[1], n[2]); err != nil {
			t.Error(err)
		}
	}

	// bind scope
	binds = [][3]int64{
		{scope_idx["scope1"], role_idx["r1"], tag_idx["a=1,b=2"]},
		{scope_idx["scope2"], role_idx["r2"], tag_idx["a=1"]},
		{scope_idx["scope3"], role_idx["r3"], tag_idx["a=1,b=2,c=2"]},
		{scope_idx["scope41"], role_idx["r4"], tag_idx["a=1"]},
		{scope_idx["scope42"], role_idx["r4"], tag_idx["a=1,b=2"]},
	}
	for _, n := range binds {
		if err := admin.BindScopeRole(n[0], n[1], n[2]); err != nil {
			t.Error(err)
		}
	}

	// case1~4
	cases := []struct {
		name  string
		uid   int64
		sid   int64
		tid   int64
		want  int64
		wante error
	}{
		//case1
		{name: "case1-1", uid: user_idx["u1"], sid: scope_idx["scope1"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case1-2", uid: user_idx["u1"], sid: scope_idx["scope1"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case1-3", uid: user_idx["u1"], sid: scope_idx["scope1"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case1-4", uid: user_idx["u1"], sid: scope_idx["scope1"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1,b=2"], wante: nil},
		//case2
		{name: "case2-1", uid: user_idx["u1"], sid: scope_idx["scope2"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case2-2", uid: user_idx["u1"], sid: scope_idx["scope2"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case2-3", uid: user_idx["u1"], sid: scope_idx["scope2"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1"], wante: nil},
		{name: "case2-4", uid: user_idx["u1"], sid: scope_idx["scope2"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1"], wante: nil},
		//case3
		{name: "case3-1", uid: user_idx["u1"], sid: scope_idx["scope3"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case3-2", uid: user_idx["u1"], sid: scope_idx["scope3"], tid: tag_idx["a=1,b=2"], want: 0, wante: EACCES},
		{name: "case3-3", uid: user_idx["u1"], sid: scope_idx["scope3"], tid: tag_idx["a=1,b=2,c=1"], want: 0, wante: EACCES},
		{name: "case3-4", uid: user_idx["u1"], sid: scope_idx["scope3"], tid: tag_idx["a=1,b=2,c=2"], want: 0, wante: EACCES},
		//case4
		{name: "case4-1", uid: user_idx["u1"], sid: scope_idx["scope41"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case4-2", uid: user_idx["u1"], sid: scope_idx["scope41"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-3", uid: user_idx["u1"], sid: scope_idx["scope41"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-4", uid: user_idx["u1"], sid: scope_idx["scope41"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1"], wante: nil},
		{name: "case4-5", uid: user_idx["u1"], sid: scope_idx["scope42"], tid: tag_idx["a=1"], want: 0, wante: EACCES},
		{name: "case4-6", uid: user_idx["u1"], sid: scope_idx["scope42"], tid: tag_idx["a=1,b=2"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case4-7", uid: user_idx["u1"], sid: scope_idx["scope42"], tid: tag_idx["a=1,b=2,c=1"], want: tag_idx["a=1,b=2"], wante: nil},
		{name: "case4-8", uid: user_idx["u1"], sid: scope_idx["scope42"], tid: tag_idx["a=1,b=2,c=2"], want: tag_idx["a=1,b=2"], wante: nil},
	}
	for _, c := range cases {
		if got, gote := access(c.uid, c.sid,
			c.tid); got != c.want || gote != c.wante {

			t.Errorf("%s access(%d,%d,%d) = %d, %v; want %d %v",
				c.name, c.uid, c.sid, c.tid,
				got, gote, c.want, c.wante)

		}

	}

}
