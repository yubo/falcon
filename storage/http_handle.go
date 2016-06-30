/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/specs"
)

func count_handler(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().Unix() - 3600
	count := 0

	for _, v := range cache.hash {
		if v.putTs > ts {
			count++
		}
	}
	w.Write([]byte(fmt.Sprintf("%d\n", count)))
}

func recv_hanlder(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Path[len("/api/recv/"):]
	args := strings.Split(urlParam, "/")

	argsLen := len(args)
	if !(argsLen == 6 || argsLen == 7) {
		renderDataJson(w, "bad args")
		return
	}

	endpoint := args[0]
	metric := args[1]
	ts, _ := strconv.ParseInt(args[2], 10, 64)
	step, _ := strconv.ParseInt(args[3], 10, 32)
	dstype := args[4]
	value, _ := strconv.ParseFloat(args[5], 64)
	tags := make(map[string]string)
	if argsLen == 7 {
		tags = dictedTagstring(args[6])
	}

	item := &specs.MetaData{
		Endpoint:    endpoint,
		Metric:      metric,
		Timestamp:   ts,
		Step:        step,
		CounterType: dstype,
		Value:       value,
		Tags:        tags,
	}
	gitem, err := convert2GraphItem(item)
	if err != nil {
		renderDataJson(w, err)
		return
	}

	handleItems([]*specs.GraphItem{gitem})
	renderDataJson(w, "ok")
}

func recv2_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0 && len(r.Form["v"]) > 0 &&
		len(r.Form["ts"]) > 0 && len(r.Form["step"]) > 0 && len(r.Form["type"]) > 0) {
		renderDataJson(w, "bad args")
		return
	}
	endpoint := r.Form["e"][0]
	metric := r.Form["m"][0]
	value, _ := strconv.ParseFloat(r.Form["v"][0], 64)
	ts, _ := strconv.ParseInt(r.Form["ts"][0], 10, 64)
	step, _ := strconv.ParseInt(r.Form["step"][0], 10, 32)
	dstype := r.Form["type"][0]

	tags := make(map[string]string)
	if len(r.Form["t"]) > 0 {
		tagstr := r.Form["t"][0]
		tagVals := strings.Split(tagstr, ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}

	item := &specs.MetaData{
		Endpoint:    endpoint,
		Metric:      metric,
		Timestamp:   ts,
		Step:        step,
		CounterType: dstype,
		Value:       value,
		Tags:        tags,
	}
	gitem, err := convert2GraphItem(item)
	if err != nil {
		renderDataJson(w, err)
		return
	}

	handleItems([]*specs.GraphItem{gitem})
	renderDataJson(w, "ok")
}

func updateAll_handler(w http.ResponseWriter, r *http.Request) {
	//	go UpdateIndexAllByDefaultStep()
	renderDataJson(w, "ok")
}

func updateAll_concurrent_handler(w http.ResponseWriter, r *http.Request) {
	//renderDataJson(w, GetConcurrentOfUpdateIndexAll())
	renderDataJson(w, "ok")
}

/*
func update_handler(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Path[len("/index/update/"):]
	args := strings.Split(urlParam, "/")

	argsLen := len(args)
	if !(argsLen == 4 || argsLen == 5) {
		renderDataJson(w, "bad args")
		return
	}
	endpoint := args[0]
	metric := args[1]
	step, _ := strconv.ParseInt(args[2], 10, 32)
	dstype := args[3]
	tags := make(map[string]string)
	if argsLen == 5 {
		tagVals := strings.Split(args[4], ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}
	err := UpdateIndexOne(endpoint, metric, tags, dstype, int(step))
	if err != nil {
		renderDataJson(w, fmt.Sprintf("%v", err))
		return
	}

	renderDataJson(w, "ok")
}
*/

func reload_handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
		applyConfigFile(&configOpts, configFile)
		renderDataJson(w, "ok")
	} else {
		renderDataJson(w, "no privilege")
	}
}

func stat_handler(w http.ResponseWriter, r *http.Request) {
	renderDataJson(w, stat_handle())
}

