/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

// add "PluginAlarm true;"config file at ctrl section to enable this plugin
import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon/alarm/api"
	"time"
)

const (
	STATUS_CREATING   = iota
	STATUS_PENDING
	STATUS_PROCESSING
	STATUS_SOLVED
	STATUS_UPGRADED
	STATUS_IGNORED
)

//HandlerMatter 处理Matter
func CreateMatter(event *Event) (matterIds []int64, err error) {

	op := &MatterOperator{}
	op.DB, err = orm.GetDB("alarm")
	if err != nil {
		return
	}

	matters := op.GenMatter(event)
	unSolvedMatters := op.GetUnsolvedMatter(event)
	for _, matter := range matters {
		//op.Matter = matter
		if matterRecord, dup := checkMatterDuplicates(unSolvedMatters, matter); dup {
			matterRecord.EndTime = event.DbTs.Unix() + matterRecord.TimeWindow
			op.RefreshMatter(matterRecord)
			op.RecordEvent(matterRecord, event)
			matterIds = append(matterIds, matterRecord.Id)
		} else {
			matterID := op.WriteMatter(matter)
			matter.Id = matterID
			matterIds = append(matterIds, matterID)
			err := op.RecordUserMatter(matter)
			if err != nil {
				fmt.Println(err)
			}
			op.RecordEvent(matter, event)
		}

	}
	return matterIds, nil
}

//checkMatterDuplicates 检查matter是否存在
func checkMatterDuplicates(matters []Matter, matter Matter) (Matter, bool) {
	for _, m := range matters {
		if matter.Aggregation == m.Aggregation {
			matter.Id = m.Id
			return matter, true
		} else if matter.Uic == m.Uic {
			matter.Id = m.Id
			return matter, true
		}
	}
	return matter, false
}

type Matter struct {
	Id          int64       // id
	Status      int         // Current status
	Aggregation string      // Aggregation
	Aggregate   string      // Aggregate
	Uic         string      // Uic
	Users       []*api.User // Users to notic
	TimeWindow  int64       // Close matter TimeWindow
	StartTime   int64       // matter start ts
	EndTime     int64       // matter end time
	Name        string      //matter name
}

type MatterOperator struct {
	DB *sql.DB
}

func (op *MatterOperator) GenMatter(event *Event) []Matter {

	var matters []Matter
	var users []*api.User
	var uic string

	for _, tag := range strings.Split(event.Pdl, ";") {
		err := db_ctrl.Raw("select uic from action where id = ?", event.ActionId).QueryRow(&uic)
		if err != nil {
			glog.Error("action %s no found", event.ActionId)
			continue
		}

		db_ctrl.Raw(fmt.Sprintf("select name, email, phone from user where id in (SELECT user_id FROM team_user tu left join team t on t.id = tu.team_id where t.name in ('%s'))", strings.Replace(uic, ",", "','", -1))).QueryRows(&users)

		ts := event.DbTs.Unix()

		matter := Matter{
			Status:      STATUS_PENDING,
			Aggregation: tag,
			Uic:         uic,
			TimeWindow:  1200,
			StartTime:   ts,
			EndTime:     ts + 1200,
			Users:       users,
			Name:        event.Endpoint,
		}

		matters = append(matters, matter)
	}
	return matters

}

// WriteMatter 写入matter
func (op *MatterOperator) WriteMatter(matter Matter) int64 {
	log.Printf("creating matter")
	sql, _ := op.DB.Prepare("INSERT alarm_event.matter set status=?,aggregation=?,uic=?,starttime=?,endtime=?,name=?")
	res, err := sql.Exec(STATUS_PENDING, matter.Aggregation, matter.Uic, matter.StartTime, matter.EndTime, matter.Name)
	if err != nil {
		fmt.Println(err)
	}
	matterID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	return matterID
}

//RefreshMatter 刷新matter结束时间1
func (op *MatterOperator) RefreshMatter(matter Matter) (err error) {
	log.Printf("refreshing  matter")
	sql, _ := op.DB.Prepare("UPDATE alarm_event.matter SET endtime=? WHERE id=?;")
	_, err = sql.Exec(matter.EndTime, matter.Id)
	return
}

//RecordEvent 记录event与matter关系
func (op *MatterOperator) RecordEvent(matter Matter, event *Event) (err error) {
	sql, _ := op.DB.Prepare("INSERT INTO alarm_event.matter_event SET matter=?,event=? ON DUPLICATE KEY UPDATE event=?")
	_, err = sql.Exec(matter.Id, event.Id, event.Id)
	return
}

