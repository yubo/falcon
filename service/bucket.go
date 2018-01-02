/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"sync"

	"github.com/yubo/falcon"
)

type bucketEntry struct { // bucket_t
	sync.RWMutex
	itemMap map[string]*itemEntry
}

func (p *bucketEntry) addItem(item *falcon.Item) (*itemEntry, error) {
	p.Lock()
	defer p.Unlock()

	e := &itemEntry{
		createTs:  timer.now(),
		shardId:   item.ShardId,
		endpoint:  item.Endpoint,
		metric:    item.Metric,
		tags:      item.Tags,
		typ:       item.Type,
		timestamp: make([]int64, CACHE_SIZE),
		value:     make([]float64, CACHE_SIZE),
	}

	p.itemMap[item.Key()] = e

	return e, nil
}

func (p *bucketEntry) getItem(key string) (*itemEntry, error) {
	p.RLock()
	defer p.RUnlock()

	if ie, ok := p.itemMap[key]; ok {
		return ie, nil
	}
	return nil, falcon.ErrNoExits
}

/*
 * not idxq.size --
 */

func (p *bucketEntry) unlink(key string) *itemEntry {
	p.Lock()
	defer p.Unlock()

	ie, ok := p.itemMap[key]
	if !ok {
		return nil
	}
	ie.Lock()
	defer ie.Unlock()

	delete(p.itemMap, key)
	ie.list.Del()

	return ie
}
