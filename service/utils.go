/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/yubo/gotool/list"
)

/*
func falconTicker(t time.Duration, debug int) <-chan time.Time {
	if debug > 1 {
		return time.NewTicker(t / DEBUG_MULTIPLES).C
	} else {
		return time.NewTicker(t).C
	}
}
*/

func dictedTagstring(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}
	s = strings.Replace(s, " ", "", -1)

	tag_dict := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {
			tag_dict[tag_pair[0]] = tag_pair[1]
		}
	}
	return tag_dict
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func backtrace(i int) string {
	pc, file, no, ok := runtime.Caller(i)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("#%02d %s %s:%d\n", i, details.Name(), file, no)
	}
	return ""
}

func list_entry(l *list.ListHead) *itemEntry {
	return (*itemEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*itemEntry)(nil)).list))))
}

// tags must sorted order by tag key
func tagsMatch(pattern, tags string) bool {
	if len(pattern) == 0 {
		return true
	}

	ts := strings.Split(pattern, ",")
	i := 0
	tags += ","

	for _, t := range ts {
		i = strings.Index(tags[i:], t+",")
		if i < 0 {
			return false
		}

	}
	return true
}
