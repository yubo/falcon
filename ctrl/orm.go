/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"database/sql"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

var (
	Orm *OrmModule
)

type OrmModule struct {
	Ctrl   orm.Ormer
	CtrlDb *sql.DB
	Idx    orm.Ormer
	Alarm  orm.Ormer
}

func (p *OrmModule) PreStart(ctrl *Ctrl) error {
	Orm = p
	return nil
}

func (p *OrmModule) Start(ctrl *Ctrl) (err error) {
	conf := &ctrl.Conf.Ctrl
	dbMaxConn, _ := conf.Int(C_DB_MAX_CONN)
	dbMaxIdle, _ := conf.Int(C_DB_MAX_IDLE)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", conf.Str(C_DSN), dbMaxIdle, dbMaxConn)
	if err != nil {
		return err
	}
	p.Ctrl = orm.NewOrm()
	p.CtrlDb, err = orm.GetDB()
	if err != nil {
		return err
	}

	/*
		if p.Ctrl, p.CtrlDb, err = falcon.NewOrm("ctrl_falcon",
			conf.Str(C_DSN), dbMaxIdle,
			dbMaxConn); err != nil {
			return err
		}
	*/
	if p.Idx, _, err = falcon.NewOrm("ctrl_index",
		conf.Str(C_IDX_DSN), dbMaxIdle,
		dbMaxConn); err != nil {
		return err
	}
	if p.Alarm, _, err = falcon.NewOrm("ctrl_alarm",
		conf.Str(C_ALARM_DSN), dbMaxIdle,
		dbMaxConn); err != nil {
		return err
	}
	return nil
}

func (p *OrmModule) Stop(ctrl *Ctrl) error {
	return nil
}

func (p *OrmModule) Reload(ctrl *Ctrl) error {
	return nil
}
