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

func (p *bucketEntry) addItem(item *Item) (ie *itemEntry, err error) {
	p.Lock()
	defer p.Unlock()

	if ie, err = itemEntryNew(item); err != nil {
		return
	}

	p.itemMap[string(item.Key)] = ie

	return ie, nil
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
