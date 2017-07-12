/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/getsentry/raven-go"
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/api/models"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("5")

	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	user := env("MYSQL_TEST_USER", "root")
	pass := env("MYSQL_TEST_PASS", "123456")
	prot := env("MYSQL_TEST_PROT", "tcp")
	addr := env("MYSQL_TEST_ADDR", "localhost:33061")
	dbname := env("MYSQL_TEST_DBNAME", "falcon")
	alarm_dbname := env("MYSQL_TEST_ALARM_DBNAME", "alarm_event")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)
	alarm_dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, alarm_dbname)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	if err := orm.RegisterDataBase("default", "mysql", dsn, 7, 7); err != nil {
		glog.Error("register database default error", err)
		return
	}
	models.Db.Ctrl = orm.NewOrm()

	if err := orm.RegisterDataBase("alarm", "mysql", alarm_dsn, 7, 7); err != nil {
		glog.Error("register database alarm error", err)
		return
	}
	models.Db.Alarm = orm.NewOrm()
	models.Db.Alarm.Using("alarm")
	raven.SetDSN("http://65fdec6391f0446eb08a7018d8ff2c07:6fa70ca785204a98a382606b186865d8@c3-op-sentry.bj/12")

	glog.V(4).Info("alarm", models.Db.Alarm)
}

func TestAlarm(t *testing.T) {
	//testTagTree(t)
	//testPopulate(t)
	//testToken(t)
	running := make(chan struct{})
	defer close(running)

	_start(running)

	time.Sleep(time.Second * 10000000)
}

func TestOK(t *testing.T) {
	//testTagTree(t)
	//testPopulate(t)
	//testToken(t)
	running := make(chan struct{})
	defer close(running)
	var err error
	mop := &MatterOperator{}
	mop.DB, err = orm.GetDB("alarm")
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range mop.GetExpiredMatter() {
		fmt.Println(m.Id)
		pf := P0p1Filter{}
		MarkClosed(pf, m.Id, mop)
		//go mop.GetLatestEventStatus(m.Id)
		fmt.Println()

	}
}
