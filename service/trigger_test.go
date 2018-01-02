/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

var (
	test_db_init bool
	test_db      orm.Ormer
)

type testItem struct {
	endpoint string
	metric   string
	tags     string
	typ      falcon.ItemType
}

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

	test_db, err = falcon.NewOrm("test_service_sync", dsn, 10, 10)
	if err != nil {
		return
	}
	test_db_init = true
}

func TestTirggerDb(t *testing.T) {
	var (
		err           error
		shard         *ShardModule
		trigger       *Trigger
		treeNodes     []*TreeNode
		tagHosts      []*TagHost
		eventTriggers []*EventTrigger
	)

	if !test_db_init {
		t.Logf("test db not inited, skip test sync\n")
		return
	}

	shard = &ShardModule{
		bucketMap: make(map[int32]*bucketEntry),
	}

	trigger = &Trigger{}

	if treeNodes, err = getTreeNodes(test_db); err != nil {
		t.Error(err)
	}
	if err = setTreeNodes(treeNodes, trigger); err != nil {
		t.Error(err)
	}

	if tagHosts, err = getTagHosts(test_db); err != nil {
		t.Error(err)
	}
	if err = setTagHosts(tagHosts, trigger); err != nil {
		t.Error(err)
	}

	if eventTriggers, err = getEventTriggers(test_db); err != nil {
		t.Error(err)
	}
	if err = setEventTriggers(eventTriggers, trigger); err != nil {
		t.Error(err)
	}

	if err = setServiceShards(shard, trigger); err != nil {
		t.Error(err)
	}
}

func testGenerateShard(items []*testItem) *ShardModule {

	shard := &ShardModule{
		bucketMap: make(map[int32]*bucketEntry),
	}
	bucket := &bucketEntry{itemMap: make(map[string]*itemEntry)}

	for _, v := range items {
		item := &itemEntry{
			endpoint: []byte(v.endpoint),
			metric:   []byte(v.metric),
			tags:     []byte(v.tags),
			typ:      v.typ,
		}
		bucket.itemMap[item.key()] = item
	}
	shard.bucketMap[0] = bucket

	return shard
}

func TestTrigger(t *testing.T) {
	treeNodes := []*TreeNode{
		{2, "cap=xiaomi", 1, nil, nil, nil},
		{3, "cap=xiaomi,owt=inf", 2, nil, nil, nil},
		{4, "cap=xiaomi,owt=miliao", 2, nil, nil, nil},
		{5, "cap=xiaomi,owt=miliao,pdl=op", 4, nil, nil, nil},
		{6, "cap=xiaomi,owt=miliao,pdl=micloud", 4, nil, nil, nil},
	}
	tagHosts := []*TagHost{
		{2, "xiaomi.bj"},
		{3, "inf.xiaomi.bj"},
		{4, "miliao.xiaomi.bj"},
		{5, "op.miliao.xiaomi.bj"},
		{6, "micloud.miliao.xiaomi.bj"},
		{5, "op_micloud.miliao.xiaomi.bj"},
		{6, "op_micloud.miliao.xiaomi.bj"},
	}
	eventTriggers := []*EventTrigger{
		/*id, parent_id, tag_id, priority, name, metric, tags, func, op, vlue, msg */
		{1, 0, 1, 1, "cpu", "cpu.busy", "", "all(#3)", ">", "90", "", "cpu busy over 90%", nil, nil},
		{2, 1, 1, 2, "", "cpu.busy", "", "all(#3)", ">", "80", "", "cpu busy over 80%", nil, nil},
		{3, 0, 5, 1, "cpu", "cpu.busy", "", "all(#3)", ">", "60", "", "cpu busy over 60%", nil, nil},
		{4, 3, 5, 2, "", "cpu.busy", "", "all(#3)", ">", "50", "", "cpu busy over 50%", nil, nil},
		{5, 0, 6, 1, "cpu", "cpu.busy", "", "all(#3)", ">", "99", "", "cpu busy over 99%", nil, nil},
	}

	testItems := []*testItem{
		{"xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"inf.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"op.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"micloud.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"op_micloud.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
	}

	trigger := &Trigger{}

	shard := testGenerateShard(testItems)

	if err = setTreeNodes(treeNodes, trigger); err != nil {
		t.Error(err)
	}
	if err = setTagHosts(tagHosts, trigger); err != nil {
		t.Error(err)
	}
	if err = setEventTriggers(eventTriggers, trigger); err != nil {
		t.Error(err)
	}
	if err = setServiceShards(shard, trigger); err != nil {
		t.Error(err)
	}

	t.Logf("=== host tag\n")
	for k, v := range trigger.HostTnodes {
		t.Logf("%s\n", k)
		for _, v1 := range v {
			t.Logf("    %s\n", v1.Name)
		}
	}

	t.Logf("=== trigger item\n")

	cnt := [2]int{0, len(tagHosts)*2 - 2}
	for _, node := range trigger.TnodeIds {
		if len(node.ETriggerMetric) == 0 {
			continue
		}
		t.Logf("tag[%s]\n", node.Name)
		for _, triggers := range node.ETriggerMetric {
			for _, trigger := range triggers {
				t.Logf("    %s %s%s%s msg '%s' tags '%s'\n",
					trigger.Metric, trigger.Func, trigger.Op, trigger.Value,
					trigger.Msg, trigger.Tags)
				for _, item := range trigger.items {
					cnt[0]++
					t.Logf("        %s %s %s\n", item.endpoint, item.metric, item.tags)
				}
			}
		}
	}
	if cnt[0] != cnt[1] {
		t.Errorf("item trigger number got %d want %d\n", cnt[0], cnt[1])
	}

}

func TestTriggerExprEval(t *testing.T) {
}
