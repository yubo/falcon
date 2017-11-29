/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	test_db_init bool
)

const (
	sqlPath = "/tmp/falcon.sql"
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
	test_db_init = true

}

func testDbReset(t *testing.T) {
	op := &Operator{O: orm.NewOrm()}
	for _, table := range dbTables {
		if _, err := op.O.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			t.Fatal(err)
		}
	}
}

func testOrm(t *testing.T) {
	if !test_db_init {
		return
	}

	t.Logf("=== run")
	testDbReset(t)
	op := &Operator{O: orm.NewOrm()}
	e := &Trigger{
		//UNIQUE INDEX `index_triggers_tag_name` (`tag_id`, `name`)
		TagId: 1,
		Name:  "test",
	}

	id, err := op.SqlInsert("insert triggers (tag_id, name) values (?, ?)", e.TagId, e.Name)
	if err != nil {
		t.Fatal(err)
	}

	// test un
	if _, err := op.SqlInsert("insert triggers (tag_id, cname) values (?, ?)", e.TagId, e.Name); err == nil {
		t.Fatalf("insert row again got nil, want err\n")
	}

	// test un is null; ugly
	if _, err := op.SqlInsert("insert triggers (tag_id) values (?)", e.TagId); err != nil {
		t.Fatalf("insert null un row got %v, want nil\n", err)
	}
	if _, err := op.SqlInsert("insert triggers (name) values (?)", e.Name); err != nil {
		t.Fatalf("insert null un row got %v, want nil\n", err)
	}

	if err := op.SqlRow(e, "select id, name, tag_id from triggers where id = ?", id); err != nil {
		t.Fatal(err)
	}

	if n, err := op.SqlExec("delete from triggers where id = ?", id); err != nil {
		t.Fatalf("delete row %d got %d, %v want 1, nil\n", id, n, err)
	}

	if n, err := op.SqlExec("delete from triggers where id = ?", id); n != 0 || err != err {
		t.Fatalf("delete row %d again, got %d, %v want 0, nil", id, n, err)
	}

	if err := op.SqlRow(e, "select id, name, tag_id from triggers where id = ?", id); err != orm.ErrNoRows {
		t.Fatalf("select no exists row, got %v, want %v", err, orm.ErrNoRows)
	}

}

func testPopulate(t *testing.T) {

	if !test_db_init {
		t.Logf("test db not inited, skip test populate\n")
		return
	}

	op := &Operator{O: orm.NewOrm()}
	op.resetDb(sqlPath)

	o := orm.NewOrm()
	sys, _ := GetUser(1, o)
	op = &Operator{
		O:     o,
		User:  sys,
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN,
	}

	if _, err := op.populate(); err != nil {
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	op := &Operator{O: orm.NewOrm()}
	op.resetDb(sqlPath)

	testPopulate(t)

	//testOrm(t)
	//testTag(t)
	//testSettings(t)
	//testToken(t)
	//testMi(t)
	//testUtils(t)

}
