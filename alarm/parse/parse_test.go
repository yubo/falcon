/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package parse

import (
	"testing"

	"github.com/golang/glog"
)

var text = []byte(`
{
	debug;
	leasekey	"/open-falcon/alarm/online/test.falcon";
	leasettl	20;
	apiAddr		"unix:./var/alarm.rpc";
	httpAddr	"unix:./var/alarm.rpc.gw";
	workerProcesses	2;	// upstream 的并发连接数
	burstSize	2;	// client put burst size to remote service
	callTimeout	5000;	// 请求超时时间
	upstream	stdout;
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
