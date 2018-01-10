/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
	"golang.org/x/net/context"
)

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

type TreeNode struct {
	Id             int64
	Name           string
	ParentId       int64
	childs         []*TreeNode
	actionTriggers []*ActionTrigger
	roleUser       map[int64][]int64 // map[roleId]
	roleToken      map[int64][]int64 // map[roleId]
}

type Trigger struct {
	nodes      map[int64]*TreeNode
	users      map[int64]*User
	eventNodes map[string][]*TreeNode // map[event.Key]
}

type TriggerModule struct {
	sync.RWMutex
	syncInterval int
	putChan      chan *Event
	ctx          context.Context
	cancel       context.CancelFunc
	trigger      *Trigger
	db           orm.Ormer
}

func (p *TriggerModule) prestart(a *Alarm) error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.putChan = a.appPutChan
	return nil
}

func (p *TriggerModule) start(s *Alarm) (err error) {
	dbmaxidle, _ := s.Conf.Configer.Int(C_DB_MAX_IDLE)
	dbmaxconn, _ := s.Conf.Configer.Int(C_DB_MAX_CONN)
	p.syncInterval, _ = s.Conf.Configer.Int(C_SYNC_INTERVAL)
	p.putChan = s.appPutChan

	p.db, err = falcon.NewOrm("alarm_sync",
		s.Conf.Configer.Str(C_SYNC_DSN), dbmaxidle, dbmaxconn)
	if err != nil {
		return err
	}

	go p.syncWorker()

	return nil
}

func getTreeNodes(db orm.Ormer) (rows []*TreeNode, err error) {
	_, err = db.Raw("SELECT b.id, b.name, c.id as parent_id FROM tag_rel a JOIN tag b ON a.tag_id = b.id LEFT JOIN tag c ON a.sup_tag_id = c.id WHERE a.offset = 1 and b.type = 0 and c.type = 0 ORDER BY tag_id").QueryRows(&rows)
	return
}

func setTreeNodes(rows []*TreeNode, trigger *Trigger) {
	nodes := make(map[int64]*TreeNode)
	nodes[1] = &TreeNode{
		Id:   1,
		Name: "",
	}

	for _, row := range rows {

		nodes[row.Id] = row
		if _, ok := nodes[row.ParentId]; ok {
			nodes[row.ParentId].childs = append(nodes[row.ParentId].childs, row)
		} else {
			glog.V(4).Infof(MODULE_NAME+"%s miss parent %d",
				row.Name, row.ParentId)
		}
	}

	trigger.nodes = nodes
}

func getUsers(db orm.Ormer) (rows []*User, err error) {
	_, err = db.Raw("select a.id, a.name, a.cname, a.email, a.phone from user a join (select distinct sub_id from tpl_rel where type_id = ?) b on a.id = b.sub_id", models.TPL_REL_T_ACL_TOKEN).QueryRows(&rows)
	return
}

func setUsers(rows []*User, trigger *Trigger) {
	users := make(map[int64]*User)

	for _, row := range rows {
		users[row.Id] = row
	}

	trigger.users = users
}

func getTagRoleUser(db orm.Ormer) (rows []*TagRoleUser, err error) {
	_, err = db.Raw("select tag_id, tpl_id as role_id, sub_id as user_id from tpl_rel where type_id = ?", models.TPL_REL_T_ACL_USER).QueryRows(&rows)
	return
}

func adjustRoleUser(node *TreeNode, base map[int64][]int64) {
	for k, v := range base {
		if _, ok := node.roleUser[k]; !ok {
			node.roleUser[k] = v
		}
	}
	for _, child := range node.childs {
		adjustRoleUser(child, node.roleUser)
	}
}

func setTagRoleUser(rows []*TagRoleUser, trigger *Trigger) {
	nodes := trigger.nodes
	for _, row := range rows {
		node, ok := nodes[row.TagId]
		if !ok {
			continue
		}
		node.roleUser[row.UserId] = append(node.roleUser[row.UserId], row.UserId)
	}

	adjustRoleUser(nodes[1], nil)
}

func getTagRoleToken(db orm.Ormer) (rows []*TagRoleToken, err error) {
	_, err = db.Raw("select tag_id, tpl_id as role_id, sub_id as token_id from tpl_rel where type_id = ?", models.TPL_REL_T_ACL_TOKEN).QueryRows(&rows)
	return
}

func adjustRoleToken(node *TreeNode, base map[int64][]int64) {
	for k, v := range base {
		if _, ok := node.roleToken[k]; !ok {
			node.roleToken[k] = v
		}
	}
	for _, child := range node.childs {
		adjustRoleUser(child, node.roleUser)
	}
}

func setTagRoleToken(rows []*TagRoleToken, trigger *Trigger) {
	nodes := trigger.nodes
	for _, row := range rows {
		node, ok := nodes[row.TagId]
		if !ok {
			continue
		}
		node.roleToken[row.TokenId] = append(node.roleToken[row.TokenId], row.TokenId)
	}
	adjustRoleToken(nodes[1], nil)
}

func getActionTriggers(db orm.Ormer) (rows []*ActionTrigger, err error) {
	_, err = db.Raw("SELECT id, tag_id, token_id, order_id, expr, action_flag, action_script FROM action_trigger ORDER BY order_id, id").QueryRows(&rows)
	return
}

func adjustTrigger(node *TreeNode, base []*ActionTrigger) {
	node.actionTriggers = append(node.actionTriggers, base...)
	for _, child := range node.childs {
		adjustTrigger(child, node.actionTriggers)
	}
}

func setActionTriggers(rows []*ActionTrigger, trigger *Trigger) (err error) {
	nodes := trigger.nodes
	for _, row := range rows {
		nodes[row.TagId].actionTriggers = append(nodes[row.TagId].actionTriggers, row)
	}

	adjustTrigger(nodes[1], nil)
	return nil
}

func (p *TriggerModule) syncWorker() {
	var (
		err            error
		trigger        *Trigger
		treeNodes      []*TreeNode
		users          []*User
		tagRoleUser    []*TagRoleUser
		tagRoleToken   []*TagRoleToken
		actionTriggers []*ActionTrigger
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
			setTreeNodes(treeNodes, trigger)

			if users, err = getUsers(p.db); err != nil {
				continue
			}
			setUsers(users, trigger)

			if tagRoleUser, err = getTagRoleUser(p.db); err != nil {
				continue
			}
			setTagRoleUser(tagRoleUser, trigger)

			if tagRoleToken, err = getTagRoleToken(p.db); err != nil {
				continue
			}
			setTagRoleToken(tagRoleToken, trigger)

			if actionTriggers, err = getActionTriggers(p.db); err != nil {
				continue
			}
			if err = setActionTriggers(actionTriggers, trigger); err != nil {
				glog.V(3).Info("%s set actionTriggers error %s", MODULE_NAME, err)
				continue
			}

			p.Lock()
			p.trigger = trigger
			p.Unlock()

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
