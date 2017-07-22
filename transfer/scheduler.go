/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"stathat.com/c/consistent"
)

type scheduler interface {
	sched(string) chan *falcon.Item
	addChan(string, chan *falcon.Item) error
}

type schedConsistent struct {
	name    string
	consist *consistent.Consistent
	chans   map[string]chan *falcon.Item
}

func newSchedConsistent() *schedConsistent {
	sched := &schedConsistent{
		name:    "consistent",
		consist: consistent.New(),
		chans:   make(map[string]chan *falcon.Item),
	}
	sched.consist.NumberOfReplicas = falcon.REPLICAS
	return sched
}

func (s *schedConsistent) addChan(key string,
	ch chan *falcon.Item) error {
	s.consist.Add(key)
	s.chans[key] = ch
	return nil
}

func (s *schedConsistent) sched(key string) chan *falcon.Item {
	node, _ := s.consist.Get(key)
	glog.V(4).Infof(MODULE_NAME+"node %s", node)
	return s.chans[node]
}
