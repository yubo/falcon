/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"
	"sync"

	"github.com/golang/glog"
)

var (
	moduleCache [CTL_M_SIZE]cache
	cacheTree   *cacheTreeT
)

type cacheTreeT struct {
	sync.RWMutex
	m map[int64]*TreeNode
	c chan struct{}
}

type cache struct {
	sync.RWMutex
	enable bool
	id     map[int64]interface{}
	key    map[string]interface{}
}

func init() {
	for i := 0; i < CTL_M_SIZE; i++ {
		moduleCache[i] = cache{
			id:  make(map[int64]interface{}),
			key: make(map[string]interface{}),
		}
	}
}

func (c *cache) set(id int64, p interface{}, keys ...string) {
	if c.enable {
		c.Lock()
		defer c.Unlock()
		c.id[id] = p
		for _, key := range keys {
			c.key[key] = p
		}
	}
}

func (c *cache) get(id int64) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.id[id]
}

func (c *cache) getByKey(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.key[key]
}

func (c *cache) del(id int64, keys ...string) {
	if c.enable {
		c.Lock()
		defer c.Unlock()
		delete(c.id, id)
		for _, key := range keys {
			delete(c.key, key)
		}
	}
}

func (tree *cacheTreeT) get(id int64) *TreeNode {
	tree.Lock()
	defer tree.Unlock()
	glog.V(4).Infof("%s get cache tree %d/%d", MODULE_NAME, id, len(tree.m))
	return tree.m[id]
}

func (tree *cacheTreeT) build() {
	select {
	case tree.c <- struct{}{}:
	default:
	}
}

func (tree *cacheTreeT) _build() {
	var ns []TreeNode

	glog.V(10).Infof("%s cacheTree build entering", MODULE_NAME)

	_, err := Db.Ctrl.Raw("SELECT a.tag_id, a.sup_tag_id, b.name FROM tag_rel a JOIN tag b ON a.tag_id = b.id WHERE a.offset = 1 and b.type = 0 ORDER BY tag_id").QueryRows(&ns)
	if err != nil {
		return
	}

	tree.Lock()
	defer tree.Unlock()

	tree.m[1] = &TreeNode{
		TagId: 1,
		Name:  "",
		Label: "",
		Read:  true,
		//Operate: true,
	}

	for idx, _ := range ns {
		n := &ns[idx]
		n.Label = n.Name[strings.LastIndexAny(n.Name, ",")+1:]
		n.Read = true
		//n.Operate = true

		tree.m[n.TagId] = n
		if _, ok := tree.m[n.SupTagId]; ok {
			tree.m[n.SupTagId].Child = append(tree.m[n.SupTagId].Child, n)
		} else {
			glog.V(5).Infof("%s %s(%d) miss suptagid %d",
				MODULE_NAME, n.Name, n.TagId, n.SupTagId)
		}
	}

}
