/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
	"golang.org/x/net/context"
)

type Action struct {
	users      []*User
	eventEntry *eventEntry
	flag       uint64
}

type User struct {
	Id    int64
	Name  string
	Cname string
	Email string
	Phone string
}

type TagRoleUser struct {
	TagId  int64
	RoleId int64
	UserId int64
}

type TagRoleToken struct {
	TagId   int64
	RoleId  int64
	TokenId int64
}

type Node struct {
	sync.RWMutex
	Name           string
	Id             int64
	ParentId       int64
	childs         []*Node
	actionTriggers []*ActionTrigger
	roleUser       map[int64]map[int64]bool // map[roleId][userId]
	roleToken      map[int64][]int64        // map[roleId]
	tokenUser      map[int64]map[int64]bool // map[tokenId][userId]
}

func (p *Node) processEvent(e *eventEntry) *ActionTrigger {
	p.Lock()
	defer p.Unlock()

	for _, t := range p.actionTriggers {
		if t.Dispatch(e) {
			return t
		}
	}
	return nil
}

func (p *Node) String() (s string) {
	s = fmt.Sprintf("%s\n", p.Name)
	for rid, users := range p.roleUser {
		s += fmt.Sprintf("\trole user %d [", rid)
		for uid, _ := range users {
			s += fmt.Sprintf("%d,", uid)
		}
		s = s[:len(s)-1]
		s += fmt.Sprintf("]\n")
	}
	for rid, ts := range p.roleToken {
		s += fmt.Sprintf("\trole token %d [", rid)
		for _, tid := range ts {
			s += fmt.Sprintf("%d,", tid)
		}
		s = s[:len(s)-1]
		s += fmt.Sprintf("]\n")
	}
	for tid, users := range p.tokenUser {
		s += fmt.Sprintf("\ttoken user %d [", tid)
		for uid, _ := range users {
			s += fmt.Sprintf("%d,", uid)
		}
		s = s[:len(s)-1]
		s += fmt.Sprintf("]\n")
	}
	for _, child := range p.childs {
		s += child.String()
	}
	return s
}

type Trigger struct {
	db    orm.Ormer
	nodes map[int64]*Node
	users map[int64]*User
}

type TriggerModule struct {
	sync.RWMutex
	lru             *queue //lru queue
	syncInterval    int
	workerProcesses int
	putChan         chan *Event
	actionChan      chan *Action
	cleanChan       chan *eventEntry
	ctx             context.Context
	cancel          context.CancelFunc
	trigger         *Trigger
	eventNodes      map[string]map[int64]*eventEntry // map[event.Key][nodeId]
	db              orm.Ormer
}

func (p *TriggerModule) prestart(a *Alarm) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.putChan = a.putEventChan
	p.actionChan = a.actionChan
	p.cleanChan = a.delEventEntryChan
	p.eventNodes = make(map[string]map[int64]*eventEntry)
	p.lru = &a.lru
	return nil
}

