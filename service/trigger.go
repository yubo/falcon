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
	"golang.org/x/net/context"
)

type TreeNode struct {
	Id             int64
	Name           string
	ParentId       int64
	Child          []*TreeNode
	ETrigger       map[string]*EventTrigger   // map[EventTrigger.Name]
	ETriggerMetric map[string][]*EventTrigger // map[EventTrigger.Metric]
}

type TagHost struct {
	TagId    int64
	HostName string
}

type Trigger struct {
	TnodeIds   map[int64]*TreeNode
	Tnodes     map[string]*TreeNode
	HostTnodes map[string][]*TreeNode // map[hostname]
}

type TriggerModule struct {
	sync.RWMutex

	syncInterval  int
	judgeInterval int
	judgeNum      int
	alarmNum      int
	eventCh       chan *event

	ctx     context.Context
	cancel  context.CancelFunc
	trigger *Trigger
	service *Service
	db      orm.Ormer
}

type action struct {
}

func (p *TriggerModule) prestart(s *Service) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.service = s
	return nil
}

func (p *TriggerModule) start(s *Service) (err error) {
	dbmaxidle, _ := s.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := s.Conf.Configer.Int(C_DB_MAX_CONN)
	p.syncInterval, _ = s.Conf.Configer.Int(C_SYNC_INTERVAL)
	p.judgeInterval, _ = s.Conf.Configer.Int(C_JUDGE_INTERVAL)
	p.judgeNum, _ = s.Conf.Configer.Int(C_JUDGE_NUM)
	p.alarmNum, _ = s.Conf.Configer.Int(C_ALARM_NUM)
	p.eventCh = make(chan *event, 32)

	p.db, err = falcon.NewOrm("service_sync",
		s.Conf.Configer.Str(C_SYNC_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	go p.syncWorker()
	go p.judgeWorker()
	go p.alarmWorker()

	return nil
}

func getTreeNodes(db orm.Ormer) (rows []*TreeNode, err error) {
	_, err = db.Raw("SELECT b.id, b.name, c.id as parent_id FROM tag_rel a JOIN tag b ON a.tag_id = b.id LEFT JOIN tag c ON a.sup_tag_id = c.id WHERE a.offset = 1 and b.type = 0 and c.type = 0 ORDER BY tag_id").QueryRows(&rows)
	return
}

func setTreeNodes(rows []*TreeNode, trigger *Trigger) (err error) {
	tagIds := make(map[int64]*TreeNode)
	tags := make(map[string]*TreeNode)
	tags[""] = &TreeNode{
		Id:             1,
		Name:           "",
		ETrigger:       make(map[string]*EventTrigger),
		ETriggerMetric: make(map[string][]*EventTrigger),
	}
	tagIds[1] = tags[""]

	for _, row := range rows {

		row.ETrigger = make(map[string]*EventTrigger)
		row.ETriggerMetric = make(map[string][]*EventTrigger)

		tags[row.Name] = row
		tagIds[row.Id] = row
		if _, ok := tagIds[row.ParentId]; ok {
			tagIds[row.ParentId].Child = append(tagIds[row.ParentId].Child, row)
		} else {
			glog.V(4).Infof(MODULE_NAME+"%s miss parent %d",
				row.Name, row.ParentId)
		}
	}

	trigger.Tnodes = tags
	trigger.TnodeIds = tagIds
	return
}

func getTagHosts(db orm.Ormer) (rows []*TagHost, err error) {
	_, err = db.Raw("select a.tag_id, h.name as host_name from tag_host a left join host h on a.host_id = h.id  where tag_id ORDER BY a.tag_id").QueryRows(&rows)
	return
}

func setTagHosts(rows []*TagHost, trigger *Trigger) (err error) {
	tags := trigger.TnodeIds

	hostTags := make(map[string][]*TreeNode)
	for _, row := range rows {
		if _, ok := tags[row.TagId]; ok {
			hostTags[row.HostName] = append(hostTags[row.HostName],
				tags[row.TagId])
		} else {
			glog.V(4).Infof(MODULE_NAME+"miss tag %d\n", row.TagId)
		}
	}

	trigger.HostTnodes = hostTags
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

func reduceTrigger(tag *TreeNode, base map[string]*EventTrigger) {
	tag.ETrigger = triggerOverride(tag.ETrigger, base)
	for _, child := range tag.Child {
		reduceTrigger(child, tag.ETrigger)
	}
}

func reduceTriggerMetric(tag *TreeNode) {
	for _, v := range tag.ETrigger {
		tmp := *v
		tag.ETriggerMetric[v.Metric] = append(tag.ETriggerMetric[v.Metric], &tmp)
		for _, v1 := range v.Child {
			tmp := *v1
			tag.ETriggerMetric[v1.Metric] = append(tag.ETriggerMetric[v1.Metric], &tmp)
		}
	}
	for _, child := range tag.Child {
		reduceTriggerMetric(child)
	}
}

func getEventTriggers(db orm.Ormer) (rows []*EventTrigger, err error) {
	_, err = db.Raw("SELECT a.id AS id, a.parent_id AS parent_id, a.tag_id AS tag_id, a.priority AS priority, a.name AS name, a.metric AS metric, a.tags AS tags, a.expr AS expr, a.msg AS msg FROM event_trigger a LEFT JOIN event_trigger b ON a.parent_id = b.id where a.tag_id > 0 ORDER BY a.id").QueryRows(&rows)
	return
}

func setEventTriggers(rows []*EventTrigger, trigger *Trigger) (err error) {
	tags := trigger.TnodeIds
	index := make(map[int64]*EventTrigger)
	for _, row := range rows {
		if err := row.exprPrepare(); err != nil {
			continue
		}
		if row.ParentId > 0 {
			if _, ok := index[row.ParentId]; ok {
				index[row.ParentId].Child = append(index[row.ParentId].Child, row)
			} else {
				glog.V(4).Infof(MODULE_NAME+"%s %d miss parent %d",
					row.Metric, row.Id, row.ParentId)
			}
		} else {
			index[row.Id] = row
			tags[row.TagId].ETrigger[row.Name] = row
		}
	}

	reduceTrigger(tags[1], nil)
	reduceTriggerMetric(tags[1])

	return nil
}

func setServiceItems(items map[string]*itemEntry, trigger *Trigger) error {
	for _, item := range items {
		// endpoint ?
		nodes, ok := trigger.HostTnodes[string(item.endpoint)]
		if !ok {
			continue
		}

		for _, node := range nodes {
			// metric ?
			ets, ok := node.ETriggerMetric[string(item.metric)]
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

func (p *TriggerModule) syncWorker() {
	var (
		err           error
		trigger       *Trigger
		treeNodes     []*TreeNode
		tagHosts      []*TagHost
		eventTriggers []*EventTrigger
	)
	ticker := time.NewTicker(time.Second * time.Duration(p.syncInterval)).C
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			glog.V(3).Info(MODULE_NAME + "sync")
			// tree nodes
			trigger = &Trigger{}

			if treeNodes, err = getTreeNodes(p.db); err != nil {
				continue
			}
			if err = setTreeNodes(treeNodes, trigger); err != nil {
				glog.V(4).Info("%s set treeNodes error %s", MODULE_NAME, err)
				continue
			}

			if tagHosts, err = getTagHosts(p.db); err != nil {
				continue
			}
			if err = setTagHosts(tagHosts, trigger); err != nil {
				glog.V(3).Info("%s set tagHosts error %s", MODULE_NAME, err)
				continue
			}

			if eventTriggers, err = getEventTriggers(p.db); err != nil {
				continue
			}
			if err = setEventTriggers(eventTriggers, trigger); err != nil {
				glog.V(3).Info("%s set eventTriggers error %s", MODULE_NAME, err)
				continue
			}

			if err = setServiceShards(p.service.shard, trigger); err != nil {
				glog.V(3).Info("%s mount item error %s", MODULE_NAME, err)
				continue
			}

			p.Lock()
			p.trigger = trigger
			p.Unlock()

		}
	}

}

func judgeTagNode(node *TreeNode, eventCh chan *event) {
	for _, triggers := range node.ETriggerMetric {
		for _, trigger := range triggers {
			for _, item := range trigger.items {
				event := trigger.Exec(item)
				if event != nil {
					eventCh <- event
				}
			}
		}
	}
}

func (p *TriggerModule) judgeWorker() {
	ticker := time.NewTicker(time.Second * time.Duration(p.judgeInterval)).C
	taskCh := make(chan *TreeNode, p.judgeNum)

	for i := 0; i < p.judgeNum; i++ {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case node := <-taskCh:
					judgeTagNode(node, p.eventCh)
				}
			}

		}()
	}

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			glog.V(3).Info("%s trigger judege worker entering", MODULE_NAME)

			p.RLock()
			trigger := p.trigger
			p.RUnlock()

			for _, node := range trigger.TnodeIds {
				taskCh <- node
			}

		}
	}

}

func (p *TriggerModule) alarmWorker() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case e := <-p.eventCh:
			glog.V(3).Info("%s recv event %#v", MODULE_NAME, e)
		}

	}
}

func (p *TriggerModule) stop(s *Service) error {
	p.cancel()
	return nil
}

func (p *TriggerModule) reload(s *Service) error {
	return nil
}
