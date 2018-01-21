/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

var (
	test_db_init bool
	test_db      orm.Ormer
)

func init() {
	var err error

	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	user := env("MYSQL_TEST_USER", "falcon")
	pass := env("MYSQL_TEST_PASS", "1234")
	prot := env("MYSQL_TEST_PROT", "tcp")
	addr := env("MYSQL_TEST_ADDR", "localhost:3306")
	dbname := env("MYSQL_TEST_DBNAME", "falcon")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)

	test_db, _, err = falcon.NewOrm("test_action_sync", dsn, 10, 10)
	if err != nil {
		return
	}
	test_db_init = true
}

func testTirggerDb(t *testing.T) {

	if !test_db_init {
		t.Logf("test db not inited, skip test sync\n")
		return
	}

	trigger := &Trigger{db: test_db}
	if err := trigger.updateNode(); err != nil {
		t.Error(err)
	}
	if err := trigger.updateUser(); err != nil {
		t.Error(err)
	}
	if err := trigger.updateTagRoleUser(); err != nil {
		t.Error(err)
	}
	if err := trigger.updateTagRoleToken(); err != nil {
		t.Error(err)
	}

	if err := trigger.updateTagTokenUser(); err != nil {
		t.Error(err)
	}

	if err := trigger.updateActionTrigger(); err != nil {
		t.Error(err)
	}
}

const (
	USER1   = 1 // user
	USER2   = 2
	USER3   = 3
	ADMIN   = 1 // role
	DEV     = 2
	SRE     = 3
	WRITE   = 1 // token
	READ    = 2
	ROOT    = 1 // node id
	XIAOMI  = 2
	INF     = 3
	MILIAO  = 4
	OP      = 5
	MICLOUD = 6
)

