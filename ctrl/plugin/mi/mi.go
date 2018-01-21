/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package mi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
	"golang.org/x/net/context"
)

const (
	C_MI_NORNS_INTERVAL = "minornsinterval"
	C_MI_NORNS_URL      = "minornsurl"

	C_MI_NORNS_INTERVAL_DEF = 5
	C_MI_NORNS_URL_DEF      = "http://norns.dev/api/v1/tagstring/cop.xiaomi/hostinfos"
)

var (
	module = &MiModule{}
)

type MiModule struct {
	ctx        context.Context
	cancel     context.CancelFunc
	miSyncer   miSync
	miNornsUrl string
}

type NornsHosts struct {
	Id       int
	XmanId   string
	Hostname string
	Type     string
	Status   string
	Loc      string
	Idc      string
}

type NornsHostInfo struct {
	Name   string
	Tags   []string
	Detail NornsHosts
}

type miSync struct {
	ctx        *miSyncContent
	localhosts map[string]*Host
}

type miHost struct {
	name   string
	tags   []string
	typ    string
	status string
	loc    string
	idc    string
}

type miNode struct {
	name  string
	host  []*miHost
	child []*miNode
}

func (p *miNode) str(i int) string {
	indent := strings.Repeat(" ", i*2)
	ret := fmt.Sprintf("%s+%s\n", indent, p.name)
	for i := 0; i < len(p.host); i++ {
		ret += fmt.Sprintf("%s%s\n", indent, p.host[i].name)
	}
	for i := 0; i < len(p.child); i++ {
		ret += fmt.Sprintf("%s", p.child[i].str(i+1))
	}
	return ret
}

func (p *miNode) String() string {
	return p.str(0)
}

type miSyncContent struct {
	node map[string]*miNode
	host map[string]*miHost
	h    []*miHost
}

func (p *miSyncContent) getNode(tag string) *miNode {
	node, ok := p.node[tag]
	if ok {
		return node
	}

	node = &miNode{
		name:  tag,
		host:  make([]*miHost, 0),
		child: make([]*miNode, 0),
	}
	p.node[tag] = node

	n := strings.LastIndex(tag, ",")
	if n < 0 {
		return node
	}

	pnode := p.getNode(tag[:n])
	pnode.child = append(pnode.child, node)
	return node
}

func init() {
	ctrl.RegisterModule(&MiModule{})
}

func (p *MiModule) PreStart(ctrl *Ctrl) error {
	ctrl.Hooks.RateLimitsAccess = rateLimitsAccess
	module = p
	return nil
}

func (p *MiModule) Start(ctrl *Ctrl) error {
	conf := &ctrl.Conf.Ctrl
	p.ctx, p.cancel = context.WithCancel(context.Background())

	url := conf.DefaultString(C_MI_NORNS_URL, C_MI_NORNS_URL_DEF)
	interval, _ := conf.DefaultInt(C_MI_NORNS_INTERVAL, C_MI_NORNS_INTERVAL_DEF)

	syncWorker(url, interval, p.ctx)

	return nil
}

func (p *MiModule) Stop(ctrl *Ctrl) error {
	p.cancel()
	return nil
}

func (p *MiModule) Reload(ctrl *Ctrl) error {
	return nil
}

func miGetTagHostByTag(tagstring, query string, localhosts map[string]*Host) (ret []*ctrl.TagHostApiGet, err error) {
	var ok bool
	var tag string

	if strings.IndexRune(tagstring, '=') == -1 {
		err = errors.New("tag should xxx=xxx formats")
		return
	}
	tags := strings.Split(tagstring, ",")

	for _, host := range miSyncer.ctx.h {
		if query != "" && !strings.Contains(host.name, query) {
			continue
		}
		for _, tag = range host.tags {
			ok = true
			for _, s := range tags {
				if !strings.Contains(tag+",", s+",") {
					ok = false
					break
				}
			}
			if ok {
				break
			}
		}
		if !ok {
			continue
		}
		if v, ok := localhosts[host.name]; ok {
			th := &ctrl.TagHostApiGet{
				TagName:  tag,
				HostName: host.name,
				HostId:   v.Id,
			}
			ret = append(ret, th)
		}
	}
	return
}

func miGetTagHostCnt(tag, query string, deep bool) (int64, error) {
	var hosts []*ctrl.TagHostApiGet
	var err error

	if miSyncer.ctx == nil {
		return 0, nil
	}

	hosts, err = miGetTagHostByTag(tag, query, miSyncer.localhosts)
	if err != nil {
		return 0, err
	}
	return int64(len(hosts)), err
}

