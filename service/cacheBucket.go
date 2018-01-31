/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"sync"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/tsdb"
)

const (
	CACHE_BUCKET_DISABLE = iota
	CACHE_BUCKET_ENABLE
)

type cacheBucket struct { // bucket_t
	sync.RWMutex
	state   int
	entries map[string]*cacheEntry
}

func (p *cacheBucket) getState() int {
	p.RLock()
	defer p.RUnlock()
	return p.state
}

func (p *cacheBucket) setState(state int) error {
	if state < CACHE_BUCKET_DISABLE || state > CACHE_BUCKET_ENABLE {
		return falcon.EINVAL
	}
	p.Lock()
	p.state = state
	p.Unlock()
	return nil
}

func (p *cacheBucket) createCacheEntry(dp *tsdb.DataPoint) (e *cacheEntry, err error) {
	p.Lock()
	defer p.Unlock()

	if e, err = newCacheEntry(dp); err != nil {
		return
	}

	p.entries[string(dp.Key.Key)] = e

	return e, nil
}

func (p *cacheBucket) getCacheEntry(key string) (*cacheEntry, error) {
	p.RLock()
	defer p.RUnlock()

	if e, ok := p.entries[key]; ok {
		return e, nil
	}
	return nil, falcon.ErrNoExits
}

func (p *cacheBucket) delEntry(key string, q *queue) *cacheEntry {
	p.Lock()

	e, ok := p.entries[key]
	if !ok {
		p.Unlock()
		return nil
	}
	delete(p.entries, key)
	p.Unlock()

	q.del(&e.list)

	return e
}

func (p *cacheBucket) _delEntry(key string, q *queue) *cacheEntry {

	e, ok := p.entries[key]
	if !ok {
		return nil
	}
	delete(p.entries, key)

	q.del(&e.list)

	return e
}

func (p *cacheBucket) clean(timeout int64, q *queue) {
	if p.getState() == CACHE_BUCKET_ENABLE {
		// clean expire entry
		now := timer.now()
		p.RLock()
		keys := make([]string, 0)
		for key, e := range p.entries {
			if now-e.lastTs > timeout {
				keys = append(keys, key)
			}
		}
		p.RUnlock()

		p.Lock()
		for _, key := range keys {
			p._delEntry(key, q)

		}
		p.Unlock()
	} else {
		// clean all entry
		p.Lock()
		for key, _ := range p.entries {
			p._delEntry(key, q)
		}
		p.Unlock()
	}
}
