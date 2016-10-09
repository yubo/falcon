/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

import (
	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
	"stathat.com/c/consistent"
)

type lbScheduler interface {
	sched(string) chan *specs.MetaData
	addChan(string, chan *specs.MetaData) error
}

type schedConsistent struct {
	name    string
	consist *consistent.Consistent
	chans   map[string]chan *specs.MetaData
}

func newSchedConsistent() *schedConsistent {
	sched := &schedConsistent{
		name:    "consistent",
		consist: consistent.New(),
		chans:   make(map[string]chan *specs.MetaData),
	}
	sched.consist.NumberOfReplicas = specs.REPLICAS
	return sched
}

func (s *schedConsistent) addChan(key string,
	ch chan *specs.MetaData) error {
	s.consist.Add(key)
	s.chans[key] = ch
	return nil
}

func (s *schedConsistent) sched(key string) chan *specs.MetaData {
	node, _ := s.consist.Get(key)
	glog.V(4).Infof(MODULE_NAME+"node %s", node)
	return s.chans[node]
}
