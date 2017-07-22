/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	f_network = regexp.MustCompile(`^(tcp)|(unix)+:`)
)

func Md5sum(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
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
