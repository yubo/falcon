/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

/*
import (
	"time"

	"github.com/astaxie/beego/orm"
)

type RuleAction struct {
	R Rule   `json:"rule"`
	A Action `json:"action"`
}

type Rule struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name"`
	Pid            int64     `json:"pid"`
	Action_id      int64     `json:"-"`
	Create_user_id int64     `json:"-"`
	Create_time    time.Time `json:"ctime"`
}

type Action struct {
	Id       int64  `json:"id"`
	Sendto   string `json:"sendto"`
	Url      string `json:"url"`
	SendFlag uint32 `json:"send_flag"`
	CbFlag   uint32 `json:"cb_flag"`
}

func (u *User) AddRule(r *Rule) (id int64, err error) {
	id, err = orm.NewOrm().Insert(r)
	if err != nil {
		return
	}
	r.Id = id
	cacheModule[CTL_M_RULE].set(id, r)
	DbLog(u.Id, CTL_M_RULE, id, CTL_A_ADD, jsonStr(r))
	return
}

func (u *User) GetRule(id int64) (*Rule, error) {
	if r, ok := cacheModule[CTL_M_RULE].get(id).(*Rule); ok {
		return r, nil
	}
	r := &Rule{Id: id}
	err := orm.NewOrm().Read(r, "Id")
	if err == nil {
		cacheModule[CTL_M_RULE].set(id, r)
	}
	return r, err
}

func (u *User) QueryRules(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(Rule))
	if query != "" {
		qs = qs.Filter("Name__icontains", query)
	}
	return qs
}

func (u *User) GetRulesCnt(query string) (int64, error) {
	return u.QueryRules(query).Count()
}

func (u *User) GetRules(query string, limit, offset int) (rules []*Rule, err error) {
	_, err = u.QueryRules(query).Limit(limit, offset).All(&rules)
	return
}

func (u *User) UpdateRule(id int64, _r *Rule) (r *Rule, err error) {
	if r, err = u.GetRule(id); err != nil {
		return nil, ErrNoRule
	}

	if _r.Name != "" {
		r.Name = _r.Name
	}
	if _r.Pid != 0 {
		r.Pid = _r.Pid
	}
	_, err = orm.NewOrm().Update(r)
	cacheModule[CTL_M_RULE].set(id, r)
	DbLog(u.Id, CTL_M_RULE, id, CTL_A_SET, "")
	return r, err
}

func (u *User) DeleteRule(id int64) error {
	rule, err := u.GetRule(id)
	if err != nil {
		return err
	}

	if _, err = orm.NewOrm().Delete(&Action{Id: rule.Action_id}); err != nil {
		return err
	}

	if _, err = orm.NewOrm().Delete(&Rule{Id: id}); err != nil {
		return err
	}

	cacheModule[CTL_M_RULE].del(id)
	DbLog(u.Id, CTL_M_RULE, id, CTL_A_DEL, "")

	return nil
}

func (u *User) BindUserRule(user_id, rule_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_rule_user` (`tag_id`, `rule_id`, `user_id`) VALUES (?, ?, ?)", tag_id, rule_id, user_id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *User) BindTokenRule(token_id, rule_id, tag_id int64) (err error) {
	if _, err := orm.NewOrm().Raw("INSERT INTO `tag_rule_token` (`tag_id`, `rule_id`, `token_id`) VALUES (?, ?, ?)", tag_id, rule_id, token_id).Exec(); err != nil {
		return err
	}
	return nil
}
*/
