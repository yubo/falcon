/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"encoding/json"
	"testing"

	"github.com/astaxie/beego/orm"
)

func testTagInitDb(t *testing.T, o orm.Ormer) (err error) {
	tables := []string{
		"tag_host",
		"tag_rel",
		"tag",
	}
	t.Log("enter testTokenInitDb")
	o.Raw("SET FOREIGN_KEY_CHECKS=0").Exec()
	for _, table := range tables {
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

func TestTagTree(t *testing.T) {
	if !test_db_init {
		t.Logf("token orm not inited, skip test tag tree\n")
		return
	}
	t.Logf("token orm inited,  test tag tree\n")
	o := orm.NewOrm()
	err := testTagInitDb(t, o)
	if err != nil {
		t.Error("init db failed", err)
	}
	schema, _ := NewTagSchema("a,b,c,d,")
	admin, _ := GetUser(1)

	// tag
	items := []string{
		"a=1",
		"a=1,b=1",
		"a=1,b=2",
		"a=1,b=2,c=1",
		"a=1,b=2,c=2",
	}
	for _, item := range items {
		if _, err = admin.addTag(&Tag{Name: item}, schema); err != nil {
			t.Error(err)
		}
	}

	tree, err := admin.GetTree()
	if err != nil {
		t.Error(err)
	} else {
		json.Marshal(tree)
	}
}
