/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
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
	o.Raw("insert user (name) value ('admin')").Exec()

	// init root tree tag
	o.Raw("insert tag (name) value ('')").Exec()

	return nil
}

func testPopulate(t *testing.T) {

	if !test_db_init {
		t.Logf("test db not inited, skip test populate\n")
		return
	}

	o := orm.NewOrm()
	sys, _ := GetUser(1, o)
	op := &Operator{
		O:     o,
		User:  sys,
		Token: SYS_F_A_TOKEN | SYS_F_O_TOKEN,
	}

	err := testSettingsInitDb(t, op.O)
	if err != nil {
		t.Error("init db failed", err)
	}

	if _, err := op.populate(); err != nil {
		t.Error(err)
	}
}
