/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package specs

import (
	"crypto/md5"
	"fmt"
	"io"
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
	ls := strings.Split(lines, "\n")
	indent := strings.Repeat(" ", i*IndentSize)
	for _, l := range ls {
		ret += fmt.Sprintf("%s%s\n", indent, l)
	}
	return strings.TrimRight(ret, "\n")
}
