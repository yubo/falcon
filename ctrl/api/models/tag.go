/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
)

type TagCreate struct {
	Name string `json:"name"`
	Type int64  `json:"-"`
}

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

type Tag_rel struct {
	Id       int64
	TagId    int64
	SupTagId int64
	Offset   int64
}

type TagNode struct {
	Key  string
	Must bool
}

type TagSchema struct {
	data  string
	nodes []TagNode
}

type TreeNode struct {
	TagId    int64       `json:"id"`
	SupTagId int64       `json:"-"`
	Name     string      `json:"name"`
	Label    string      `json:"label"`
	Read     bool        `json:"read"`
	Child    []*TreeNode `json:"children"`
}

type tagNode struct {
	TagId    int64
	SupTagId int64
	Name     string
	Child    []*tagNode
}

const (
	XIAOMI_SCHEMA = "cop,owt,pdl,servicegroup;service,jobgroup;job,sbs;mod;srv;grp;cluster;loc;idc;status;"
)

//cop,owt,pdl,servicegroup;service,jobgroup;job,sbs;mod;srv;grp;cluster;loc;idc;status;
// ',' : must
// ';' : or
func NewTagSchema(tag string) (*TagSchema, error) {
	var (
		i, j int
	)

	if tag == "" {
		return nil, nil
	}

	ret := &TagSchema{data: tag}
	for i, j = 0, 0; j < len(tag); j++ {
		if tag[j] == ',' {
			if i >= j {
				return nil, core.ErrParam
			}
			ret.nodes = append(ret.nodes, TagNode{
				Key:  strings.TrimSpace(tag[i:j]),
				Must: true,
			})
			i = j + 1
		} else if tag[j] == ';' {
			if i >= j {
				return nil, core.ErrParam
			}
			ret.nodes = append(ret.nodes, TagNode{
				Key:  strings.TrimSpace(tag[i:j]),
				Must: false,
			})
			i = j + 1
		}
	}
	if i != j || i == 0 {
		return nil, core.ErrParam
	}

	return ret, nil
}

func tagMap(tag string) (map[string]string, error) {
	var (
		keyZone bool
		i, j    int
		k, v    string
		ret     = make(map[string]string)
	)

	for keyZone, i, j = true, 0, 0; j < len(tag); j++ {
		if tag[j] == '=' && keyZone {
			k = strings.TrimSpace(tag[i:j])
			i = j + 1
			keyZone = false
		} else if tag[j] == ',' && !keyZone {
			v = strings.TrimSpace(tag[i:j])
			if len(k) > 0 && len(v) > 0 {
				ret[k] = v
				k, v = "", ""
			} else {
				return ret, core.ErrParam
			}
			i = j + 1
			keyZone = true
		}
	}

	v = strings.TrimSpace(tag[i:])
	if len(k) > 0 && len(v) > 0 {
		ret[k] = v
		return ret, nil
	} else {
		return ret, core.ErrParam
	}
}

func (ts *TagSchema) Fmt(tag string, force bool) (string, error) {
	var (
		ret string
		n   int
	)

	if tag == "" {
		return "", nil
	}

	m, err := tagMap(tag)
	if err != nil {
		return "", err
	}
	for _, node := range ts.nodes {
		if v, ok := m[node.Key]; ok {
			ret += fmt.Sprintf("%s=%s,", node.Key, v)
			n++
		} else if !force && node.Must {
			return ret, core.ErrParam
		}

		// done
		if n == len(m) {
			return ret[:len(ret)-1], nil
		}
	}

	// some m.key miss match
	if force && len(ret) > 1 {
		return ret[0 : len(ret)-1], nil
	}

	return ret, core.ErrParam
}

func TagRelation(t string) (ret []string) {

	if t == "" {
		return []string{""}
	}

	tags := strings.Split(t, ",")
	if len(tags) < 1 {
		return []string{""}
	}
	ret = make([]string, len(tags)+1)

	for tag, i := "", 1; i < len(ret); i++ {
		tag += tags[i-1] + ","
		ret[i] = tag[:len(tag)-1]
	}
	return ret
}

func TagParents(t string) (ret []string) {

	tags := strings.Split(t, ",")
	if len(tags) < 1 {
		return nil
	}
	ret = make([]string, len(tags))

	for tag, i := "", 1; i < len(ret); i++ {
		tag += tags[i-1] + ","
		ret[i] = tag[:len(tag)-1]
	}
	return ret
}

func TagParent(t string) string {
	if i := strings.LastIndexAny(t, ","); i > 0 {
		return t[:i]
	} else {
		return ""
	}
}

func TagLast(t string) string {
	return t[strings.LastIndexAny(t, ",")+1:]
}