func (p *TriggerModule) start(s *Alarm) (err error) {
	dbmaxidle, _ := s.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := s.Conf.Configer.Int(C_DB_MAX_CONN)
	p.workerProcesses, _ = s.Conf.Configer.Int(C_WORKER_PROCESSES)
	p.syncInterval, _ = s.Conf.Configer.Int(C_SYNC_INTERVAL)
	p.putChan = s.putEventChan

	p.db, err = falcon.NewOrm("alarm_sync",
		s.Conf.Configer.Str(C_SYNC_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	go p.syncWorker()
	go p.processWorker()
	go p.cleanWorker()

	return nil
}

func getNodes(db orm.Ormer) (rows []*Node, err error) {
	_, err = db.Raw("SELECT b.id, b.name, c.id as parent_id FROM tag_rel a JOIN tag b ON a.tag_id = b.id LEFT JOIN tag c ON a.sup_tag_id = c.id WHERE a.offset = 1 and b.type = 0 and c.type = 0 ORDER BY tag_id").QueryRows(&rows)
	return
}

func setNodes(rows []*Node, trigger *Trigger) error {
	nodes := make(map[int64]*Node)
	nodes[1] = &Node{
		Id:        1,
		Name:      "",
		roleUser:  make(map[int64]map[int64]bool),
		roleToken: make(map[int64][]int64),
		tokenUser: make(map[int64]map[int64]bool),
	}

	for _, row := range rows {
		row.roleUser = make(map[int64]map[int64]bool)
		row.roleToken = make(map[int64][]int64)
		row.tokenUser = make(map[int64]map[int64]bool)
		nodes[row.Id] = row
		if _, ok := nodes[row.ParentId]; ok {
			nodes[row.ParentId].childs = append(nodes[row.ParentId].childs, row)
		} else {
			glog.V(4).Infof("%s %s miss parent %d", MODULE_NAME,
				row.Name, row.ParentId)
		}
	}
	trigger.nodes = nodes
	return nil
}

func getUsers(db orm.Ormer) (rows []*User, err error) {
	_, err = db.Raw("SELECT a.id, a.name, a.cname, a.email, a.phone FROM user a JOIN (SELECT distinct sub_id FROM tpl_rel WHERE type_id = ?) b ON a.id = b.sub_id", models.TPL_REL_T_ACL_TOKEN).QueryRows(&rows)
	return
}

func setUsers(rows []*User, trigger *Trigger) error {
	users := make(map[int64]*User)

	for _, row := range rows {
		users[row.Id] = row
	}
	trigger.users = users
	return nil
}

func getTagRoleUser(db orm.Ormer) (rows []*TagRoleUser, err error) {
	_, err = db.Raw("SELECT tag_id, tpl_id as role_id, sub_id as user_id FROM tpl_rel WHERE type_id = ?", models.TPL_REL_T_ACL_USER).QueryRows(&rows)
	return
}

// 用户身份可以继承
func adjustRoleUser(node *Node, base map[int64]map[int64]bool) {
	for role, users := range base {
		if _, ok := node.roleUser[role]; !ok {
			node.roleUser[role] = users
			continue
		}
		// merge
		for user, _ := range users {
			node.roleUser[role][user] = true
		}
	}
	for _, child := range node.childs {
		adjustRoleUser(child, node.roleUser)
	}
}

func setTagRoleUser(rows []*TagRoleUser, trigger *Trigger) error {
	nodes := trigger.nodes
	for _, row := range rows {
		node, ok := nodes[row.TagId]
		if !ok {
			continue
		}
		if _, ok := node.roleUser[row.RoleId]; !ok {
			node.roleUser[row.RoleId] = make(map[int64]bool)
		}
		node.roleUser[row.RoleId][row.UserId] = true
	}
	adjustRoleUser(nodes[1], nil)
	return nil
}

func getTagRoleToken(db orm.Ormer) (rows []*TagRoleToken, err error) {
	_, err = db.Raw("SELECT tag_id, tpl_id as role_id, sub_id as token_id FROM tpl_rel WHERE type_id = ?", models.TPL_REL_T_ACL_TOKEN).QueryRows(&rows)
	return
}

func adjustRoleToken(node *Node, base map[int64][]int64) {
	for k, v := range base {
		if _, ok := node.roleToken[k]; !ok {
			node.roleToken[k] = v
		}
		// else override
	}
	for _, child := range node.childs {
		adjustRoleToken(child, node.roleToken)
	}
}

func setTagRoleToken(rows []*TagRoleToken, trigger *Trigger) error {
	nodes := trigger.nodes
	for _, row := range rows {
		node, ok := nodes[row.TagId]
		if !ok {
			continue
		}
		node.roleToken[row.RoleId] = append(node.roleToken[row.RoleId], row.TokenId)
	}
	adjustRoleToken(nodes[1], nil)
	return nil
}

func (trigger *Trigger) updateTagTokenUser() error {
	for _, node := range trigger.nodes {
		// token
		for roleId, tokens := range node.roleToken {
			// user
			for userId, _ := range node.roleUser[roleId] {
				for _, tokenId := range tokens {
					if _, ok := node.tokenUser[tokenId]; !ok {
						node.tokenUser[tokenId] = make(map[int64]bool)
					}
					node.tokenUser[tokenId][userId] = true
				}
			}
		}
	}
	return nil
}

func getActionTriggers(db orm.Ormer) (rows []*ActionTrigger, err error) {
	_, err = db.Raw("SELECT id, tag_id, token_id, order_id, expr, action_flag, action_script FROM action_trigger ORDER BY order_id, id").QueryRows(&rows)
	return
}

func adjustTrigger(node *Node, base []*ActionTrigger) {
	node.actionTriggers = append(node.actionTriggers, base...)
	for _, child := range node.childs {
		adjustTrigger(child, node.actionTriggers)
	}
}

func setActionTriggers(rows []*ActionTrigger, trigger *Trigger) (err error) {
	nodes := trigger.nodes
	for _, row := range rows {
		if err := row.exprPrepare(); err != nil {
			continue
		}
		nodes[row.TagId].actionTriggers = append(nodes[row.TagId].actionTriggers, row)
	}

	adjustTrigger(nodes[1], nil)
	return nil
}

func (p *Trigger) updateNode() error {
	if nodes, err := getNodes(p.db); err != nil {
		return err
	} else {
		return setNodes(nodes, p)
	}
}

func (p *Trigger) updateUser() error {
	if users, err := getUsers(p.db); err != nil {
		return err
	} else {
		return setUsers(users, p)
	}
}

func (p *Trigger) updateTagRoleUser() error {
	if tagRoleUser, err := getTagRoleUser(p.db); err != nil {
		return err
	} else {
		return setTagRoleUser(tagRoleUser, p)
	}
}

func (p *Trigger) updateTagRoleToken() error {
	if tagRoleToken, err := getTagRoleToken(p.db); err != nil {
		return err
	} else {
		return setTagRoleToken(tagRoleToken, p)
	}
}

func (p *Trigger) updateActionTrigger() error {
	if rows, err := getActionTriggers(p.db); err != nil {
		return err
	} else {
		return setActionTriggers(rows, p)
	}

}

func (p *TriggerModule) triggerSync() {
	glog.V(3).Infof("%s sync", MODULE_NAME)
	// nodes
	t := &Trigger{db: p.db}

	if t.updateNode() == nil &&
		t.updateUser() == nil &&
		t.updateTagRoleUser() == nil &&
		t.updateTagRoleToken() == nil &&
		t.updateTagTokenUser() == nil &&
		t.updateActionTrigger() == nil {

		p.Lock()
		p.trigger = t
		p.Unlock()
	}
}

func (p *TriggerModule) syncWorker() {
	ticker := time.NewTicker(time.Second * time.Duration(p.syncInterval)).C
	p.triggerSync()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker:
			p.triggerSync()
		}
	}
}