func miGetTagHost(tag, query string, deep bool,
	limit, offset int) (ret []*ctrl.TagHostApiGet, err error) {

	if miSyncer.ctx == nil {
		return
	}

	ret, err = miGetTagHostByTag(tag, query, miSyncer.localhosts)
	if err != nil {
		return
	}

	if offset > 0 {
		ret = ret[offset:]
	}
	if limit > 0 && limit < len(ret) {
		ret = ret[:limit]
	}

	for i, _ := range ret {
		if h1, err := SysOp.GetHost(ret[i].HostId); err == nil {
			ret[i].Pause = h1.Pause
			ret[i].MaintainBegin = h1.MaintainBegin
			ret[i].MaintainEnd = h1.MaintainEnd
		}
	}
	return
}

func miFetchHosts(filename string) (ctx *miSyncContent, err error) {
	var hosts []NornsHostInfo
	var b []byte

	ctx = &miSyncContent{
		node: make(map[string]*miNode),
		host: make(map[string]*miHost),
	}
	if filename == "" {
		if err = falcon.GetJson(miNornsUrl, &hosts, 2*time.Minute); err != nil {
			glog.Error(MODULE_NAME, err)
			return
		}
	} else {
		if b, err = ioutil.ReadFile(filename); err != nil {
			glog.Error(MODULE_NAME, err)
			return
		} else {
			if err = json.Unmarshal(b, &hosts); err != nil {
				glog.Error(MODULE_NAME, err)
				return
			}

		}
	}

	ctx.h = make([]*miHost, len(hosts))
	for idx, h := range hosts {

		host := &miHost{
			name:   h.Detail.Hostname,
			tags:   h.Tags,
			typ:    h.Detail.Type,
			status: h.Detail.Status,
			loc:    h.Detail.Loc,
			idc:    h.Detail.Idc,
		}
		for k, v := range host.tags {
			host.tags[k] = TagToNew(v)
		}

		for _, t := range h.Tags {
			n := strings.LastIndex(t, "_idc.")
			if n < 0 {
				n = strings.LastIndex(t, "_loc.")
			}
			if n < 0 {
				n = strings.LastIndex(t, "_status.")
			}
			if n < 0 {
				n = len(t)
			}

			node := ctx.getNode(TagToNew(t[:n]))
			node.host = append(node.host, host)
		}
		ctx.host[host.name] = host
		ctx.h[idx] = host
	}

	glog.V(4).Infof("%s get %d node from norns", MODULE_NAME, len(ctx.host))
	return
}

func getLocalHosts() (hm map[string]*Host) {
	var ha []*Host

	hm = make(map[string]*Host)

	for i := 0; ; i++ {
		n, err := ctrl.Orm.Ctrl.Raw("select id, uuid, name, type, status, loc, idc, pause, maintain_begin, maintain_end, create_time from host limit ? offset ?", 100, 100*i).QueryRows(&ha)
		if err != nil || n == 0 {
			break
		}
		for _, h := range ha {
			hm[h.Name] = h
		}
	}

	glog.V(4).Infof("%s get %d hosts from lcoal db", MODULE_NAME, len(hm))
	return
}

func syncHosts(src map[string]*miHost, dst map[string]*Host) {
	var (
		err error
		res sql.Result
	)
	add, del, add_err, del_err := 0, 0, 0, 0
	for _, host := range src {
		if _, ok := dst[host.name]; !ok {
			add++
			h := &Host{
				Name:   host.name,
				Type:   host.typ,
				Status: host.status,
				Loc:    host.loc,
				Idc:    host.idc,
			}
			res, err = ctrl.Orm.Ctrl.Raw("insert host (name, type, status, loc, idc) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", h.Name, h.Type, h.Status, h.Loc, h.Idc).Exec()
			if err != nil {
				add_err++
				continue
			}

			h.Id, err = res.LastInsertId()
			if err != nil {
				add_err++
			} else {
				dst[h.Name] = h
			}
		}
	}

	for _, host := range dst {
		if _, ok := src[host.Name]; !ok {
			del++
			_, err := ctrl.Orm.Ctrl.Raw("DELETE FROM host WHERE name = ?", host.Name).Exec()
			if err != nil {
				del_err++
			} else {
				delete(dst, host.Name)
			}
		}
	}
	glog.V(4).Infof("%s syncHosts(): insert %d/%d del %d/%d",
		MODULE_NAME, add_err, add, del_err, del)

}

func syncWorker(url string, interval int, ctx context.Context) {
	miNornsUrl = url
	miSyncer.localhosts = make(map[string]*Host)
	go func() {
		miSyncer.localhosts = getLocalHosts()
		if ctx, err := miFetchHosts(""); err == nil {
			miSyncer.ctx = ctx
			syncHosts(miSyncer.ctx.host, miSyncer.localhosts)
		}
		ticker := time.NewTicker(time.Minute * time.Duration(interval)).C
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				if ctx, err := miFetchHosts(""); err == nil {
					miSyncer.ctx = ctx
					syncHosts(miSyncer.ctx.host, miSyncer.localhosts)
				}
			}
		}
	}()
}
