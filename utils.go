/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Md5sum(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func FmtTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func IndentLines(i int, lines string) (ret string) {
	ls := strings.Split(strings.Trim(lines, "\n"), "\n")
	indent := strings.Repeat(" ", i*IndentSize)
	for _, l := range ls {
		ret += fmt.Sprintf("%s%s\n", indent, l)
	}
	return string([]byte(ret)[:len(ret)-1])
}

func readFileInt(filename string) (int, error) {
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

func GetType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		//return "*" + t.Elem().Name()
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
