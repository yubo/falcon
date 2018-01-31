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
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm"
	"github.com/yubo/falcon/lib/tsdb"
	"golang.org/x/net/context"
)

var (
	test_db_init bool
	test_db      orm.Ormer
)

type testItem struct {
	key string
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

	test_db, _, err = falcon.NewOrm("test_service_sync", dsn, 10, 10)
	if err != nil {
		return
	}
	test_db_init = true
}

func testTirggerDb(t *testing.T) {
	var (
		err           error
		shard         *CacheModule
		trigger       *Trigger
		treeNodes     []*Node
		tagHosts      []*TagHost
		eventTriggers []*EventTrigger
	)

	if !test_db_init {
		t.Logf("test db not inited, skip test sync\n")
		return
	}

	shard = &CacheModule{buckets: make([]*cacheBucket, falcon.SHARD_NUM)}

	trigger = &Trigger{}

	if treeNodes, err = getNodes(test_db); err != nil {
		t.Error(err)
	}
	if err = setNodes(treeNodes, trigger); err != nil {
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

	if err = setServiceBuckets(shard, trigger); err != nil {
		t.Error(err)
	}
}

func testFillDps(shard *CacheModule, keys []*tsdb.Key, vs []*tsdb.TimeValuePair, t *testing.T) {
	for _, key := range keys {
		for _, v := range vs {
			if _, err := shard.put(&tsdb.DataPoint{Key: key, Value: v}); err != nil {
				t.Error(err)
			}
		}
	}
}

func testTriggerExprProcess(trigger *Trigger) int {
	var cnt int

	ctx, cancel := context.WithCancel(context.Background())
	eventCh := make(chan *alarm.Event, 32)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-eventCh:
				glog.V(4).Infof("%s %s %s %d",
					string(e.Key), e.Expr,
					e.Msg, e.Timestamp)
				cnt++
			}

		}
	}()

	for _, node := range trigger.nodes {
		judgeTagNode(cache.buckets, node, eventCh)
	}

	time.Sleep(time.Millisecond * 100)
	cancel()

	return cnt
}

func testTrigger(t *testing.T) {
	test_cache_init()
	cache.prestart(cacheApp)
	cache.start(cacheApp)
	defer cache.stop(cacheApp)

	treeNodes := []*Node{
		{2, "cop=xiaomi", 1, nil, nil, nil},
		{3, "cop=xiaomi,owt=inf", 2, nil, nil, nil},
		{4, "cop=xiaomi,owt=miliao", 2, nil, nil, nil},
		{5, "cop=xiaomi,owt=miliao,pdl=op", 4, nil, nil, nil},
		{6, "cop=xiaomi,owt=miliao,pdl=micloud", 4, nil, nil, nil},
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
		{1, 0, 1, 8, "cpu", "cpu.busy", "", "count(#3,0,>)=3", "cpu busy over 0%", nil, nil, nil},
		{2, 1, 1, 2, "", "cpu.busy", "", "count(#3,80,>)=3", "cpu busy over 80%", nil, nil, nil},
		{3, 0, 5, 1, "cpu", "cpu.busy", "", "count(#3,60,>)=3", "cpu busy over 60%", nil, nil, nil},
		{4, 3, 5, 2, "", "cpu.busy", "", "count(#3,50,>)=3", "cpu busy over 50%", nil, nil, nil},
		{5, 0, 6, 1, "cpu", "cpu.busy", "", "count(#3,99,>)=3", "cpu busy over 99%", nil, nil, nil},
	}

	testKeys := []*tsdb.Key{
		{[]byte("xiaomi.bj/cpu.busy//GAUGE"), 0},
		{[]byte("inf.xiaomi.bj/cpu.busy//GAUGE"), 0},
		{[]byte("miliao.xiaomi.bj/cpu.busy//GAUGE"), 0},
		{[]byte("op.miliao.xiaomi.bj/cpu.busy//GAUGE"), 0},
		{[]byte("micloud.miliao.xiaomi.bj/cpu.busy//GAUGE"), 0},
		{[]byte("op_micloud.miliao.xiaomi.bj/cpu.busy//GAUGE"), 0},
	}

	testVs := []*tsdb.TimeValuePair{
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9},
	}

	trigger := &Trigger{}

	testFillDps(cache, testKeys, testVs, t)

	if err = setNodes(treeNodes, trigger); err != nil {
		t.Error(err)
	}
	if err = setTagHosts(tagHosts, trigger); err != nil {
		t.Error(err)
	}
	if err = setEventTriggers(eventTriggers, trigger); err != nil {
		t.Error(err)
	}
	if err = setServiceBuckets(cache, trigger); err != nil {
		t.Error(err)
	}

	glog.V(4).Infof("=== host tag\n")
	for k, v := range trigger.hostNodes {
		glog.V(5).Infof("%s\n", k)
		for _, v1 := range v {
			glog.V(5).Infof("    %s\n", v1.Name)
		}
	}

	glog.V(4).Infof("=== trigger item\n")
	cnt := [2]int{0, len(tagHosts)*2 - 2}
	for _, node := range trigger.nodes {
		if len(node.eventTriggerMetrics) == 0 {
			continue
		}
		glog.V(5).Infof("tag[%s]\n", node.Name)
		for _, triggers := range node.eventTriggerMetrics {
			for _, trigger := range triggers {
				glog.V(5).Infof("    %s %s msg '%s' tags '%s'\n",
					trigger.Metric, trigger.Expr,
					trigger.Msg, trigger.Tags)
				for _, item := range trigger.entries {
					cnt[0]++
					glog.V(5).Infof("        %s %s %s\n", item.endpoint, item.metric, item.tags)
				}
			}
		}
	}
	if cnt[0] != cnt[1] {
		t.Errorf("item trigger number got %d want %d\n", cnt[0], cnt[1])
	}

	glog.V(4).Infof("=== trigger item expr exec\n")
	eventNum := testTriggerExprProcess(trigger)
	glog.V(4).Infof("=== trigger item expr exec event : %d\n", eventNum)
}
