/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

// add "PluginAlarm true;"config file at ctrl section to enable this plugin
import (
	"time"

	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/api/models"
	"fmt"
)

const (
	MODULE_NAME = "\x1B[32m[CTRL_PLUGIN_ALARM]\x1B[0m "
)

var (
	running  chan struct{}
	enable   bool
	db_alarm orm.Ormer
	db_ctrl  orm.Ormer
)

type Event struct {
	Id       int64
	Pdl      string
	ActionId int64
	DbTs     time.Time
	Endpoint string
	/*
		EventId      string
		Status       string
		Endpoint     string
		Hostname     string
		IsCustomize  int
		Idc          string
		EventTs  string
		CurrentStep  int
		Metric       string
		PushedTags   string
		LeftVal      float64
		Fun          string
		Op           string
		RightVal     float64
		MaxStep      int64
		Priority     int64
		Note         string
		TplId        int64
		ExpressionId int64
		StrategyId   int64
		LastValues   string
	*/
}

func init() {
	running = make(chan struct{})
	ctrl.RegisterPrestart(start)
}

func start(conf *falcon.ConfCtrl) error {
	enable, _ = conf.Ctrl.Bool("PluginAlarm")
	if !enable {
		glog.V(3).Info(MODULE_NAME + " not enabled")
		return nil
	}

	close(running)
	running = make(chan struct{})

	return _start(running)
}

func _start(done chan struct{}) error {

	ticker := time.NewTicker(time.Second * 5)
	tk := time.NewTicker(time.Second * 30)

	go func() {
		db_ctrl = models.Db.Ctrl
		db_alarm = models.Db.Alarm
		id := getLastEventId()
		for {
			select {
			case <-ticker.C:
				glog.V(4).Info("alarm")

				n := syncEvent(id, 100)
				id += n

				for n == 100 {
					n = syncEvent(id, 100)
					id += n
				}
				glog.V(4).Info(strconv.FormatInt(n, 10))
			case <-tk.C:
				var err error
				mop := &MatterOperator{}
				mop.DB, err = orm.GetDB("alarm")
				if err != nil {
					fmt.Println(err)
				}
				for _, m := range mop.GetExpiredMatter() {
					pf := P0p1Filter{}
					go MarkClosed(pf, m.Id, mop)
					//go mop.GetLatestEventStatus(m.Id)

				}

			case _, ok := <-done:
				if !ok {
					return
				}
			}
		}
	}()
	return nil
}

func getLastEventId() (id int64) {
	//err := db_alarm.Raw("SELECT MAX(event) FROM matter_event ").QueryRow(&id)
	//if err != nil || id == 0 {
	db_alarm.Raw("select AUTO_INCREMENT from information_schema.tables where table_name = 'event'").QueryRow(&id)
	glog.V(4).Infof("get last event id %d", id-1)
	//}
	return id - 1
}

func syncEvent(id int64, limit int) int64 {
	var events []*Event

	//"event_id, status, endpoint, hostname, is_customize, pdl, idc, current_step, event_ts, metric, pushed_tags, left_val, func, op, right_val, max_step, priority, note, tpl_id, action_id, expression_id, strategy_id, last_values"
	_, err := db_alarm.Raw("select id,pdl,action_id,db_ts,endpoint from event where id > ? order by id limit ?", id, limit).QueryRows(&events)
	if err != nil {
		return id
	}

	for _, e := range events {
		_, err := CreateMatter(e)
		if err != nil {
			glog.Errorf(MODULE_NAME+" creatematter err %s", err.Error())
		}
	}
	return int64(len(events))
}
