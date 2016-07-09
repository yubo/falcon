/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yubo/falcon/specs"
)

func falconTicker(t time.Duration, debug int) <-chan time.Time {
	if debug > 1 {
		return time.NewTicker(t / DEBUG_MULTIPLES).C
	} else {
		return time.NewTicker(t).C
	}
}

func renderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func renderDataJson(w http.ResponseWriter, data interface{}) {
	renderJson(w, specs.Dto{Msg: "success", Data: data})
}

func renderMsgJson(w http.ResponseWriter, msg string) {
	renderJson(w, map[string]string{"msg": msg})
}

func autoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		renderMsgJson(w, err.Error())
		return
	}
	renderDataJson(w, data)
}

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

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
func key2filename(baseDir string, key string) string {
	return fmt.Sprintf("%s/%s/%s.rrd", baseDir, key[0:2], key)
}
