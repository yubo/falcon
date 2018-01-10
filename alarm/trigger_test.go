/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

/*
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

func testTirggerDb(t *testing.T) {
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

func testFillDps(shard *ShardModule, items []*testItem, dps []*falcon.DataPoint, t *testing.T) {
	for _, v := range items {
		for _, dp := range dps {
			item := &falcon.Item{
				ShardId:   0,
				Endpoint:  []byte(v.endpoint),
				Metric:    []byte(v.metric),
				Tags:      []byte(v.tags),
				Type:      v.typ,
				Value:     dp.Value,
				Timestamp: dp.Timestamp,
			}
			if _, err := shard.put(item); err != nil {
				t.Error(err)
			}
		}
	}
}

func testTriggerExprProcess(trigger *Trigger) int {
	var cnt int

	ctx, cancel := context.WithCancel(context.Background())
	eventCh := make(chan *event, 32)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-eventCh:
				glog.V(4).Infof("%s %s %s %s %d",
					e.trigger.Metric, e.trigger.Tags,
					e.trigger.Expr, e.trigger.Msg,
					e.timestamp)
				cnt++
			}

		}
	}()

	for _, node := range trigger.TnodeIds {
		judgeTagNode(node, eventCh)
	}

	time.Sleep(time.Millisecond * 100)
	cancel()

	return cnt
}

func TestTrigger(t *testing.T) {
	treeNodes := []*TreeNode{
		{2, "cap=xiaomi", 1, nil, nil, nil},
		{3, "cap=xiaomi,owt=inf", 2, nil, nil, nil},
		{4, "cap=xiaomi,owt=miliao", 2, nil, nil, nil},
		{5, "cap=xiaomi,owt=miliao,pdl=op", 4, nil, nil, nil},
		{6, "cap=xiaomi,owt=miliao,pdl=micloud", 4, nil, nil, nil},
	}
	actionTriggers := []*ActionTrigger{
		{1, 0, 1, 8, "cpu", "cpu.busy", "", "count(#3,0,>)=3", "cpu busy over 0%", nil, nil, nil},
		{2, 1, 1, 2, "", "cpu.busy", "", "count(#3,80,>)=3", "cpu busy over 80%", nil, nil, nil},
		{3, 0, 5, 1, "cpu", "cpu.busy", "", "count(#3,60,>)=3", "cpu busy over 60%", nil, nil, nil},
		{4, 3, 5, 2, "", "cpu.busy", "", "count(#3,50,>)=3", "cpu busy over 50%", nil, nil, nil},
		{5, 0, 6, 1, "cpu", "cpu.busy", "", "count(#3,99,>)=3", "cpu busy over 99%", nil, nil, nil},
	}

	testItems := []*testItem{
		{"xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"inf.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"op.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"micloud.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
		{"op_micloud.miliao.xiaomi.bj", "cpu.busy", "", falcon.ItemType_GAUGE},
	}
	testEvent := []*Event{
		{2, "xiaomi.bj/cpu.busy//GUAGE", "kh"},
		{3, "inf.xiaomi.bj/cpu.busy//GUAGE"},
		{4, "miliao.xiaomi.bj/cpu.busy//GUAGE"},
		{5, "op.miliao.xiaomi.bj/cpu.busy//GUAGE"},
		{6, "micloud.miliao.xiaomi.bj/cpu.busy//GUAGE"},
		{5, "op_micloud.miliao.xiaomi.bj/cpu.busy//GUAGE"},
		{6, "op_micloud.miliao.xiaomi.bj/cpu.busy//GUAGE"},
	}

	testDps := []*falcon.DataPoint{
		{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9},
	}

	trigger := &Trigger{}

	testFillDps(shard, testItems, testDps, t)

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

	glog.V(4).Infof("=== host tag\n")
	for k, v := range trigger.HostTnodes {
		glog.V(5).Infof("%s\n", k)
		for _, v1 := range v {
			glog.V(5).Infof("    %s\n", v1.Name)
		}
	}

	glog.V(4).Infof("=== trigger item\n")
	cnt := [2]int{0, len(tagHosts)*2 - 2}
	for _, node := range trigger.TnodeIds {
		if len(node.ETriggerMetric) == 0 {
			continue
		}
		glog.V(5).Infof("tag[%s]\n", node.Name)
		for _, triggers := range node.ETriggerMetric {
			for _, trigger := range triggers {
				glog.V(5).Infof("    %s %s msg '%s' tags '%s'\n",
					trigger.Metric, trigger.Expr,
					trigger.Msg, trigger.Tags)
				for _, item := range trigger.items {
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
*/
