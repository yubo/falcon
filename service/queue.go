/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"sync"

	"github.com/yubo/gotool/list"
)

type queue struct {
	sync.RWMutex
	//size int
	head list.ListHead
}

func (p *queue) init() {
	//p.size = 0
	p.head.Init()
}

func (p *queue) addHead(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.Add(entry)
	//p.size++
}

func (p *queue) enqueue(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.AddTail(entry)
	//p.size++
}

func (p *queue) dequeue() *list.ListHead {
	p.Lock()
	defer p.Unlock()

	if p.head.Empty() {
		return nil
	}

	entry := p.head.Next
	entry.Del()
	//p.size--
	return entry
}

func (p *queue) first() *list.ListHead {
	p.Lock()
	defer p.Unlock()

	if p.head.Empty() {
		return nil
	}

	entry := p.head.Next
	return entry
}

func (p *queue) moveTail(entry *list.ListHead) {
	p.Lock()
	defer p.Unlock()

	p.head.MoveTail(entry)
	//p.size++
}