func TestTrigger(t *testing.T) {
	nodes := []*Node{
		{sync.RWMutex{}, "cap=xiaomi", XIAOMI, ROOT, nil, nil, nil, nil, nil},
		{sync.RWMutex{}, "cap=xiaomi,owt=inf", INF, XIAOMI, nil, nil, nil, nil, nil},
		{sync.RWMutex{}, "cap=xiaomi,owt=miliao", MILIAO, XIAOMI, nil, nil, nil, nil, nil},
		{sync.RWMutex{}, "cap=xiaomi,owt=miliao,pdl=op", OP, MILIAO, nil, nil, nil, nil, nil},
		{sync.RWMutex{}, "cap=xiaomi,owt=miliao,pdl=micloud", MICLOUD, MILIAO, nil, nil, nil, nil, nil},
	}
	users := []*User{
		{USER1, "user1", "user1", "user1@example.com", "111111111"},
		{USER2, "user2", "user2", "user2@example.com", "222222222"},
		{USER3, "user3", "user3", "user3@example.com", "333333333"},
	}
	tagRoleUsers := []*TagRoleUser{
		{XIAOMI, ADMIN, USER1},
		{XIAOMI, SRE, USER2},
		{OP, ADMIN, USER2},
		{MILIAO, DEV, USER3},
	}
	tagRoleTokens := []*TagRoleToken{
		{XIAOMI, ADMIN, WRITE},
		{XIAOMI, ADMIN, READ},
		{XIAOMI, DEV, READ},
		{XIAOMI, SRE, READ},
		{OP, DEV, WRITE},
		{OP, DEV, READ},
	}

	actionTriggers := []*ActionTrigger{
		// 优先级为 0,1　的事件，给具有　WRITE 权限的人发送邮件和短信
		{1, XIAOMI, WRITE, 100, "value(priority)<2", ACTION_F_EMAIL | ACTION_F_SMS, "", nil, nil},
		{2, MILIAO, WRITE, 100, "value(priority)=0", ACTION_F_EMAIL | ACTION_F_SMS | ACTION_F_SCRIPT, "curl http://alarm.example.com/?key=${key}&msg=${msg}", nil, nil},
		{3, MILIAO, WRITE, 100, "value(priority)=1", ACTION_F_EMAIL | ACTION_F_SMS, "", nil, nil},
		{2, MILIAO, WRITE, 100, "value(priority)<4", ACTION_F_EMAIL, "", nil, nil},
	}

	trigger := &Trigger{}

	if setNodes(nodes, trigger) != nil ||
		setUsers(users, trigger) != nil ||
		setTagRoleUser(tagRoleUsers, trigger) != nil ||
		setTagRoleToken(tagRoleTokens, trigger) != nil ||
		trigger.updateTagTokenUser() != nil ||
		setActionTriggers(actionTriggers, trigger) != nil {
		t.Error("env init failed")
	}

	glog.V(4).Infof("=== tag token user\n")

	cases := []struct {
		tag   int64
		token int64
		user  int64
		want  bool
	}{
		{XIAOMI, READ, USER1, true},
		{XIAOMI, READ, USER2, true},
		{XIAOMI, READ, USER3, false},
		{XIAOMI, WRITE, USER1, true},
		{XIAOMI, WRITE, USER2, false},
		{XIAOMI, WRITE, USER3, false},

		{INF, READ, USER1, true},
		{INF, READ, USER2, true},
		{INF, READ, USER3, false},
		{INF, WRITE, USER1, true},
		{INF, WRITE, USER2, false},
		{INF, WRITE, USER3, false},

		{MILIAO, READ, USER1, true},
		{MILIAO, READ, USER2, true},
		{MILIAO, READ, USER3, true},
		{MILIAO, WRITE, USER1, true},
		{MILIAO, WRITE, USER2, false},
		{MILIAO, WRITE, USER3, false},

		{OP, READ, USER1, true},
		{OP, READ, USER2, true},
		{OP, READ, USER3, true},
		{OP, WRITE, USER1, true},
		{OP, WRITE, USER2, true},
		{OP, WRITE, USER3, true},

		{MICLOUD, READ, USER1, true},
		{MICLOUD, READ, USER2, true},
		{MICLOUD, READ, USER3, true},
		{MICLOUD, WRITE, USER1, true},
		{MICLOUD, WRITE, USER2, false},
		{MICLOUD, WRITE, USER3, false},
	}

	glog.V(5).Infof("%s", trigger.nodes[1])

	for _, tc := range cases {
		if got := trigger.nodes[tc.tag].tokenUser[tc.token][tc.user]; got != tc.want {
			t.Errorf("tag %d token %d user %d = %v; want %v", tc.tag, tc.token, tc.user, got, tc.want)
		}
	}

	eventCases := []struct {
		tag       int64
		key       string
		expr      string
		msg       string
		timestamp int64
		value     float64
		priority  int
		want      int
	}{
		{XIAOMI, "xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 0, 1},
		{INF, "inf.xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 0, 1},
		{MILIAO, "miliao.xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 0, 1},
		{OP, "op.miliao.xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 0, 3},
		{MICLOUD, "micloud.miliao.xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 0, 1},
		{MICLOUD, "micloud.miliao.xiaomi.bj/cpu.busy//GUAGE", "count(#3,0,>)=3", "cpu busy over 0%", 0, 0, 3, 1},
	}

	for _, tc := range eventCases {
		e := &eventEntry{
			lastTs:    time.Now().Unix(),
			tagId:     tc.tag,
			key:       tc.key,
			expr:      tc.expr,
			msg:       tc.msg,
			timestamp: tc.timestamp,
			value:     tc.value,
			priority:  tc.priority,
		}
		node := trigger.nodes[tc.tag]
		if got := actionGenerate(node.processEvent(e), node, e, trigger.users); len(got.users) != tc.want {
			t.Errorf("actionGenerate tag %d key %s user num got %d; want %d", tc.tag, tc.key, len(got.users), tc.want)
		}
	}

}
