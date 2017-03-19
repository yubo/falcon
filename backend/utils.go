/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/yubo/falcon"
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
	renderJson(w, falcon.Dto{Msg: "success", Data: data})
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
