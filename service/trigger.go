/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/alarm"
	"golang.org/x/net/context"
)

type Node struct {
	Id                  int64
	Name                string
	ParentId            int64
	childs              []*Node
	eventTriggers       map[string]*EventTrigger   // map[EventTrigger.Name]
	eventTriggerMetrics map[string][]*EventTrigger // map[EventTrigger.Metric]
}

type TagHost struct {
	TagId    int64
	HostName string
}

type Trigger struct {
	db        orm.Ormer
	nodes     map[int64]*Node
	hostNodes map[string][]*Node // map[hostname]
}

type TriggerModule struct {
	sync.RWMutex

	eventChan chan *alarm.Event
	ctx       context.Context
	cancel    context.CancelFunc
	trigger   *Trigger
	service   *Service
}

func (p *TriggerModule) prestart(s *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.service = s
	return nil
}

func (p *TriggerModule) start(s *Service) (err error) {
	dbmaxidle, _ := s.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := s.Conf.Configer.Int(C_DB_MAX_CONN)
	syncInterval, _ := s.Conf.Configer.Int(C_SYNC_INTERVAL)
	judgeInterval, _ := s.Conf.Configer.Int(C_JUDGE_INTERVAL)
	judgeNum, _ := s.Conf.Configer.Int(C_JUDGE_NUM)
	eventChan := s.eventChan
	p.trigger = &Trigger{}

	db, err := falcon.NewOrm("service_sync",
		s.Conf.Configer.Str(C_SYNC_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	go p.syncWorker(syncInterval, db)
	go p.judgeWorker(judgeInterval, judgeNum, eventChan)

	return nil
}

func (p *TriggerModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *TriggerModule) reload(s *Service) error {
	return nil
}

func getNodes(db orm.Ormer) (rows []*Node, err error) {
	_, err = db.Raw("SELECT b.id, b.name, c.id as parent_id FROM tag_rel a JOIN tag b ON a.tag_id = b.id LEFT JOIN tag c ON a.sup_tag_id = c.id WHERE a.offset = 1 and b.type = 0 and c.type = 0 ORDER BY tag_id").QueryRows(&rows)
	return
}

func setNodes(rows []*Node, trigger *Trigger) (err error) {
	nodes := make(map[int64]*Node)
	nodes[1] = &Node{
		Id:                  1,
		Name:                "",
		eventTriggers:       make(map[string]*EventTrigger),
		eventTriggerMetrics: make(map[string][]*EventTrigger),
	}

	for _, row := range rows {

		row.eventTriggers = make(map[string]*EventTrigger)
		row.eventTriggerMetrics = make(map[string][]*EventTrigger)

		//nodes[row.Name] = row
		nodes[row.Id] = row
		if _, ok := nodes[row.ParentId]; ok {
			nodes[row.ParentId].childs = append(nodes[row.ParentId].childs, row)
		} else {
			glog.V(4).Infof("%s %s miss parent %d", MODULE_NAME,
				row.Name, row.ParentId)
		}
	}

	trigger.nodes = nodes
	return
}

func getTagHosts(db orm.Ormer) (rows []*TagHost, err error) {
	_, err = db.Raw("select a.tag_id, h.name as host_name from tag_host a left join host h on a.host_id = h.id  where tag_id ORDER BY a.tag_id").QueryRows(&rows)
	return
}

func setTagHosts(rows []*TagHost, trigger *Trigger) (err error) {
	nodes := trigger.nodes

	hostNodes := make(map[string][]*Node)
	for _, row := range rows {
		if _, ok := nodes[row.TagId]; ok {
			hostNodes[row.HostName] = append(hostNodes[row.HostName],
				nodes[row.TagId])
		} else {
			glog.V(4).Infof("%s miss tag %d\n", MODULE_NAME, row.TagId)
		}
	}

	trigger.hostNodes = hostNodes
	return
}

func triggerOverride(cover, base map[string]*EventTrigger) map[string]*EventTrigger {
	ret := make(map[string]*EventTrigger)
	for k, v := range base {
		ret[k] = v
	}
	for k, v := range cover {
		ret[k] = v
	}
	return ret
}

func adjustTrigger(tag *Node, base map[string]*EventTrigger) {
	tag.eventTriggers = triggerOverride(tag.eventTriggers, base)
	for _, child := range tag.childs {
		adjustTrigger(child, tag.eventTriggers)
	}
}

func adjustTriggerMetric(tag *Node) {
	for _, v := range tag.eventTriggers {
		tmp := *v
		tag.eventTriggerMetrics[v.Metric] = append(tag.eventTriggerMetrics[v.Metric], &tmp)
		for _, v1 := range v.Child {
			tmp := *v1
			tag.eventTriggerMetrics[v1.Metric] = append(tag.eventTriggerMetrics[v1.Metric], &tmp)
		}
	}
	for _, child := range tag.childs {
		adjustTriggerMetric(child)
	}
}

func getEventTriggers(db orm.Ormer) (rows []*EventTrigger, err error) {
	_, err = db.Raw("SELECT a.id AS id, a.parent_id AS parent_id, a.tag_id AS tag_id, a.priority AS priority, a.name AS name, a.metric AS metric, a.tags AS tags, a.expr AS expr, a.msg AS msg FROM event_trigger a LEFT JOIN event_trigger b ON a.parent_id = b.id where a.tag_id > 0 ORDER BY a.id").QueryRows(&rows)
	return
}

func setEventTriggers(rows []*EventTrigger, trigger *Trigger) (err error) {
	nodes := trigger.nodes
	index := make(map[int64]*EventTrigger)
	for _, row := range rows {
		if err := row.exprPrepare(); err != nil {
			continue
		}
		if row.ParentId > 0 {
			if _, ok := index[row.ParentId]; ok {
				index[row.ParentId].Child = append(index[row.ParentId].Child, row)
			} else {
				glog.V(4).Infof("%s %s %d miss parent %d", MODULE_NAME,
					row.Metric, row.Id, row.ParentId)
			}
		} else {
			index[row.Id] = row
			nodes[row.TagId].eventTriggers[row.Name] = row
		}
	}

	adjustTrigger(nodes[1], nil)
	adjustTriggerMetric(nodes[1])

	return nil
}

func setServiceItems(items map[string]*itemEntry, trigger *Trigger) error {
	for _, item := range items {
		// endpoint ?
		nodes, ok := trigger.hostNodes[string(item.endpoint)]
		if !ok {
			continue
		}

		for _, node := range nodes {
			// metric ?
			ets, ok := node.eventTriggerMetrics[string(item.metric)]
			if !ok {
				continue
			}

			// tags ?
			for _, et := range ets {
				if !tagsMatch(et.Tags, string(item.tags)) {
					continue
				}

				// match !
				et.items = append(et.items, item)
			}
		}

	}
	return nil
}

func setServiceShards(shard *ShardModule, trigger *Trigger) error {
	bucketMap := make(map[int32]*bucketEntry)

	// copy it first
	shard.RLock()
	for k, v := range shard.bucketMap {
		bucketMap[k] = v
	}
	shard.RUnlock()

	// process
	for _, v := range bucketMap {
		items := make(map[string]*itemEntry)
		v.RLock()
		for k1, v1 := range v.itemMap {
			items[k1] = v1
		}
		v.RUnlock()
		setServiceItems(items, trigger)
	}

	return nil
}

func (p *Trigger) updateNode() error {
	if nodes, err := getNodes(p.db); err != nil {
		glog.V(5).Infof("%v", err)
		return err
	} else {
		return setNodes(nodes, p)
	}
}

func (p *Trigger) updateTagHost() error {
	if tagHosts, err := getTagHosts(p.db); err != nil {
		glog.V(5).Infof("%v", err)
		return err
	} else {
		return setTagHosts(tagHosts, p)
	}
}

func (p *Trigger) updateEventTrigger() error {
	if eventTriggers, err := getEventTriggers(p.db); err != nil {
		glog.V(5).Infof("%v", err)
		return err
	} else {
		return setEventTriggers(eventTriggers, p)
	}
}

func (p *TriggerModule) triggerSync(db orm.Ormer) {
	glog.V(3).Infof("%s sync", MODULE_NAME)
	// nodes
	t := &Trigger{db: db}
	if t.updateNode() == nil &&
		t.updateTagHost() == nil &&
		t.updateEventTrigger() == nil &&
		setServiceShards(p.service.shard, t) == nil {

		p.Lock()
		p.trigger = t
		p.Unlock()
	}

}

func (p *TriggerModule) syncWorker(interval int, db orm.Ormer) {
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C
	p.triggerSync(db)

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			p.triggerSync(db)
		}
	}
}

func judgeTagNode(node *Node, eventChan chan *alarm.Event) {
	for _, eventTriggers := range node.eventTriggerMetrics {
		for _, eventTrigger := range eventTriggers {
			statsInc(ST_JUDGE_CNT, len(eventTrigger.items))
			for _, item := range eventTrigger.items {
				event := eventTrigger.Dispatch(item)
				if event != nil {
					select {
					case eventChan <- event:
						statsInc(ST_EVENT_CNT, 1)
					default:
						statsInc(ST_EVENT_DROP_CNT, 1)
					}
				}
			}
		}
	}
}

func (p *TriggerModule) judgeWorker(interval, n int, ch chan *alarm.Event) {
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C
	taskCh := make(chan *Node, n)

	for i := 0; i < n; i++ {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case node := <-taskCh:
					judgeTagNode(node, ch)
				}
			}

		}()
	}

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			glog.V(3).Infof("%s trigger judege worker entering", MODULE_NAME)

			p.RLock()
			t := p.trigger
			p.RUnlock()

			for _, node := range t.nodes {
				taskCh <- node
			}

		}
	}

}