func actionGenerate(p *ActionTrigger, node *Node, e *eventEntry, users map[int64]*User) *Action {

	if p == nil || p.ActionFlag == 0 || len(node.tokenUser[p.TokenId]) == 0 {
		return nil
	}

	action := &Action{flag: p.ActionFlag, eventEntry: e}

	for uid, _ := range node.tokenUser[p.TokenId] {
		if user, ok := users[uid]; ok {
			action.users = append(action.users, user)
		}
	}
	return action
}

func (p *TriggerModule) createEvent(event *Event) {
	e := &eventEntry{
		lastTs:    time.Now().Unix(),
		tagId:     event.TagId,
		key:       string(event.Key),
		expr:      string(event.Expr),
		msg:       string(event.Msg),
		timestamp: event.Timestamp,
		value:     event.Value,
		priority:  int(event.Priority),
	}
	p.RLock()
	t := p.trigger
	p.RUnlock()

	// check node
	node := t.nodes[e.tagId]
	if node == nil {
		return
	}

	// add eventEntry
	if p.eventNodes[e.key] == nil {
		p.eventNodes[e.key] = make(map[int64]*eventEntry)
	}
	p.eventNodes[e.key][e.tagId] = e

	// add lru queue
	p.lru.enqueue(&e.list)

	// generate action
	if action := actionGenerate(node.processEvent(e),
		node, e, t.users); action != nil {
		p.actionChan <- action
	}

}

func (p *TriggerModule) updateEvent(entry *eventEntry) {
	p.lru.Lock()
	entry.Lock()
	defer entry.Unlock()
	defer p.lru.Unlock()

	entry.lastTs = time.Now().Unix()
	entry.list.Del()
	p.lru.head.AddTail(&entry.list)
}

func (p *TriggerModule) deleteEvent(e *eventEntry) {
	// entry alread del from lru list
	p.Lock()
	defer p.Unlock()

	delete(p.eventNodes[e.key], e.tagId)
	if len(p.eventNodes[e.key]) == 0 {
		delete(p.eventNodes, e.key)
	}
}

func (p *TriggerModule) precessEvent(event *Event) {
	key := string(event.Key)

	p.RLock()
	nodes := p.eventNodes[key]
	p.RUnlock()

	if nodes == nil || nodes[event.TagId] == nil {
		p.createEvent(event)
	} else {
		p.updateEvent(nodes[event.TagId])
	}
}

func (p *TriggerModule) processWorker() {
	ch := p.putChan
	for i := 0; i < p.workerProcesses; i++ {
		go func() {
			for {
				select {
				case <-p.ctx.Done():
					return
				case event := <-ch:
					p.precessEvent(event)
				}
			}

		}()
	}
}

func (p *TriggerModule) cleanWorker() {
	ch := p.cleanChan
	for {
		select {
		case <-p.ctx.Done():
			return
		case entry := <-ch:
			p.deleteEvent(entry)
		}
	}
}

func (p *TriggerModule) stop(a *Alarm) error {
	p.cancel()
	return nil
}

func (p *TriggerModule) reload(a *Alarm) error {
	return nil
}
