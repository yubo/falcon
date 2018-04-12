/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/lib/core"
)

type dbModule struct {
}

func (p *dbModule) PreStart(ctrl *Ctrl) error {
	return nil
}

func (p *dbModule) Start(ctrl *Ctrl) (err error) {
	conf := ctrl.Conf
	db := &models.CtrlDb{}

	dbMaxConn := conf.DbMaxConn
	dbMaxIdle := conf.DbMaxIdle

	orm.RegisterDriver("mysql", orm.DRMySQL)
	err = orm.RegisterDataBase("default", "mysql", conf.Dsn, dbMaxIdle, dbMaxConn)
	if err != nil {
		return err
	}
	db.Ctrl = orm.NewOrm()
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
	if db.Idx, _, err = core.NewOrm("ctrl_index",
		conf.IdxDsn, dbMaxIdle,
		dbMaxConn); err != nil {
		return err
	}
	if db.Alarm, _, err = core.NewOrm("ctrl_alarm",
		conf.AlarmDsn, dbMaxIdle,
		dbMaxConn); err != nil {
		return err
	}
	ctrl.db = db
	return nil
}

func (p *dbModule) Stop(ctrl *Ctrl) error {
	return nil
}

func (p *dbModule) Reload(ctrl *Ctrl) error {
	return nil
}