func history_handler(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Path[len("/history/"):]
	args := strings.Split(urlParam, "/")

	argsLen := len(args)
	endpoint := args[0]
	metric := args[1]
	tags := make(map[string]string)
	if argsLen > 2 {
		tagVals := strings.Split(args[2], ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}
	renderDataJson(w, getAllItems(specs.Checksum(endpoint, metric, tags)))
}

func history2_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0) {
		renderDataJson(w, "bad args")
		return
	}
	endpoint := r.Form["e"][0]
	metric := r.Form["m"][0]

	tags := make(map[string]string)
	if len(r.Form["t"]) > 0 {
		tagstr := r.Form["t"][0]
		tagVals := strings.Split(tagstr, ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}

	renderDataJson(w, getAllItems(specs.Checksum(endpoint, metric, tags)))
}

func last_handler(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Path[len("/last/"):]
	args := strings.Split(urlParam, "/")

	argsLen := len(args)
	endpoint := args[0]
	metric := args[1]
	tags := make(map[string]string)
	if argsLen > 2 {
		tagVals := strings.Split(args[2], ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}
	renderDataJson(w, getLastItem(specs.Checksum(endpoint, metric, tags)))
}

func last2_handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !(len(r.Form["e"]) > 0 && len(r.Form["m"]) > 0) {
		renderDataJson(w, "bad args")
		return
	}
	endpoint := r.Form["e"][0]
	metric := r.Form["m"][0]

	tags := make(map[string]string)
	if len(r.Form["t"]) > 0 {
		tagstr := r.Form["t"][0]
		tagVals := strings.Split(tagstr, ",")
		for _, tag := range tagVals {
			tagPairs := strings.Split(tag, "=")
			if len(tagPairs) == 2 {
				tags[tagPairs[0]] = tagPairs[1]
			}
		}
	}

	renderDataJson(w, getLastItem(specs.Checksum(endpoint, metric, tags)))
}

func convert2GraphItem(d *specs.MetaData) (*specs.GraphItem, error) {
	item := &specs.GraphItem{}

	item.Endpoint = d.Endpoint
	item.Metric = d.Metric
	item.Tags = d.Tags
	item.Timestamp = d.Timestamp
	item.Value = d.Value
	item.Step = int(d.Step)
	if item.Step < MIN_STEP {
		item.Step = MIN_STEP
	}
	item.Heartbeat = item.Step * 2

	if d.CounterType == GAUGE {
		item.DsType = d.CounterType
		item.Min = "U"
		item.Max = "U"
	} else if d.CounterType == COUNTER {
		item.DsType = DERIVE
		item.Min = "0"
		item.Max = "U"
	} else if d.CounterType == DERIVE {
		item.DsType = DERIVE
		item.Min = "0"
		item.Max = "U"
	} else {
		return item, fmt.Errorf("not_supported_counter_type")
	}

	item.Timestamp = item.Timestamp - item.Timestamp%int64(item.Step)

	return item, nil
}

func httpRoutes() {
	http.HandleFunc("/count", count_handler)

	// 接收数据 endpoint metric ts step dstype value [tags]
	http.HandleFunc("/api/recv/", recv_hanlder)

	http.HandleFunc("/v2/api/recv", recv2_handler)

	http.HandleFunc("/index/updateAll", updateAll_handler)

	// 获取索引全量更新的并行数
	http.HandleFunc("/index/updateAll/concurrent", updateAll_concurrent_handler)

	// 更新一条索引数据,用于手动建立索引 endpoint metric step dstype tags
	//http.HandleFunc("/index/update/", update_handler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(falcon.VERSION))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		renderDataJson(w, config())
	})

	http.HandleFunc("/config/reload", reload_handler)

	http.HandleFunc("/stat", stat_handler)

	// items.history
	http.HandleFunc("/history/", history_handler)

	http.HandleFunc("/v2/history", history2_handler)

	// items.last
	http.HandleFunc("/last/", last_handler)

	http.HandleFunc("/v2/last", last2_handler)

}
