/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl"
)

type LogApi struct {
	LogId  int64
	Module int64
	Id     int64
	User   string
	Action int64
	Data   string
	Time   time.Time
}

type LogApiGet struct {
	LogId  int64     `json:"id"`
	Module string    `json:"module"`
	Id     int64     `json:"tid"`
	User   string    `json:"user"`
	Action string    `json:"action"`
	Data   string    `json:"data"`
	Time   time.Time `json:"time"`
}

func logSql(begin, end string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if begin != "" {
		sql2 = append(sql2, "a.time >= ?")
		sql3 = append(sql3, begin)
	}
	if end != "" {
		sql2 = append(sql2, "a.time <= ?")
		sql3 = append(sql3, end)
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}

func (op *Operator) GetLogsCnt(begin, end string) (cnt int64, err error) {
	sql, sql_args := logSql(begin, end)
	err = op.O.Raw("SELECT count(*) FROM log a "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetLogs(begin, end string, limit, offset int) (ret []*LogApi, err error) {
	sql, sql_args := logSql(begin, end)
	sql = "select a.id as log_id, a.module, a.module_id as id, b.name as user, a.action, a.data, a.time from log a left join user b on a.user_id = b.id " + sql + " ORDER BY a.id DESC LIMIT ? OFFSET ?"
	sql_args = append(sql_args, limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

/* for demo */
func (op *Operator) populate() (interface{}, error) {
	var (
		err       error
		items     []string
		id        int64
		tag_idx   = make(map[string]int64)
		user_idx  = make(map[string]int64)
		role_idx  = make(map[string]int64)
		token_idx = make(map[string]int64)
		host_idx  = make(map[string]int64)
		tpl_idx   = make(map[string]int64)
		test_user = "test01"
	)
	tag_idx["/"] = 1

	// user
	items = []string{
		test_user,
		"user0",
		"user1",
		"user2",
		"user3",
		"user4",
		"user5",
		"user6",
	}
	for _, item := range items {
		if id, err = op.CreateUser(&UserApiAdd{Name: item, Uuid: item}); err != nil {
			return nil, err
		}
		user_idx[item] = id
		glog.Infof("add user(%s)\n", item)
	}

	// tag
	items = []string{
		"cop=xiaomi",
		"cop=xiaomi,owt=inf",
		"cop=xiaomi,owt=miliao",
		"cop=xiaomi,owt=miliao,pdl=op",
		"cop=xiaomi,owt=miliao,pdl=micloud",
	}
	for _, item := range items {
		glog.Infof("add tag(%s)\n", item)
		if tag_idx[item], err = op.CreateTag(&TagCreate{Name: item}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// tag host
	items2 := [][2]string{
		{"cop=xiaomi", "mi1.bj"},
		{"cop=xiaomi", "mi2.bj"},
		{"cop=xiaomi", "mi3.bj"},
		{"cop=xiaomi,owt=inf", "inf1.bj"},
		{"cop=xiaomi,owt=inf", "inf2.bj"},
		{"cop=xiaomi,owt=inf", "inf3.bj"},
		{"cop=xiaomi,owt=miliao", "miliao1.bj"},
		{"cop=xiaomi,owt=miliao", "miliao2.bj"},
		{"cop=xiaomi,owt=miliao", "miliao3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=op", "miliao.op3.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud1.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud2.bj"},
		{"cop=xiaomi,owt=miliao,pdl=micloud", "miliao.cloud3.bj"},
	}
	for _, item2 := range items2 {
		glog.Infof("add host(%s, %s)\n", item2[1], item2[0])
		if host_idx[item2[1]], err = op.CreateHost(&HostCreate{Name: item2[1]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}

		if _, err = op.CreateTagHost(&RelTagHostApiAdd{TagId: tag_idx[item2[0]],
			HostId: host_idx[item2[1]]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		glog.Infof("add tag(%s)\n", item)
		if id, err = op.AddAction(&Action{}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
		if tpl_idx[item], err = op.AddTemplate(&Template{Name: item,
			ActionId: id}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}
	// template strategy
	items2 = [][2]string{
		{"tpl1", "cpu.busy"},
		{"tpl1", "cpu.cnt"},
		{"tpl1", "cpu.idle"},
		{"tpl2", "cpu.busy"},
		{"tpl2", "cpu.cnt"},
		{"tpl2", "cpu.idle"},
		{"tpl3", "cpu.busy"},
		{"tpl3", "cpu.cnt"},
		{"tpl3", "cpu.idle"},
	}
	for _, item2 := range items2 {
		glog.Infof("add strategy(%s, %s)\n", item2[0], item2[1])
		if _, err = op.AddStrategy(&Strategy{Metric: item2[1],
			TplId: tpl_idx[item2[0]]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// clone template
	items = []string{
		"tpl1",
		"tpl2",
		"tpl3",
	}
	for _, item := range items {
		glog.Infof("clone template(%s)\n", item)
		if _, err = op.CloneTemplate(tpl_idx[item]); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// bind tag template
	items2 = [][2]string{
		{"cop=xiaomi", "tpl1"},
		{"cop=xiaomi", "tpl2"},
		{"cop=xiaomi", "tpl3"},
		{"cop=xiaomi,owt=inf", "tpl1"},
		{"cop=xiaomi,owt=inf", "tpl2"},
		{"cop=xiaomi,owt=inf", "tpl3"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl1"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl2"},
		{"cop=xiaomi,owt=miliao,pdl=op", "tpl3"},
	}
	for _, item2 := range items2 {
		glog.Infof("add tag tpl(%s, %s)\n", item2[0], item2[1])
		if _, err = op.CreateTagTpl(&RelTagTpl{TagId: tag_idx[item2[0]],
			TplId: tpl_idx[item2[1]]}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// role
	items = []string{
		"adm",
		"sre",
		"dev",
		"usr",
	}
	for _, item := range items {
		glog.Infof("add role(%s)\n", item)
		if role_idx[item], err = op.CreateRole(&RoleCreate{Name: item}); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// token
	for i := SYS_IDX_R_TOKEN; i < SYS_IDX_TOKEN_SIZE; i++ {
		token_idx[tokenName[i]] = int64(i)
	}
	/*
		items = []string{
			SYS_R_TOKEN,
			SYS_O_TOKEN,
			SYS_A_TOKEN,
		}
		for _, item := range items {
			if token_idx[item], err = op.CreateToken(&TokenCreate{Name: item}); err != nil {
				return nil, err
			}
			glog.Infof("add token(%s)\n", item)
		}
	*/

	// bind user
	binds := [][3]string{
		{"cop=xiaomi,owt=miliao", test_user, "adm"},
		{"cop=xiaomi,owt=miliao", test_user, "sre"},
		{"cop=xiaomi,owt=miliao", test_user, "dev"},
		{"cop=xiaomi,owt=miliao", test_user, "usr"},
	}
	for _, s := range binds {
		glog.Infof("bind tag(%s) user(%s) role(%s)\n", s[0], s[1], s[2])
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[0]], role_idx[s[2]],
			user_idx[s[1]], TPL_REL_T_ACL_USER); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	// bind token
	binds = [][3]string{
		{SYS_O_TOKEN, "adm", "/"},
		{SYS_R_TOKEN, "adm", "/"},
		{SYS_A_TOKEN, "adm", "/"},
		{SYS_O_TOKEN, "sre", "/"},
		{SYS_R_TOKEN, "sre", "/"},
		{SYS_R_TOKEN, "dev", "/"},
		{SYS_R_TOKEN, "usr", "/"},
		{SYS_O_TOKEN, "adm", "cop=xiaomi,owt=miliao"},
		{SYS_O_TOKEN, "dev", "cop=xiaomi,owt=miliao,pdl=op"},
		{SYS_O_TOKEN, "usr", "cop=xiaomi"},
		{SYS_A_TOKEN, "usr", "cop=xiaomi,owt=miliao"},
	}
	for _, s := range binds {
		glog.Infof("bind tag(%s) token(%s) role(%s)\n", s[2], s[0], s[1])
		if _, err := addTplRel(op.O, op.User.Id, tag_idx[s[2]], role_idx[s[1]],
			token_idx[s[0]], TPL_REL_T_ACL_TOKEN); err != nil {
			glog.Error(err.Error())
			return nil, err
		}
	}

	return "populate db done", nil
}

func (op *Operator) ResetDb(populate bool) (interface{}, error) {
	var err error
	var cmds []string

	file := ctrl.Configure.Ctrl.Str(ctrl.C_DB_SCHEMA)
	if file == "" {
		return "", fmt.Errorf("please config ctrl:dbSchema  file path")
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(buf), "\n")
	for cmd, in, i := "", false, 0; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			continue
		}

		if in {
			cmd += " " + strings.TrimSpace(line)
			if line[len(line)-1] == ';' {
				cmds = append(cmds, cmd)
				in = false
			}
		} else {
			n := strings.Index(line, " ")
			if n <= 0 {
				continue
			}

			switch line[:n] {
			case "SET", "CREATE", "INSERT", "DROP":
				cmd = line
				if line[len(line)-1] == ';' {
					cmds = append(cmds, cmd)
				} else {
					in = true
				}
			}
		}
	}

	for i := 0; i < len(cmds); i++ {
		_, err := op.O.Raw(cmds[i]).Exec()
		if err != nil {
			glog.Error(MODULE_NAME+" sql %s ret %s", cmds[i], err.Error())
		}
	}

	/*
		for _, table := range dbTables {
			if _, err = op.O.Raw("TRUNCATE TABLE `" + table + "`").Exec(); err != nil {
				return nil, err
			}
		}

		// init root tree tag
		op.SqlExec("insert tag (id, name) values (1, '')")
		op.SqlExec("insert tag_rel (tag_id, sup_tag_id) values (1, 1)")
		op.SqlExec("insert user (id, uuid, name, cname, email) values (1, 'root@localhost', 'system', 'system', 'root@localhost')")
		op.SqlExec("insert token (id, name, cname, note) values " +
			"(1, 'falcon_read', 'read', 'read') " +
			"(2, 'falcon_operate', 'operate', 'operate') " +
			"(3, 'falcon_admin', 'admin', 'admin')")

		op.SqlExec("alter table host auto_increment=1000")
		op.SqlExec("alter table token auto_increment=1000")
		op.SqlExec("alter table role auto_increment=1000")
		op.SqlExec("alter table tag auto_increment=1000")
		op.SqlExec("alter table user auto_increment=1000")
	*/

	// reset cache
	// ugly hack
	initCache(ctrl.Configure)

	if populate {
		return op.populate()
	}

	return "reset db done", nil
}
