/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func testSettingsInitDb(t *testing.T, o orm.Ormer) (err error) {
	t.Log("enter testSettingsInitDb")
	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range dbTables {
		if _, err = o.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
			return
		}
	}
	o.Raw("SET FOREIGN_KEY_CHECKS=1").Exec()

	// init admin
	o.Insert(&User{Name: "admin"})

	// init root tree tag
	o.Insert(&Tag{})

	return nil
}

func TestPopulate(t *testing.T) {

	if !test_db_init {
		t.Logf("test db not inited, skip test populate\n")
		return
	}

	o := orm.NewOrm()
	sys, _ := GetUser(1, o)
	op := &Operator{
		O:    o,
		User: sys,
	}

	err := testSettingsInitDb(t, op.O)
	if err != nil {
		t.Error("init db failed", err)
	}

	if _, err := op.Populate(); err != nil {
		t.Error(err)
	}
}
