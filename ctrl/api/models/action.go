/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

/*
 * 命中action trigger 后，会导致event短路，剩余的trigger将被忽略
 * action trigger 优先级
 *  - 路径越大，优先级越高
 *  - 相同路径长度下，priority越小，优先级越高
 * 匹配规则
 *  - metric 留空，匹配所有情况， 不为空时，字符串比较
 *  - tags 留空，匹配所有情况， 不为空时，比较tag k v
 */
type ActionTrigger struct {
	Id        int64  `json:"id"`
	TagId     int64  `json:"tag_id"`
	Priority  int    `json:"priority"`
	Metric    string `json:"metric"`
	Tags      string `json:"tags"`
	Msg       string `json:"msg"`
	TimeStamp int32  `json:"timestamp"`
	Action    uint64 `json:"action"`
}

type ActionFilter struct {
	Id       int64  `json:"id"`
	TagId    int64  `json:"tag_id"`
	ParentId int64  `json:"parent_id"`
	TplId    int64  `json:"tpl_id"`
	Version  int    `json:"version"`
	Priority int    `json:"priority"`
	Name     string `json:"name"`
	Metric   string `json:"metric"`
	Tags     string `json:"tags"`
	Func     string `json:"func"`
	Op       string `json:"op"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

const (
	SYS_F_ACTION_SMS = 1 << iota
	SYS_F_ACTION_EMAIL
	SYS_F_ACTION_LOG
	SYS_F_ACTION_DROP
)