//GetMatter 获取matter
func (op *MatterOperator) GetMatter(id int64) Matter {
	sql := "select uic,aggregation,status,starttime,endtime from alarm_event.matter where id='" + strconv.FormatInt(id, 10) + "'"
	rows, _ := op.DB.Query(sql)
	defer rows.Close()
	var uic, aggregation string
	var status int
	var starttime, endtime int64
	for rows.Next() {
		rows.Scan(&uic, &aggregation, &status, &starttime, &endtime)
	}
	matter := Matter{
		Id:          id,
		Uic:         uic,
		Aggregation: aggregation,
		Status:      status,
		StartTime:   starttime,
		EndTime:     endtime,
	}
	return matter

}

//GetUnslovedMatter 获取未恢复matter
func (op *MatterOperator) GetUnsolvedMatter(event *Event) []Matter {
	sql := "select id,uic,aggregation,status,starttime,endtime from alarm_event.matter where endtime >" + strconv.FormatInt(event.DbTs.Unix(), 10)
	rows, err := op.DB.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	//defer rows.Close()
	var matter Matter
	var matters []Matter
	for rows.Next() {
		rows.Scan(&matter.Id, &matter.Uic, &matter.Aggregation, &matter.Status, &matter.StartTime, &matter.EndTime)
		matters = append(matters, matter)

	}
	return matters
}

//GetUnslovedMatter 获取未恢复matter
func (op *MatterOperator) GetExpiredMatter() []Matter {
	sql := "select id,uic,aggregation,status,starttime,endtime from alarm_event.matter where status <" + strconv.Itoa(STATUS_SOLVED) + " and endtime <" + strconv.FormatInt(time.Now().Unix(), 10)
	rows, err := op.DB.Query(sql)
	if err != nil {
		fmt.Println(err)
	}
	//defer rows.Close()
	var matter Matter
	var matters []Matter
	for rows.Next() {
		rows.Scan(&matter.Id, &matter.Uic, &matter.Aggregation, &matter.Status, &matter.StartTime, &matter.EndTime)
		matters = append(matters, matter)

	}
	return matters
}

type EventState struct {
	id       int64
	priority int
	eventid  string
	status   string
}

func (op *MatterOperator) GetLatestEventStatus(id int64) (error, map[string]EventState) {
	sql := "SELECT b.id, b.priority, b.event_id ,b.status FROM alarm_event.matter_event a,alarm_event.event b WHERE a.event=b.id AND a.matter = " + strconv.FormatInt(id, 10) + " order by b.id desc"
	rows, err := op.DB.Query(sql)
	if err != nil {
		return err, nil
	}
	esmap := make(map[string]EventState)
	var estates []EventState
	var estate EventState
	for rows.Next() {
		rows.Scan(&estate.id, &estate.priority, &estate.eventid, &estate.status)
		esmap[estate.eventid] = estate
		estates = append(estates, estate)
	}
	defer rows.Close()
	return nil, esmap
}

type priorityFilter interface {
	pfilter(states LastEventStates) bool
}

func MarkClosed(f priorityFilter, id int64, mop *MatterOperator) error {
	_, status := mop.GetLatestEventStatus(id)
	if ok := f.pfilter(status); ok {
		mop.CloseMatter(id)
	}

	return nil
}

type P0p1Filter struct{}

func (f P0p1Filter) pfilter(states LastEventStates) (ok bool) {
	var PROBLEM = "PROBLEM"
	var statarr []string
	for k, v := range states {
		if v.priority < 5 && v.priority >= 0 {
			fmt.Println(k, v)
			statarr = append(statarr, v.status)
		}
	}
	for _, e := range statarr {
		if e == PROBLEM {
			return false
		}
	}
	return true
}

type LastEventStates map[string]EventState

// RecordUserMatter user与matter关系
func (op *MatterOperator) RecordUserMatter(matter Matter) (err error) {
	for _, user := range matter.Users {
		fmt.Printf("%+v\n", user)
		sql, _ := op.DB.Prepare("INSERT alarm_event.user_matter set user=?,matter=?")
		_, err := sql.Exec(user.Name, matter.Id)
		return err
	}
	return
}

// 关闭matter
func (op *MatterOperator) CloseMatter(id int64) (err error) {
	sql, _ := op.DB.Prepare("UPDATE alarm_event.matter SET status = 3  WHERE id =?")
	_, err = sql.Exec(id)
	return err
}
