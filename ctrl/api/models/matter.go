/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import "fmt"

const (
	STATUS_CREATING = iota
	STATUS_PENDING
	STATUS_PROCESSING
	STATUS_SOLVED
	STATUS_UPGRADED
	STATUS_IGNORED
)

type Matter struct {
	Id          int64
	Status      int
	Aggregation string // Aggregation
	//Tag         string
	Starttime int64
	Endtime   int64
	Name      string
}

type Events struct {
	Eventid    string
	Endpoint   string
	Metric     string
	Func       string
	Priority   string
	Status     string
	Timestamp  string
}

type Claim struct {
	Matter    int64
	User      string
	Timestamp int64
	Commit    string
}

func (op *Operator) QueryMatters(status int, per int, offset int) ([]Matter, error) {
	var matters []Matter
	_, err := Db.Alarm.Raw("SELECT a.matter AS id ,  b.status, b.aggregation, b.starttime , b.endtime, b.name FROM alarm_event.user_matter a ,alarm_event.matter b WHERE a.user=? AND b.status=? AND b.id=a.matter order by b.id desc limit ? offset  ?", op.User.Name, status, per, offset).QueryRows(&matters)
	fmt.Println(matters)
	fmt.Println(op.User.Name)
	return matters, err
}

func (op *Operator) QueryEventsByMatter(matterID int64, per, offset int) []Events {
	var events []Events
	_, err := Db.Alarm.Raw("SELECT distinct a.event AS eventid, b.endpoint, b.metric,b.func,b.priority, b.status,b.db_ts as timestamp FROM alarm_event.matter_event a, alarm_event.event b WHERE a.matter=? AND a.event =b.id order by a.id desc limit ? offset ?", matterID, per, offset).QueryRows(&events)
	if err != nil {
		fmt.Println(err)
	}
	encountered := map[Events]bool{}
	result := []Events{}
	for v := range events {
		if encountered[events[v]] == true {
		} else {
			encountered[events[v]] = true
			result = append(result, events[v])
		}
	}
	return result
}

func (op *Operator) QueryEventsCntByMatter(matterID int64) (int, error) {
	var events []Events
	_, err := Db.Alarm.Raw("SELECT distinct a.event AS eventid FROM alarm_event.matter_event a, alarm_event.event b WHERE a.matter=? AND a.event =b.id order by a.id desc", matterID).QueryRows(&events)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	encountered := map[Events]bool{}
	result := []Events{}
	for v := range events {
		if encountered[events[v]] == true {
		} else {
			encountered[events[v]] = true
			result = append(result, events[v])
		}
	}
	cnt := len(result)
	return cnt, err
}

func (op *Operator) UpdateMatter(id int64, _m Matter) error {
	if _, err := op.GetMatter(id); err != nil {
		return err
	}
	_, err := Db.Alarm.Raw("UPDATE alarm_event.matter SET status = ?,name= ? WHERE id =?", _m.Status, _m.Name, id).Exec()
	return err

}
func (op *Operator) GetMatter(id int64) (Matter, error) {
	var matter Matter
	err := Db.Alarm.Raw("SELECT id , status, tag, uic , starttime , endtime, name FROM alarm_event.matter WHERE id=? ", id).QueryRow(&matter)
	return matter, err
}

func (op *Operator) GetMatterCnt(status int) (int64, error) {
	var matterCnt int64
	err := Db.Alarm.Raw("SELECT count(*) AS cnt FROM alarm_event.user_matter a ,alarm_event.matter b WHERE a.user=? AND b.status=? AND b.id=a.matter order by b.id desc", op.User.Name, status).QueryRow(&matterCnt)
	return matterCnt, err
}

func (op *Operator) AddClaim(claim Claim) error {
	op.O.Using("alarm")
	op.O.Begin()
	_, err := op.O.Raw("UPDATE alarm_event.matter SET status = ? WHERE id =?", STATUS_PROCESSING, claim.Matter).Exec()
	_, err = op.O.Raw("INSERT INTO alarm_event.matter_claim (`matter`,`user`,`comment`,`timestamp`) VALUES (?,?,?,?)", claim.Matter, claim.User, claim.Commit, claim.Timestamp).Exec()
	if err != nil {
		err = op.O.Rollback()
	} else {
		err = op.O.Commit()
	}
	op.O.Using("default")
	return err
}
