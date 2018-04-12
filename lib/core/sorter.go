/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package core

import (
	"sort"
	"strings"
)

type SortInterface interface {
	Id() int
}

type sorter struct {
	data []SortInterface
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (s *sorter) Sort() {
	sort.Sort(s)
}

func (s *sorter) Len() int {
	return len(s.data)
}

func (s *sorter) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s *sorter) Less(i, j int) bool {
	p, q := s.data[i], s.data[j]
	if p.Id() < q.Id() {
		return true
	}
	if p.Id() > q.Id() {
		return false
	}
	switch strings.Compare(GetType(p), GetType(q)) {
	case -1:
		return true
	//case 1:
	//	return false
	default:
		return false
	}
}

func Sort(data []SortInterface) error {
	s := sorter{data: data}
	s.Sort()
	return nil
}

func moduleSort(data []Module) error {
	tmp := make([]SortInterface, len(data))
	for k, v := range data {
		tmp[k] = v.(SortInterface)
	}
	if err := Sort(tmp); err != nil {
		return err
	}
	for k, v := range tmp {
		data[k] = v.(Module)
	}
	return nil
}
