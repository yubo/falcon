/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io/ioutil"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

var (
	f_network   = regexp.MustCompile(`^(tcp)|(unix)+:`)
	crc64_table = crc64.MakeTable(crc64.ECMA)
	crc32_table = crc32.MakeTable(0xD5828281)
)

func Md5sum(raw []byte) string {
	h := md5.New()
	h.Write(raw)
	//io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sum64(raw []byte) uint64 {
	h := crc64.New(crc64_table)
	h.Write(raw)
	//io.WriteString(h, raw)
	return h.Sum64()
}

func Sum32(raw []byte) uint32 {
	return crc32.Checksum(raw, crc32_table)
}

func FmtTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func ReadFileInt(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	if i, err := strconv.Atoi(strings.TrimSpace(string(data))); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func GetType(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		//return "*" + t.Elem().Name()
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func PreByte(data []byte, pos int) (int, byte) {
	for i := pos; i >= 0; i-- {
		c := data[i]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		return i, c
	}
	return -1, ' '
}

func IndentLines(i int, lines string) (ret string) {
	ls := strings.Split(strings.Trim(lines, "\n"), "\n")
	indent := strings.Repeat(" ", i*IndentSize)
	for _, l := range ls {
		ret += fmt.Sprintf("%s%s\n", indent, l)
	}
	return string([]byte(ret)[:len(ret)-1])
}

func AddrIsDisable(addr string) bool {
	if addr == "" || addr == "disable" {
		return true
	}
	return false
}

func Dialer(addr string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	return d.Dial(ParseAddr(addr))
}

func ParseAddr(url string) (net, addr string) {
	if f := f_network.Find([]byte(url)); f != nil {
		return url[:len(f)-1], url[len(f):]
	}
	return "tcp", url
}

func sortTags(s []byte) []byte {
	str := strings.Replace(string(s), " ", "", -1)
	if str == "" {
		return []byte{}
	}

	tags := strings.Split(str, ",")
	sort.Strings(tags)
	return []byte(strings.Join(tags, ","))
}

func Override(dst, src interface{}) error {
	srv := reflect.ValueOf(src).Elem()
	srt := srv.Type()

	drv := reflect.ValueOf(dst).Elem()
	drt := drv.Type()

	if !drv.CanSet() {
		return errors.New("dst can't set")
	}

	for i := 0; i < srv.NumField(); i++ {
		fname := srt.Field(i).Name

		if _, ok := drt.FieldByName(fname); !ok {
			continue
		}

		sf := srv.Field(i)
		df := drv.FieldByName(fname)

		if !df.CanSet() {
			continue
		}

		if sf.Type().Kind() != df.Type().Kind() {
			continue
		}

		switch df.Type().Kind() {
		case reflect.Struct, reflect.Map:
			// skip
		default:
			df.Set(sf)
		}

	}
	return nil
}

func NewOrm(name, dsn string, maxIdleConns, maxOpenConns int) (o orm.Ormer, err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err = db.Ping(); err != nil {
		return
	}

	return orm.NewOrmWithDB("mysql", name, db)
}