func (op *Operator) createTag(t *TagCreate, schema *TagSchema) (id int64, err error) {
	if schema != nil {
		t.Name, err = schema.Fmt(t.Name, false)
		if err != nil {
			return
		}
	}

	op.O.Begin()
	defer func() {
		if err != nil {
			op.O.Rollback()
		} else {
			op.O.Commit()
		}
	}()

	id, err = op.SqlInsert("insert tag (name, type) values (?, ?)", t.Name, t.Type)
	if err != nil {
		return
	}
	glog.V(4).Infof("%s id: %d tag: %s\n", MODULE_NAME, id, t.Name)

	if rels := TagRelation(t.Name); len(rels) > 0 {
		var (
			ids []int64
			arg = make([]string, len(rels))
		)

		for i, v := range rels {
			arg[i] = "'" + v + "'"
		}

		_, err = op.O.Raw("select id from tag where name in (" + strings.Join(arg, ", ") + ") order by id").QueryRows(&ids)
		if err != nil {
			return
		}

		tag_rels := make([]string, len(ids))
		for i, tid := range ids {
			tag_rels[i] = fmt.Sprintf("(%d, %d, %d)", id, tid, len(ids)-1-i)
		}
		_, err = op.SqlExec("insert tag_rel (tag_id, sup_tag_id, offset) values " + strings.Join(tag_rels, ", "))
		if err != nil {
			return
		}
	}

	moduleCache[CTL_M_TAG].set(id,
		&Tag{Id: id, Name: t.Name, Type: t.Type}, t.Name)
	op.log(CTL_M_TAG, id, CTL_A_ADD, "")

	return id, err
}

func (op *Operator) CreateTag(t *TagCreate) (id int64, err error) {
	if id, err = op.createTag(t, _models.sysTagSchema); err != nil {
		return
	}
	if cacheTree != nil {
		cacheTree.build()
	}
	return
}

func (op *Operator) GetTag(id int64) (ret *Tag, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_TAG].get(id).(*Tag); ok {
		return
	}

	ret = &Tag{}
	err = op.SqlRow(ret, "select id, name, type from tag where id = ?", id)
	if err == nil {
		moduleCache[CTL_M_TAG].set(id, ret, ret.Name)
	}
	return
}

func (op *Operator) GetTagIdByName(tag string) (int64, error) {
	if t, err := op.GetTagByName(tag); err != nil {
		return 0, err
	} else {
		return t.Id, nil
	}
}

func (op *Operator) GetTagByName(tag string) (ret *Tag, err error) {
	var ok bool

	if ret, ok = moduleCache[CTL_M_TAG].getByKey(tag).(*Tag); ok {
		return
	}

	ret = &Tag{}
	err = op.SqlRow(ret, "select id, name, type from tag where name = ?", tag)
	if err == nil {
		moduleCache[CTL_M_TAG].set(ret.Id, ret, ret.Name)
	}
	return
}

func (op *Operator) GetTagsCnt(query string) (cnt int64, err error) {
	sql, sql_args := sqlName(query)
	err = op.O.Raw("SELECT count(*) FROM tag "+sql, sql_args...).QueryRow(&cnt)
	return
}

func (op *Operator) GetTags(query string, limit, offset int) (ret []*Tag, err error) {
	sql, sql_args := sqlName(query)
	sql = sqlLimit("select id, name, type from tag "+sql+" ORDER BY name", limit, offset)
	_, err = op.O.Raw(sql, sql_args...).QueryRows(&ret)
	return
}

/*
func (op *Operator) UpdateTag(id int64, t *Tag) (ret *Tag, err error) {
	_, err = op.SqlExec("update mockcfg set name = ? where id = ?", t.Name, id)
	if err != nil {
		return
	}

	if ret, err = op.GetTag(id); err != nil {
		return
	}

	moduleCache[CTL_M_TAG].set(id, ret, ret.Name)
	op.log(CTL_M_TAG, id, CTL_A_SET, "")
	return ret, err
}
*/

func (op *Operator) DeleteTag(id int64) (err error) {
	if err = op.RelCheck("SELECT count(*) FROM tag_rel where sup_tag_id = ? and offset = 1 ", id); err != nil {
		return errors.New(err.Error() + "(tag - tag)")
	}

	if err = op.RelCheck("SELECT count(*) FROM tag_host where tag_id = ?", id); err != nil {
		return errors.New(err.Error() + "(tag - host)")
	}

	if err = op.RelCheck("SELECT count(*) FROM tpl_rel where tag_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_USER); err != nil {
		return errors.New(err.Error() + "(tag - role - user)")
	}

	if err = op.RelCheck("SELECT count(*) FROM tpl_rel where tag_id = ? and type_id = ?",
		id, TPL_REL_T_ACL_TOKEN); err != nil {
		return errors.New(err.Error() + "(tag - role - token)")
	}

	if err = op.RelCheck("SELECT count(*) FROM tag_tpl where tag_id = ?", id); err != nil {
		return errors.New(err.Error() + "(tag - template)")
	}

	if _, err = op.SqlExec("delete from tag where id = ?", id); err != nil {
		return err
	}
	op.O.Raw("delete FROM tag_rel where tag_id = ?", id).Exec()
	// TODO: clean tpl_rel, tag_host, ...

	// rebuild cache tree
	if cacheTree != nil {
		cacheTree.build()
	}

	if t, ok := moduleCache[CTL_M_TAG].get(id).(*Tag); ok {
		moduleCache[CTL_M_TAG].del(id, t.Name)
	}
	op.log(CTL_M_TAG, id, CTL_A_DEL, "")

	return nil
}

