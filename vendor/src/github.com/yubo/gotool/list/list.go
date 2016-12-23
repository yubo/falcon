/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package list

type ListHead struct {
	Next, Prev *ListHead
}

func (list *ListHead) Init() {
	list.Next = list
	list.Prev = list
}

func (head *ListHead) Empty() bool {
	return head.Next == head
}

func (head *ListHead) IsFirst(list *ListHead) bool {
	return list.Prev == head
}

func (head *ListHead) IsLast(list *ListHead) bool {
	return list.Next == head
}

func _list_del(entry *ListHead) {
	entry.Next.Prev = entry.Prev
	entry.Prev.Next = entry.Next
}

func (entry *ListHead) Del() {
	_list_del(entry)
	entry.Next = nil
	entry.Prev = nil
}

func _list_add(_new, prev, next *ListHead) {
	next.Prev = _new
	_new.Next = next
	_new.Prev = prev
	prev.Next = _new
}

func (entry *ListHead) DelInit() {
	_list_del(entry)
	entry.Init()
}

func (head *ListHead) Add(_new *ListHead) {
	_list_add(_new, head, head.Next)
}

func (head *ListHead) AddTail(_new *ListHead) {
	_list_add(_new, head.Prev, head)
}

func (head *ListHead) Move(list *ListHead) {
	_list_del(list)
	head.Add(list)
}

func (head *ListHead) MoveTail(list *ListHead) {
	_list_del(list)
	head.AddTail(list)
}

func _list_splice(list, prev, next *ListHead) {
	var first, last *ListHead
	if list.Empty() {
		return
	}

	first = list.Next
	last = list.Prev
	first.Prev = prev
	prev.Next = first
	last.Next = next
	next.Prev = last
}

func (head *ListHead) Splice(list *ListHead) {
	_list_splice(list, head, head.Next)
}

func (head *ListHead) SpliceTail(list *ListHead) {
	_list_splice(list, head, head.Next)
}

func (head *ListHead) SpliceInit(list *ListHead) {
	_list_splice(list, head, head.Next)
	list.Init()
}

func (head *ListHead) SpliceTailInit(list *ListHead) {
	_list_splice(list, head, head.Next)
	list.Init()
}
