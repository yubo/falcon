/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/yubo/falcon"
)

type Log struct {
	Id       int64
	Module   int64
	ModuleId int64
	UserId   int64
	Action   int64
	Data     string
	Time     time.Time
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func DbLog(o orm.Ormer, uid, module, module_id, action int64, data string) {
	o.Raw("insert log (user_id, module, module_id, action, data) values (?, ?, ?, ?, ?)", uid, module, module_id, action, data).Exec()
}

func array2sql(array []int64) string {
	var ret string
	if len(array) == 0 {
		return "()"
	}

	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("%d,", array[i])
	}
	return fmt.Sprintf("(%s)", ret[:len(ret)-1])
}

func stringscmp(a, b []string) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := 0; i < len(a); i++ {
		if ret = strings.Compare(a[i], b[i]); ret != 0 {
			return
		}
	}
	return
}

func intscmp64(a, b []int64) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}

	_a := make([]int, len(a))
	for i := 0; i < len(_a); i++ {
		_a[i] = int(a[i])
	}

	_b := make([]int, len(b))
	for i := 0; i < len(_b); i++ {
		_b[i] = int(b[i])
	}

	sort.Ints(_a)
	sort.Ints(_b)

	for i := 0; i < len(_a); i++ {
		if ret = _a[i] - _b[i]; ret != 0 {
			return
		}
	}
	return
}

func intscmp(a, b []int) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	sort.Ints(a)
	sort.Ints(b)
	for i := 0; i < len(a); i++ {
		if ret = a[i] - b[i]; ret != 0 {
			return
		}
	}
	return
}

func jsonStr(i interface{}) string {
	if ret, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return string(ret)
	}
}

func MdiffStr(src, dst []string) (add, del []string) {
	_src := make(map[string]bool)
	_dst := make(map[string]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}

func MdiffInt(src, dst []int64) (add, del []int64) {
	_src := make(map[int64]bool)
	_dst := make(map[int64]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}

// FromRequest extracts the user IP address from req, if present.
func FromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

func GetIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}

func getJson(url string, resp interface{}, timeout time.Duration) error {
	cli := &http.Client{Timeout: timeout}
	r, err := cli.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(resp)
}

func postJson(url string, param interface{}, resp interface{}) error {
	cli := &http.Client{Timeout: 60 * time.Second}
	r, err := cli.Post(
		url,
		"application/json",
		bytes.NewBuffer([]byte(jsonStr(param))),
	)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(resp)
}

// cap=xiaomi,owt=inf -> cap.xiaomi_owt.inf
func TagToOld(t string) string {
	ret := make([]byte, len(t))

	copy(ret, t)
	for keyZone, i := true, 0; i < len(ret); i++ {
		if ret[i] == '=' && keyZone {
			ret[i] = '.'
			keyZone = false
		} else if ret[i] == ',' && !keyZone {
			ret[i] = '_'
			keyZone = true
		}
	}
	return string(ret)
}

// cap.xiaomi_owt.inf -> cap=xiaomi,owt=inf
func TagToNew(t string) string {
	ret := make([]byte, len(t))

	copy(ret, t)
	for keyZone, i := true, 0; i < len(ret); i++ {
		if ret[i] == '.' && keyZone {
			ret[i] = '='
			keyZone = false
		} else if ret[i] == '_' && !keyZone {
			ret[i] = ','
			keyZone = true
		}
	}
	return string(ret)
}

func GetKind(v interface{}) reflect.Kind {
	kind := reflect.ValueOf(v).Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		return GetKind(reflect.ValueOf(v).Elem())
	}
	return kind
}

func PruneNilMsg(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	kind := rv.Kind()

	if (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice) &&
		reflect.ValueOf(v).IsNil() {

		switch GetKind(v) {
		case reflect.Slice:
			return []struct{}{}
		case reflect.Struct:
			return struct{}{}
		}
	}
	return v
}

func (op *Operator) RelCheck(sql string, args ...interface{}) (err error) {
	var n int64

	if err = op.O.Raw(sql, args...).QueryRow(&n); err != nil {
		return
	}

	if n > 0 {
		return falcon.ErrInUse
	}
	return nil
}

func (op *Operator) SqlRow(container interface{}, sql string, args ...interface{}) error {
	return op.O.Raw(sql, args...).QueryRow(container)
}

// returns the integer generated by the database
func (op *Operator) SqlInsert(sql string, args ...interface{}) (int64, error) {
	res, err := op.O.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Affected returns the number of rows affected
func (op *Operator) SqlExec(sql string, args ...interface{}) (int64, error) {
	res, err := op.O.Raw(sql, args...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func sqlLimit(sql string, limit, offset int) string {
	if limit > 0 {
		return fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, limit, offset)
	}
	return sql
}

func sqlName(query string) (where string, args []interface{}) {
	sql2 := []string{}
	sql3 := []interface{}{}
	if query != "" {
		sql2 = append(sql2, "name like ? ")
		sql3 = append(sql3, "%"+query+"%")
	}
	if len(sql2) != 0 {
		where = "WHERE " + strings.Join(sql2, " AND ")
		args = sql3
	}
	return
}
