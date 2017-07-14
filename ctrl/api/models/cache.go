/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl"
	"github.com/yubo/falcon/ctrl/config"
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
	glog.V(4).Infof(MODULE_NAME+"get cache tree %d/%d", id, len(tree.m))
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

	glog.V(4).Infof(MODULE_NAME + "cacheTree build entering")

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
			glog.V(4).Infof(MODULE_NAME+"%s(%d) miss suptagid %d",
				n.Name, n.TagId, n.SupTagId)
		}
	}

}

/* called by initModels() */
func initCache(c *config.ConfCtrl) error {
	for _, module := range strings.Split(
		c.Ctrl.Str(ctrl.C_CACHE_MODULE), ",") {
		for k, v := range moduleName {
			if v == module {
				moduleCache[k] = cache{
					enable: true,
					id:     make(map[int64]interface{}),
					key:    make(map[string]interface{}),
				}
				break
			}
		}
	}

	// build host cache
	if moduleCache[CTL_M_HOST].enable {
		var items []*Host
		for i := 0; ; i++ {
			n, err := Db.Ctrl.Raw("select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host limit ? offset ?", 100, 100*i).QueryRows(&items)
			if err != nil || n == 0 {
				break
			}

			for _, h := range items {
				moduleCache[CTL_M_HOST].set(h.Id, h, h.Name)
			}
		}
	}
	// build tag cache
	if moduleCache[CTL_M_TAG].enable {
		var items []*Tag
		for i := 0; ; i++ {
			n, err := Db.Ctrl.Raw("select id, name, type, create_time from tag LIMIT ? OFFSET ?", 100, 100*i).QueryRows(&items)
			if err != nil || n == 0 {
				break
			}

			for _, t := range items {
				moduleCache[CTL_M_TAG].set(t.Id, t, t.Name)
			}
		}
	}

	go func() {
		cacheTree = &cacheTreeT{
			m: make(map[int64]*TreeNode),
			c: make(chan struct{}, 1),
		}
		cacheTree._build()
		ticker := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ticker.C:
				cacheTree._build()
			case <-cacheTree.c:
				cacheTree._build()
			}
		}
	}()

	return nil
}