/*******************************************************************************
 ************************ tag - tag ********************************************
 ******************************************************************************/

func cloneTreeNode(t *TreeNode, depth int) (out *TreeNode) {

	if t == nil {
		return
	}

	out = &TreeNode{
		TagId:    t.TagId,
		SupTagId: t.SupTagId,
		Name:     t.Name,
		Label:    t.Label,
		Read:     t.Read,
	}

	if depth == 0 {
		return
	}

	out.Child = make([]*TreeNode, len(t.Child))

	for i, _ := range t.Child {
		out.Child[i] = cloneTreeNode(t.Child[i], depth-1)
	}

	return
}

func pruneTagTree(nodes map[int64]*TreeNode, idx int64) (tree *TreeNode) {
	n, ok := nodes[idx]
	if !ok {
		return nil
	}

	n.Label = n.Name[strings.LastIndexAny(n.Name, ",")+1:]

	return tree
}

func (op *Operator) GetOpTag(expand bool) ([]int64, error) {
	if expand {
		return userHasTokenTagExpend(op.O, op.User.Id, SYS_O_TOKEN)
	} else {
		return userHasTokenTag(op.O, op.User.Id, SYS_O_TOKEN)
	}
}

func (op *Operator) GetReadTag(expand bool) ([]int64, error) {
	if expand {
		return userHasTokenTagExpend(op.O, op.User.Id, SYS_R_TOKEN)
	} else {
		return userHasTokenTag(op.O, op.User.Id, SYS_R_TOKEN)
	}
}

func (op *Operator) GetTreeNode(tagId int64, depth int, direct bool) (tree *TreeNode) {
	var (
		ids, names []string
		isChild    []bool
	)

	if direct {
		return cloneTreeNode(cacheTree.get(tagId), depth)
	}

	_, err := op.O.Raw("SELECT group_concat(d1.sup_tag_id order by d1.`offset` desc) as ids, group_concat(d3.name order by d1.`offset` desc SEPARATOR ',,') as tags from tag_rel d1 join (SELECT c1.tag_id, c1.sup_tag_id, c1.offset from tag_rel c1 join (SELECT distinct b1.user_tag_id FROM (SELECT a1.tag_id AS user_tag_id, a2.tag_id AS token_tag_id, a1.tpl_id AS role_id, a1.sub_id AS user_id, a2.sub_id AS token_id FROM tpl_rel a1 JOIN tpl_rel a2 ON a1.type_id = ? AND a1.sub_id = ? AND a2.type_id = ?  AND a2.sub_id = ? AND a1.tpl_id = a2.tpl_id) b1 JOIN tag_rel b2 ON b1.user_tag_id = b2.tag_id AND b1.token_tag_id = b2.sup_tag_id ) c2 on c1.tag_id = c2.user_tag_id WHERE c1.sup_tag_id = ?) d2 left join tag d3 on d1.sup_tag_id = d3.id WHERE d1.tag_id = d2.tag_id AND d1.offset <= d2.offset GROUP BY d1.tag_id ", TPL_REL_T_ACL_USER, op.User.Id, TPL_REL_T_ACL_TOKEN, SYS_R_TOKEN, tagId).QueryRows(&ids, &names)
	if err != nil {
		return nil
	}

	// set child falg to remove from tree
	isChild = make([]bool, len(names))
	for i := 0; i < len(names); i++ {
		for j := 0; j < i; j++ {
			if strings.HasPrefix(names[i], names[j]+",,") {
				isChild[i] = true
			}
		}
	}

	nmap := make(map[int64]*TreeNode)

	for i := 0; i < len(ids); i++ {
		if isChild[i] {
			continue
		}
		name := strings.Split(names[i], ",,")
		_id := strings.Split(ids[i], ",")
		id := make([]int64, len(_id))

		for j := 0; j < len(id); j++ {
			id[j], _ = strconv.ParseInt(_id[j], 10, 0)
		}

		if !(len(id) > 0 && id[0] > 0) {
			return nil
		}

		for j := 0; j < len(id); j++ {
			if _, ok := nmap[id[j]]; ok {
				continue
			}

			n := &TreeNode{
				TagId: id[j],
				Name:  name[j],
				Label: name[j][strings.LastIndexAny(name[j], ",")+1:],
			}

			if j == len(id)-1 {
				if m := depth - len(id); m > 0 {
					n = cloneTreeNode(cacheTree.get(id[len(id)-1]), m)
				} else {
					n.Read = true
				}
			}
			nmap[id[j]] = n

			if j-1 >= 0 {
				nmap[id[j-1]].Child = append(nmap[id[j-1]].Child, n)
			}
		}

	}
	return nmap[tagId]
}
