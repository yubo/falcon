/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package core

import (
	"testing"

	"github.com/golang/glog"
)

type bar struct {
	id int
}

func (m *bar) Id() int {
	return m.id
}

type foo struct {
	id int
}

func (m *foo) Id() int {
	return m.id
}

func iSort(data []interface{}) error {
	tmp := make([]SortInterface, len(data))
	for k, v := range data {
		tmp[k] = v.(SortInterface)
	}
	if err := Sort(tmp); err != nil {
		return err
	}
	for k, v := range tmp {
		data[k] = v
	}
	return nil
}
func TestSort(t *testing.T) {
	data := []interface{}{
		&foo{2},
		&bar{2},
		&foo{3},
		&foo{4},
		&bar{1},
	}
	iSort(data)
	for _, m := range data {
		glog.V(3).Infof("%d %s\n",
			m.(SortInterface).Id(), GetType(m))
	}

}
