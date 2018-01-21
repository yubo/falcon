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
	dpEntryMap map[string]*dpEntry
}

func (p *bucketEntry) createDpEntry(dp *DataPoint) (ie *dpEntry, err error) {
	p.Lock()
	defer p.Unlock()

	if ie, err = dpEntryNew(dp); err != nil {
		return
	}

	p.dpEntryMap[string(dp.Key.Key)] = ie

	return ie, nil
}

func (p *bucketEntry) getDpEntry(key string) (*dpEntry, error) {
	p.RLock()
	defer p.RUnlock()

	if ie, ok := p.dpEntryMap[key]; ok {
		return ie, nil
	}
	return nil, falcon.ErrNoExits
}

/*
 * not idxq.size --
 */

func (p *bucketEntry) unlink(key string) *dpEntry {
	p.Lock()
	defer p.Unlock()

	ie, ok := p.dpEntryMap[key]
	if !ok {
		return nil
	}
	ie.Lock()
	defer ie.Unlock()

	delete(p.dpEntryMap, key)
	ie.list.Del()

	return ie
}
