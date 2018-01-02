/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

type event struct {
}

type EventTrigger struct {
	Id       int64
	ParentId int64
	TagId    int64
	Priority int
	Name     string
	Metric   string
	Tags     string
	Func     string
	Op       string
	Value    string
	Expr     string
	Msg      string
	Child    []*EventTrigger
	items    []*itemEntry
}

func (p *EventTrigger) Exec(item *itemEntry) *event {
	return nil
}

/*
expr:
true
all(#3) > hj
true && false

// all(#3) > 3
// min(#4) > 3
// max(#5) > 3
//
func (p *EventTrigger) Parse() error {
	return nil
}
*/
